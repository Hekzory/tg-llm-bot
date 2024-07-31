.PHONY: 
		build up upd down restart reup reupd

build:
		docker-compose build

up:
		docker-compose up

upd:
		docker-compose up -d

down:
		docker-compose down

restart:
		docker-compose stop
		docker-compose up -d

reup:
		docker-compose down
		docker-compose build
		docker-compose up

reupd:
		docker-compose down
		docker-compose build
		docker-compose up -d
