SOURCE_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")

markdown_stew: $(SOURCE_FILES) ./index.templ index_templ.go
	go build -o ./markdown_stew ./cmd/markdown_stew/main.go
index_templ.go: index.templ
	templ generate -f index.templ
clean:
	rm -f ./markdown_stew *.html
install: ./markdown_stew
	cp ./markdown_stew ~/.local/bin
