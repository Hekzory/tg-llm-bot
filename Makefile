.PHONY: 
		build up upd down restart reup reupd docs

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
		COMPOSE_BAKE=true docker-compose down
		COMPOSE_BAKE=true docker-compose build
		COMPOSE_BAKE=true docker-compose up

reupd:
		COMPOSE_BAKE=true docker-compose down
		COMPOSE_BAKE=true docker-compose build
		COMPOSE_BAKE=true docker-compose up -d

docs:
		@echo "Starting godoc server on http://localhost:6061"
		@cd go && go install golang.org/x/tools/cmd/godoc@latest && ~/go/bin/godoc -http=:6061

clean:
		COMPOSE_BAKE=true docker-compose down
		docker system prune -a