package main

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM, // GitHub Flavored Markdown: tables, strikethrough, task lists
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(), // line breaks in markdown become <br> in HTML
		html.WithXHTML(),     // self-closing tags like <br />
	),
)

func RenderMarkdown(body string) (string, error) {

	var buf bytes.Buffer
	if err := md.Convert([]byte(body), &buf); err != nil {
		return "", fmt.Errorf("rendering markdown: %w", err)
	}

	return buf.String(), nil
}
