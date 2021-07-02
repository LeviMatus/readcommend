.PHONY: all build install clean test
all: clean test build

clean:
	go clean ./...
	go clean -testcache
	rm -f readcommend

test:
	go test ./... -bench=.

test-vendor:
	go test -mod vendor ./... -bench=.

build: clean test
	go build -o readcommend -v main.go

build-vendor: clean test-vendor
	go build -mod vendor -o readcommend -v main.go

install:
	go build -i -o readcommend main.go
	mv readcommend $(GOPATH)/bin

install-vendor:
	go build -mod vendor -i -o readcommend main.go
	mv readcommend $(GOPATH)/bin