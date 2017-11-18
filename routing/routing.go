package routing

import (
	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/app"
	"github.com/magrandera/PiaaS/views"
)

// SetupRouting : sets up the api urls
func SetupRouting(router *mux.Router) {
	router.HandleFunc("/", views.HomePage)

	router.HandleFunc("/api/app", app.GetApplications).Methods("GET")
	router.HandleFunc("/api/app/{name}", app.GetApplication).Methods("GET")
	router.HandleFunc("/api/app/{name}", app.CreateApplication).Methods("POST")
	router.HandleFunc("/api/app/{name}", app.DeleteApplication).Methods("DELETE")
	router.HandleFunc("/api/app/{name}", app.DeleteApplication).Methods("PUT")
}
