bd-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bd-linux-amd64 cmd/bd/main.go

bd-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bd-darwin-amd64 cmd/bd/main.go
