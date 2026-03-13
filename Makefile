.PHONY: install fmt fmt-check lint test build check clean

install:
	go install .

fmt:
	goimports -w .

fmt-check:
	@test -z "$$(goimports -d .)" || (goimports -d . && exit 1)

lint:
	golangci-lint run ./...

test:
	go test ./...

build:
	go build -o huey .

check: fmt-check lint test

clean:
	rm -f huey
	go clean
