package main

import (
	"net/http"
	"github.com/drink-events-backend/cmd/routers"
)

func main() {
	r := routers.InitRouter()

	server := http.Server{
		Addr: ":3050",
		Handler: r,
	}
	
	server.ListenAndServe()
}