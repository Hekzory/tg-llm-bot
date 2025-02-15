.PHONY: 
		build up upd down restart reup reupd

build:
		COMPOSE_BAKE=true docker-compose build

up:
		COMPOSE_BAKE=true docker-compose up

upd:
		COMPOSE_BAKE=true docker-compose up -d

down:
		COMPOSE_BAKE=true docker-compose down

restart:
		COMPOSE_BAKE=true docker-compose stop
		COMPOSE_BAKE=true docker-compose up -d

reup:
		docker-compose down
		COMPOSE_BAKE=true docker-compose build
		COMPOSE_BAKE=true docker-compose up

reupd:
		docker-compose down
		COMPOSE_BAKE=true docker-compose build
		COMPOSE_BAKE=true docker-compose up -d
