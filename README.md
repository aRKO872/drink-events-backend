# drink-events-backend

## For continuos file changes :
Use command `docker compose watch`, and run on another terminal.

## Adding Migrations :
Use command : 
goose -dir "./migration/path" create migration-name sql

For Example: 
goose -dir "./internal/migrations/" create add_users_table sql