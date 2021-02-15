binary = m-check

build:
	go build -o cmd/$(binary) cmd/m-check/m-check.go

run:
	go run cmd/m-check/m-check.go

compile:
	# Cross compilation for building the m-check binary
	GOOS=windows GOARCH=amd64 go build -o ./cmd/$(binary)_windows_amd64.exe cmd/m-check/m-check.go
	GOOS=linux GOARCH=amd64 go build -o ./cmd/$(binary)_linux_amd64 cmd/m-check/m-check.go
	GOOS=darwin GOARCH=amd64 go build -o ./cmd/$(binary)_darwin_amd64 cmd/m-check/m-check.go