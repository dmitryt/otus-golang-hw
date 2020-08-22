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

.PHONY: build
