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
	goose -dir ./migrations postgres "user=db_calendar_user password=db_calendar_pass dbname=db_calendar sslmode=disable" up

.PHONY: build
