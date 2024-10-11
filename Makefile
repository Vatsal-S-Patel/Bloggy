
createmigration:
	migrate create -ext sql -dir internal/migrations -seq ${name}

migrateup:
	migrate -path internal/migrations -database $(POSTGRES_URL_BLOGGY) -verbose up ${n}

migratedown:
	migrate -path internal/migrations -database $(POSTGRES_URL_BLOGGY) -verbose down ${n}

removedirtyread:
	migrate -path internal/migrations -database $(POSTGRES_URL_BLOGGY) force 1