.PHONY: check
check:
	go fmt ./...
	golangci-lint run

.PHONY: test
test:
	go test -v ./...