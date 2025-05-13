ifneq (,$(wildcard ./.env))
	include .env
	export
	ENV_FILE_PARAM = --env-file .env
endif

start-db:
	docker compose -f docker-compose-db.yml up
stop-db:
	docker compose -f docker-compose-db.yml up

run:
	docker compose -f docker-compose.yml up

clean:
	docker compose down --volumes --remove-orphans


test:
	go test ./... -v

coverage:
	go test -coverprofile=coverage.out ./...

cover-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"



