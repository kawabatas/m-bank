.PHONY: create
## create: creates db and inserts seed
create:
	go run ./cmd/create

.PHONY: serve
## serve: runs server
serve:
	# go run main.go
	go build
	./m-bank

.PHONY: test
## test: runs test
test:
	go test ./...

.PHONY: gen-rest
## gen-rest: generates server code from swagger spec
gen-rest:
	# configure_xxx.goに編集を加えていないので削除
	rm -r gen/
	mkdir gen/
	# ./gen/* 自動生成コードの出力先(configure_xxx.go以外は編集しない)
	swagger generate server -f swagger.yml -A bank -t gen --exclude-main --strict-additional-properties

.PHONY: gen-mock
## gen-mock: generates gomock
gen-mock:
	mockgen -destination=domain/mock/balance_repository.go -package=mock github.com/kawabatas/m-bank/domain/repository BalanceRepository
	mockgen -destination=domain/mock/payment_transaction_repository.go -package=mock github.com/kawabatas/m-bank/domain/repository PaymentTransactionRepository

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
