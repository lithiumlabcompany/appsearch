build:
	docker-compose build
down:
	docker-compose down
test: build down
	docker-compose run --rm app
tdd: build down
	docker-compose run --rm app air
