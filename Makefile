IMAGE = mprokopov/getversion
TAG = latest
DOCKER_IMAGE = $(IMAGE):$(TAG)

all: mac linux windows

linux: getversion.linux

windows: getversion.exe

mac: getversion.mac

docker: getversion.linux Dockerfile
	docker build -t $(DOCKER_IMAGE) .

docker-push:
	docker push $(DOCKER_IMAGE)

getversion.linux: main.go
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o getversion.linux

getversion.mac: main.go
	go build -o getversion.mac
	ln -sf getversion.mac getversion

getversion.exe: main.go
	GOOS=windows GOARCH=386 go build -o getversion.exe

clean:
	rm -f getversion.mac getversion.linux getversion getversion.exe
