binary_name := "gh-purview"
waybar_binary := "gh-purview-waybar"
build_dir := "bin"
main_path := "./cmd/app"
waybar_path := "./cmd/waybar"

go := "go"
goflags := "-trimpath"
ldflags := "-ldflags '-s -w'"

default: build

build: clean fmt vet test
    {{ go }} build {{ goflags }} {{ ldflags }} -o {{ build_dir }}/{{ binary_name }} {{ main_path }}

build-waybar: clean fmt vet test
    {{ go }} build {{ goflags }} {{ ldflags }} -o {{ build_dir }}/{{ waybar_binary }} {{ waybar_path }}

install-waybar: build-waybar
    {{ go }} build {{ goflags }} {{ ldflags }} -o `go env GOPATH`/bin/{{ waybar_binary }} {{ waybar_path }}
    @echo "Installed {{ waybar_binary }} to `go env GOPATH`/bin"

clean:
    rm -rf {{ build_dir }}

fmt:
    {{ go }} fmt ./...

vet:
    {{ go }} vet ./...

test:
    {{ go }} test -v -race ./...

install: build
    @mkdir -p {{ env_var('HOME') }}/.local/share/gh/extensions/{{ binary_name }}
    @cp {{ build_dir }}/{{ binary_name }} {{ env_var('HOME') }}/.local/share/gh/extensions/{{ binary_name }}/{{ binary_name }}
    @echo "Installed gh extension: {{ binary_name }}"

uninstall:
    rm -rf {{ env_var('HOME') }}/.local/share/gh/extensions/{{ binary_name }}
    @echo "Uninstalled gh extension: {{ binary_name }}"

deps:
    {{ go }} mod download
    {{ go }} mod tidy

help:
    @echo "Available recipes:"
    @echo "  build           - Build TUI for current platform"
    @echo "  build-waybar    - Build Waybar module for current platform"
    @echo "  clean           - Remove build artifacts"
    @echo "  fmt             - Format code"
    @echo "  vet             - Run go vet"
    @echo "  test            - Run tests"
    @echo "  deps            - Download and tidy dependencies"
    @echo "  install         - Install TUI as gh extension locally"
    @echo "  install-waybar  - Install Waybar module to GOBIN"
    @echo "  uninstall       - Remove local gh extension"
