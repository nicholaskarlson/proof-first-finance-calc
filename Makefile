.PHONY: test fmt build demo clean

VERSION ?= dev

test:
	go test -count=1 ./...

fmt:
	gofmt -w cmd internal tests

build:
	mkdir -p bin
	go build -ldflags "-X main.version=$(VERSION)" -o bin/fincalc ./cmd/fincalc

demo:
	go run ./cmd/fincalc demo --out ./out/demo

clean:
	rm -rf ./bin ./out
