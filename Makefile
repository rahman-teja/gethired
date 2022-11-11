.PHONY: install test-dev test cover run-dev

install:
	go mod download

run-dev:
	go run ./cmd/api/ up

create-database:
	docker run -p 3306:3306 -e MYSQL_DATABASE=upload -e MARIADB_ALLOW_EMPTY_ROOT_PASSWORD=1 -d mariadb

build:
	go build -o bin/api ./cmd/api 
