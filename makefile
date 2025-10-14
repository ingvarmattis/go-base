local-deps-up:
	docker compose -f ./build/local/docker-compose.yaml up -d

local-deps-down:
	docker compose -f ./build/local/docker-compose.yaml down

create-migration:
	migrate create -ext sql -dir ./build/migrations <migration-name>

migration-up:
	migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path ./build/migrations up

migration-down:
	migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path ./build/migrations down
