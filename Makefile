.PHONY:doc-api certs test
.DEFAULT_GOAL=run

OUTPUT=out/bin
BIN_NAME=axis

build:
	mkdir -p ${OUTPUT}

tdd:
	go test -v -cover -bench=. ./...

test-coverage:
	go test ./.. -coverprofile=coverage.out
	go tool cover -html=coverage.out

bench:
	go test ./... -v -bench=. -benchmem

certs:
	openssl genrsa -out certs/id_rsa 4096
	openssl rsa -in certs/id_rsa -pubout -out certs/id_rsa.pub
