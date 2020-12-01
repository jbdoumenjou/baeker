.PHONY: clean

export GO111MODULE=on

GOFILES := $(shell git ls-files '*.go')

clean:
	rm -rf out/ dist/

test: clean
	go test -v -cover ./...

build: clean
	go build -o dist/baeker

check:
	golangci-lint run

imports:
	@goimports -w $(GOFILES)

fmt:
	@gofmt -s -l -w $(GOFILES)

