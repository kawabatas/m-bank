.PHONY: create
## create: creates db and inserts seed
create:
	go run ./cmd/create

.PHONY: serve
## serve: runs server
serve:
	go run main.go

.PHONY: test
## test: runs test
test:
	go test ./...

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
