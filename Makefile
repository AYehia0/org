# the make file is used for the development stage only
compose-build:
	docker-compose -f docker-compose.dev.yaml build
compose-up:
	docker-compose -f docker-compose.dev.yaml up -d && docker-compose -f docker-compose.dev.yaml logs -f
compose-down:
	docker-compose -f docker-compose.dev.yaml down

.PHONY:
	compose-down compose-up
