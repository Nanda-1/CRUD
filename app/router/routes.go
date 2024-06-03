package router

import (
	"CRUD/app/controllers"
	"CRUD/app/middleware"
	"database/sql"

	"github.com/gorilla/mux"
)

// Router sets up the routes for the application
func Router(database *sql.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/post", middleware.Authentication(controllers.CreatePostHandler(database))).Methods("POST")
	router.HandleFunc("/api/post/{post_id}", middleware.Authentication(controllers.GetPostHandler(database))).Methods("GET")
	router.HandleFunc("/api/post/{post_id}", middleware.Authentication(controllers.UpdatePostHandler(database))).Methods("PUT")
	router.HandleFunc("/api/post/{post_id}", middleware.Authentication(controllers.DeletePostHandler(database))).Methods("DELETE")
	return router
}
