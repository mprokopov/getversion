all: main.go
	go build
linux: main.go
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build

windows: main.go
	GOOS=windows GOARCH=386 go build -o getVersion.exe
