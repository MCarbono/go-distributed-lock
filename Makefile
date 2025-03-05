
run:
	go run main.go --env=local

infra:
	docker-compose up -d

infra_down:
	docker-compose down

build:
	docker-compose -f docker-compose.production.yml build

run_prod:
	docker-compose -f docker-compose.production.yml up -d

infra_down_prod:
	docker-compose -f docker-compose.production.yml down 