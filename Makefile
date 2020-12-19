.PHONY: build
.PHONY: run
.PHONY: build-run

build:
	@go build -o ./build/syncwave ./cmd/syncwave

run:
	@./build/syncwave execute --config-file config.yaml --scenario scenario.yaml

build-run: build run

