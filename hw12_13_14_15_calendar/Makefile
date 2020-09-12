cdir = $(shell pwd)

build:
	mkdir -p bin
	go build -o ./bin ./...

test:
	go test -race ./...

lint:
	golangci-lint run ./...

launch:
	./bin/calendar -config=./configs/calendar.yml & \
	./bin/calendar_scheduler -config=./configs/scheduler.yml & \
	./bin/calendar_sender -config=./configs/sender.yml

run:
	make build && make launch

migrate:
	./scripts/migrations.sh

up:
	docker-compose -f deployments/docker-compose.yml up -d
	
down:
	docker-compose -f deployments/docker-compose.yml down

generate:
	protoc \
	-I /usr/local/include \
	-I ${cdir}/service/schema \
	--go_out=plugins=grpc:${cdir}/service/server \
	--grpc-gateway_out=logtostderr=true:${cdir}/service/server \
	${cdir}/service/schema/calendar.proto

.PHONY: build
