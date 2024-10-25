package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/purarue/markdown_stew"
)

type Config struct {
	files      []string
	language   string
	title      string
	darkMode   bool
	embedCss   string
	faviconUrl string
}

func ParseConfig() *Config {
	language := flag.String("language", "en", "language for HTML page")
	title := flag.String("title", "", "title to use")
	darkMode := flag.Bool("dark-mode", false, "default to dark mode")
	faviconUrl := flag.String("favicon-url", "favicon.ico", "URL for favicon")
	embedCss := flag.String("css", "", "raw css string to embed into page")
	flag.Usage = func() {
		// print usage
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <files...>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	files := flag.Args()
	if len(files) < 1 {
		fmt.Fprintf(os.Stderr, "Error: no files provided\n")
		flag.Usage()
		os.Exit(1)
	}

	embedCssWrapped := ""
	if len(strings.TrimSpace(*embedCss)) > 0 {
		embedCssWrapped = fmt.Sprintf("<style>%s</style>", *embedCss)
	}

	var titleTrimmed string
	if len(strings.TrimSpace(*title)) == 0 {
		titleTrimmed = ""
	} else {
		titleTrimmed = *title
	}

	return &Config{
		files:      files,
		language:   *language,
		title:      titleTrimmed,
		darkMode:   *darkMode,
		faviconUrl: *faviconUrl,
		embedCss:   embedCssWrapped,
	}
}

func stew() error {
	config := ParseConfig()

	picoText, err := markdown_stew.PicoCSSText()
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading pico.css: %s", err))
	}

	var tmpls []markdown_stew.Template
	for _, file := range config.files {
		tl, err := markdown_stew.ReadTemplate(file)
		if err != nil {
			return errors.New(fmt.Sprintf("Error reading file: %s", err))
		}
		tmpls = append(tmpls, *tl)
	}

	// check for duplicate slugs
	// map from slug to name of file
	seen := make(map[string]string)
	for _, tmpl := range tmpls {
		if _, ok := seen[tmpl.Slug]; ok {
			return errors.New(fmt.Sprintf("Duplicate slug: %s, used in both: %s and %s", tmpl.Slug, seen[tmpl.Slug], tmpl.Filename))
		}
		seen[tmpl.Slug] = tmpl.Filename
	}

	picoWrapped := fmt.Sprintf("<style>%s</style>", picoText)

	template := markdown_stew.Index(tmpls, config.title, config.language, config.darkMode, picoWrapped, config.embedCss, config.faviconUrl)
	template.Render(context.Background(), os.Stdout)
	return nil
}

func main() {
	if err := stew(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
