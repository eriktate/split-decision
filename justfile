set dotenv-load

run:
	air

migrate DIR:
	dbmate -d ./migrations {{DIR}}

build:
	go build -o splitdecision ./cmd/serve/main.go

connectdb:
	psql $DATABASE_URL
