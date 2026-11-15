.PHONY: check test vet fmt website

test:
	go test -count=1 ./...

vet:
	go vet ./...

fmt:
	gofmt -l .

website:
	asuka build --platform web website

check: fmt vet test
	asuka vet ./...
	asuka build --platform web website
