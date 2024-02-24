APP_BINARY=accountApp

postgres:
	docker run --name postgres -p 5401:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres

createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres users

dropdb:
	docker exec -it postgres dropdb users

migrateup:
	migrate -path pkg/db/migrations -database "postgresql://postgres:password@localhost:5401/users?sslmode=disable" -verbose up

migrateup1:
	migrate -path pkg/db/migrations -database "postgresql://postgres:password@localhost:5401/users?sslmode=disable" -verbose up 1
	
migratedown:
	migrate -path pkg/db/migrations -database "postgresql://postgres:password@localhost:5401/users?sslmode=disable" -verbose down

migratedown1:
	migrate -path pkg/db/migrations -database "postgresql://postgres:password@localhost:5401/users?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v ${CURDIR}:/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

build:
	chdir . && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${APP_BINARY} .

mock:
	mockgen -package mockdb -destination pkg/db/mock/store.go github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc Store

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

.PHONY: postgres test sqlc createdb dropdb mock migratedown migrateup migratedown2 migrateup1 server proto build