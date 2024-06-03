package main

import (
	"CRUD/app/db"
	"CRUD/app/router"
	"log"
	"net/http"
)

func main() {
	database := db.Connect()
	defer database.Close()

	r := router.Router(database)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
