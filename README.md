# markdown_stew

A single-page HTML generator that takes `.md`, `.txt` or `.html` files as input, and creates a static file with no dependencies which lets you switch between rendered versions of those files.

So - it just throws all the files you give it into a single file, hence the name.

I use this for notes or documentation, where I want to be able to switch between a raw/rendered versions of the file, or [create a page I can scroll through on my phone](https://github.com/seanbreckenridge/dotfiles/blob/master/.local/scripts/notes_rendered/Makefile), but don't want to use a full-blown wiki or CMS.

As an example of the output, see [here](https://sean.fish/p/markdown_stew_example.html)

<img src="https://raw.githubusercontent.com/seanbreckenridge/markdown_stew/main/.github/demo.png" width="500px">

If you want to embed all the images into the HTML file itself, you can use [`monolith`](https://github.com/Y2Z/monolith) on the output HTML file.

## Install

```sh
go install 'github.com/seanbreckenridge/markdown_stew/cmd/markdown_stew@latest'
```

Manually:

```sh
git clone https://github.com/seanbreckenridge/markdown_stew
cd ./markdown_stew
go build -o markdown_stew ./cmd/markdown_stew
# copy onto your $PATH somewhere
cp ./markdown_stew ~/.local/bin
```

If you want to hack on the template, see [templ docs](https://templ.guide/quick-start/installation), `make` will update the generated golang template file [`index_templ.go`](./index_templ.go)

## Usage

```
Usage: markdown_stew [flags] <files...>
  -css string
    	raw css string to embed into page
  -dark-mode
    	default to dark mode
  -language string
    	language for HTML page (default "en")
  -title string
    	title to use
```

To re-render on save, I use [`entr`](http://entrproject.org/):

```sh
ls *.md | entr -c sh -c 'markdown_stew *.md >out.html; reload-browser'
```
