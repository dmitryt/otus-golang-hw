cdir = $(shell pwd)

build:
	mkdir -p bin
	go build -o ./bin ./...

test:
	go test -race ./...

lint:
	golangci-lint run ./...

run:
	make build && ./bin/calendar -config=./configs/calendar.yml

migrate:
	./scripts/migrations.sh

generate:
	protoc \
	-I /usr/local/include \
	-I ${cdir}/service/schema \
	--go_out=plugins=grpc:${cdir}/service/server \
	--grpc-gateway_out=logtostderr=true:${cdir}/service/server \
	${cdir}/service/schema/calendar.proto

.PHONY: build
