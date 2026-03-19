package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const reloadScript = `
<script>
(function() {
	var lastModified = null;
	setInterval(function() {
		fetch('/__reload')
			.then(function(r) { return r.text(); })
			.then(function(t) {
				if (lastModified && lastModified !== t) {
					window.location.reload();
				}
				lastModified = t;
			});
	}, 500);
})();
</script>
`

type reloadState struct {
	mu      sync.Mutex
	version string
}

func (r *reloadState) bump() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.version = time.Now().Format(time.RFC3339Nano)
}

func (r *reloadState) current() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.version
}

func Serve(cfg *Config, contentDir, outputDir, templatesDir, staticDir string) error {
	state := &reloadState{}
	state.bump()

	// Initial build
	if err := Build(contentDir, outputDir, templatesDir, cfg); err != nil {
		return fmt.Errorf("initial build: %w", err)
	}

	// Watch for changes in a background goroutine
	go watch([]string{contentDir, templatesDir, staticDir}, func() {
		fmt.Println("change detected, rebuilding...")
		if err := Build(contentDir, outputDir, templatesDir, cfg); err != nil {
			fmt.Printf("rebuild failed: %v\n", err)
		}
		state.bump()
	})

	mux := http.NewServeMux()

	// Reload endpoint polled by the browser script
	mux.HandleFunc("/__reload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, state.current())
	})

	// Serve output directory with reload script injected
	mux.Handle("/", injectReload(http.FileServer(http.Dir(outputDir))))

	addr := ":3000"
	fmt.Printf("serving at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func injectReload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only inject into HTML requests
		ext := filepath.Ext(r.URL.Path)
		if ext != "" && ext != ".html" {
			next.ServeHTTP(w, r)
			return
		}

		// Use 200 as default status to avoid "WriteHeader code 0" panic
		rec := &responseRecorder{header: make(http.Header), status: 200}
		next.ServeHTTP(rec, r)

		body := append(rec.body, []byte(reloadScript)...)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(rec.status)
		w.Write(body)
	})
}

type responseRecorder struct {
	header http.Header
	body   []byte
	status int
}

func (r *responseRecorder) Header() http.Header    { return r.header }
func (r *responseRecorder) WriteHeader(status int) { r.status = status }
func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return len(b), nil
}

// watch polls the given directories for file changes and calls onChange when detected.
func watch(dirs []string, onChange func()) {
	modTimes := map[string]time.Time{}

	for {
		time.Sleep(500 * time.Millisecond)

		changed := false
		for _, dir := range dirs {
			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				prev, ok := modTimes[path]
				if !ok || info.ModTime().After(prev) {
					modTimes[path] = info.ModTime()
					if ok {
						changed = true
					}
				}
				return nil
			})
		}

		if changed {
			onChange()
		}
	}
}
