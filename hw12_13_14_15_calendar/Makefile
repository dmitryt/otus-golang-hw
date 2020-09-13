cdir = $(shell pwd)

build:
	mkdir -p bin
	go build -o ./bin ./...

test:
	go test -race -v `go list ./... | grep -v integration-tests`

lint:
	golangci-lint run ./...

launch:
	./bin/calendar -config=./configs/calendar.yml & \
	./bin/calendar_scheduler -config=./configs/scheduler.yml & \
	./bin/calendar_sender -config=./configs/sender.yml

run:
	make build && make launch

migrate:
	goose -dir ./migrations postgres "user=\${DB_USER} password=\${DB_PASSWORD} dbname=\${DB_NAME} sslmode=disable\" up

up:
	docker-compose -f deployments/docker-compose.yml up -d

down:
	docker-compose -f deployments/docker-compose.yml down

integration-tests:
	docker-compose -f deployments/docker-compose.yml up -d --build
	ginkgo ./integration-tests
	docker-compose -f deployments/docker-compose.yml down

generate:
	protoc \
	-I /usr/local/include \
	-I ${cdir}/service/schema \
	--go_out=plugins=grpc:${cdir}/service/server \
	--grpc-gateway_out=logtostderr=true:${cdir}/service/server \
	${cdir}/service/schema/calendar.proto

.PHONY: build
