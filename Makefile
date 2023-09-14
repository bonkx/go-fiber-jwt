build:
	docker-compose down
	docker-compose build
run:
	docker-compose up --remove-orphans
down:
	docker-compose down