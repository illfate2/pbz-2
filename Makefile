bin/museum:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/museum ./cmd/

.PHONY: dep
dep:
	go mod tidy
