.PHONY: check
check:
	go fmt ./...
	golangci-lint run

.PHONY: test
test:
	go test -v ./...

.PHONY: generate
generate:
	oapi-codegen --config=./api/config.yaml ./api/openapi.yaml

.PHONY: up
up:
	docker-compose up --build -d
	go run main.go