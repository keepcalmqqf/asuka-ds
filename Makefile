.PHONY: check test vet fmt site

test:
	go test -count=1 ./...

vet:
	go vet ./...

fmt:
	gofmt -l .

site:
	asuka build --platform web site

check: fmt vet test
	asuka vet ./...
	asuka build --platform web site
