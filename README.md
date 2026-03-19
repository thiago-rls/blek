# Blek

I build static site generators when I want to dig deeper into a new language. It's the project I know well enough to focus on the language itself rather than figuring out what to build. Built it in Go, liked it enough, now it's my blog.

It does what I need, nothing more. 

## What it does

- `build`, `serve`, `clean`, `new`, and `init` commands
- Development server with live reload
- RSS feed (`feed.xml`) generated automatically from posts
- Go's `html/template` — no new syntax to learn, full control
- Posts (listed, dated) and standalone pages (About, Projects, whatever you want)
- Markdown content

## Installation
```bash
go install github.com/thiago-rls/blek@latest
```

If `blek` isn't found after installation, your Go bin directory probably isn't in your `PATH`:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Getting started
```bash
mkdir my-blog && cd my-blog
blek init .
blek new post "First Post"
blek new page "About" (Optional)
blek serve
```

When you're ready to publish:
```bash
blek build
```

Output goes to `output/`.

## Commands

| Command | Description |
|---------|-------------|
| `blek init [dir]` | Set up a new project with default structure and config |
| `blek build` | Build the site into `output/` |
| `blek serve` | Dev server with auto-rebuild and browser reload |
| `blek new post "Title"` | New post in `content/posts/` |
| `blek new page "Title"` | New standalone page |
| `blek clean` | Delete `output/` |
| `blek version` | Print the current version |

## License

MIT — do whatever you want with it.
