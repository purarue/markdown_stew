package markdown_stew

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Template struct {
	Filename string
	Raw      string
	Rendered string
	Slug     string
	Title    string
}

func slugify(basename string) string {
	return strings.ReplaceAll(strings.ToLower(basename), " ", "-")
}

func titleFromFilename(basename string) string {
	// remove extension
	ext := filepath.Ext(basename)
	name := strings.TrimSuffix(basename, ext)
	// replace punctuation with spaces
	for _, p := range []string{".", "_", "-"} {
		name = strings.ReplaceAll(name, p, " ")
	}
	return name
}

func plaintextToHtml(plaintext string) string {
	lines := strings.Split(plaintext, "\n")
	var buf bytes.Buffer
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		buf.WriteString(fmt.Sprintf("<p>%s</p>", line))
	}
	return buf.String()
}

func ReadTemplate(filename string) (*Template, error) {
	basename := filepath.Base(filename)
	// read file into string
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// get extension, if it's .md, then render markdown into html
	// if it's .html, then just read the file into a string

	var rendered string
	ext := filepath.Ext(filename)
	if ext == ".md" {
		// render markdown into html
		var buf bytes.Buffer

		renderer := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
		)

		err := renderer.Convert(content, &buf)
		if err != nil {
			return nil, err
		}
		rendered = buf.String()
	} else if ext == ".html" {
		rendered = string(content)
	} else if ext == ".txt" {
		rendered = plaintextToHtml(string(content))
	} else {
		fmt.Fprintf(os.Stderr, "Unsupported file type for file %s: %s, assuming plaintext\n", filename, ext)
		rendered = plaintextToHtml(string(content))
	}
	return &Template{
		Filename: filename,
		Raw:      string(content),
		Rendered: rendered,
		Slug:     slugify(basename),
		Title:    titleFromFilename(basename),
	}, nil
}
