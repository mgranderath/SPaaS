package routing

import (
	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/app"
	"github.com/magrandera/PiaaS/views"
)

// SetupRouting : sets up the api urls
func SetupRouting(router *mux.Router) {
	router.HandleFunc("/", views.HomePage)

	router.HandleFunc("/api/app/{name}", app.CreateApplication).Methods("POST")
	router.HandleFunc("/api/app/{name}", app.DeleteApplication).Methods("DELETE")
}
