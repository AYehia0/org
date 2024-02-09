compose-build:
	docker-compose up --force-recreate --build -d && docker-compose logs -f
compose-up:
	docker-compose -f docker-compose.yaml up -d && docker-compose -f docker-compose.yaml logs -f
compose-down:
	docker-compose down

.PHONY:
	compose-down compose-up
