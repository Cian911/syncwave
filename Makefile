.PHONY: build
.PHONY: run
.PHONY: build-run

build:
	@go build -o syncwave ./cmd/syncwave

run:
	@./syncwave execute --config-file config.yaml --scenario scenario.yaml

build-run: build run

