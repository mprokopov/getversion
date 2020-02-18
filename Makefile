all: main.go
	go build
linux: main.go
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build
