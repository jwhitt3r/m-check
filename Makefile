binary = gondola

build:
	go build -o cmd/$(binary) cmd/gondola/main.go

run:
	go run cmd/gondola/main.go

compile:
	# Cross compilation for building the gondola binary
	GOOS=windows GOARCH=amd64 go build -o ./cmd/$(binary)_windows_amd64.exe cmd/gondola/main.go
	GOOS=linux GOARCH=amd64 go build -o ./cmd/$(binary)_linux_amd64 cmd/gondola/main.go
	GOOS=darwin GOARCH=amd64 go build -o ./cmd/$(binary)_darwin_amd64 cmd/gondola/main.go
