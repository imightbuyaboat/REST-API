.PHONY: build run test docker-up docker-down

build:
	go build -o bin/app .

run:
	go run .

test:
	go test ./tests

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down -v