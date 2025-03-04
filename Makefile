


run:
	go run main.go --env=local

infra:
	docker-compose up -d

infra-down:
	docker-compose down



