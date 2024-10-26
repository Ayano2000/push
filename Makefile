build:
	go build cmd/push/main.go -o bin/push

up:
	docker-compose up -d

down:
	docker-compose down -v