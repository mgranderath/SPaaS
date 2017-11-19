package routing

import (
	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/server/views"
)

// SetupRouting : sets up the api urls
func SetupRouting(router *mux.Router) {
	router.HandleFunc("/", views.HomePage)

	router.HandleFunc("/api/app", getApplications).Methods("GET")
	router.HandleFunc("/api/app/{name}", getApplication).Methods("GET")
	router.HandleFunc("/api/app/{name}", createApplication).Methods("POST")
	router.HandleFunc("/api/app/{name}", deleteApplication).Methods("DELETE")
	router.HandleFunc("/api/app/{name}", updateApplication).Methods("PUT")
}
