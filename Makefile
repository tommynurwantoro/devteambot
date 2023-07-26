## default arguments
config = "./config.yaml"
configLocal = "./config.local.yaml"

## test: Test golang sources code
test:
	go test -cover ./... -count=1

## install: Install module requirement applications
install:
	go mod tidy

## build: Build binary applications
build:
	go build -o bin/devteambot main.go

## run: Run binary applications but download module first
run: install build
	./bin/devteambot svc --config=$(config)

## dev: Run binary applications without download module first
dev: build
	./bin/devteambot svc --config=$(configLocal)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run with parameter options: "
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
