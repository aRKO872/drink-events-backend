# drink-events-backend

## Local Setup : 
Reach out to admin for .env file. Later :
- Run `go mod tidy`
- Simply run `docker compose up -d`
- To stop all containers from running, use `docker compose down`

## For continuous file changes :
Use command `docker compose watch`, and run on another terminal.

## To create migrations : 
migrate create -ext sql -dir internal/migrations -seq <migration_name>