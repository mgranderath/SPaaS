package routing

import (
	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/models"
	"github.com/magrandera/PiaaS/views"
)

// SetupRouting : sets up the api urls
func SetupRouting(router *mux.Router) {
	router.HandleFunc("/", views.HomePage)

	router.HandleFunc("/api/app", models.GetApplications).Methods("GET")
	router.HandleFunc("/api/app/{name}", models.GetApplication).Methods("GET")
	router.HandleFunc("/api/app/{name}", models.CreateApplication).Methods("POST")
	router.HandleFunc("/api/app/{name}", models.DeleteApplication).Methods("DELETE")
	router.HandleFunc("/api/app/{name}", models.UpdateApplication).Methods("PUT")
}
