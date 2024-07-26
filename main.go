package main

import (
	"fmt"
	"net/http"

	"github.com/drink-events-backend/cmd/routers"
	internal_database "github.com/drink-events-backend/internal"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	if setupErr := internal_database.SetupMigrations(); setupErr != nil {
		fmt.Println("error setting up migrations : ", setupErr.Error())
		return
	}

	r := routers.InitRouter()

	server := http.Server{
		Addr: ":3050",
		Handler: r,
	}
	
	server.ListenAndServe()
}