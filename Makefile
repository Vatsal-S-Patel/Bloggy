
run:
	go run ./cmd/

createmigration:
	migrate create -ext sql -dir internal/migrations -seq ${name}

migrateup:
	migrate -path internal/migrations -database $(POSTGRES_URL_BLOGGY) -verbose up ${n}

migratedown:
	migrate -path internal/migrations -database $(POSTGRES_URL_BLOGGY) -verbose down ${n}

migrateforce:
	migrate -path internal/migrations -database $(POSTGRES_URL_BLOGGY) force ${n}

lint:
	golangci-lint run