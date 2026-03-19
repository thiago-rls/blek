package main

import (
	"fmt"
	"os"
)

func main() {
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("loading config: %v\n", err)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "--serve" {
		startServer(cfg)
		return
	}

	if err := Build("content", "output", "templates", cfg); err != nil {
		fmt.Printf("build failed: %v\n", err)
		return
	}

	fmt.Println("done")
}
