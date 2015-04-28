GO ?= go
GOPATH := $(CURDIR)/_vendor:$(GOPATH)

all: clean build image

clean:
	$(GO) clean
build:
	CGO_ENABLED=0 GOOS=linux $(GO) build -a -tags netgo -ldflags '-w'
image:
	docker build -t rlister/docker-slacks .
push:
	docker push rlister/docker-slacks
