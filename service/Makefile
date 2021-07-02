.PHONY: all build install clean test
all: clean test build

clean:
	go clean ./...
	goclean -testcache
	rm -f readcommend

test:
	go test -v ./...

build: clean
	go build -o readcommend -v service/cmd/api/main.go

install: clean test
	go build -i -o readcommend service/cmd/api/main.go
	mv readcommend $(GOPATH)/bin