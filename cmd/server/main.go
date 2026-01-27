package main

import (
	"fmt"
	"net/http"

	"github.com/rahulp18/bookeasy-backend/internal/db"
	"github.com/rahulp18/bookeasy-backend/internal/routes"
)

func main() {
	db.Connect()

	mux := http.NewServeMux()

	routes.Register(mux)
	fmt.Println("Server is listing on port 4200")
	http.ListenAndServe(":4200", mux)

}
