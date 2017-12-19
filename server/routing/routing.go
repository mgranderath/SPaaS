package routing

import (
	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/server/views"
)

// SetupRouting : sets up the api urls
func SetupRouting(router *mux.Router) {
	router.HandleFunc("/", views.HomePage)

	router.HandleFunc("/api/v1/app", getApplications).Methods("GET")
	router.HandleFunc("/api/v1/app/{name}", getApplication).Methods("GET")
	router.HandleFunc("/api/v1/app/{name}", createApplication).Methods("POST")
	router.HandleFunc("/api/v1/app/{name}", deleteApplication).Methods("DELETE")
	router.HandleFunc("/api/v1/app/{name}/start", startApplication).Methods("POST")
	router.HandleFunc("/api/v1/app/{name}/stop", stopApplication).Methods("POST")
	// router.HandleFunc("/api/v1/app/{name}", updateApplication).Methods("PUT")

	router.HandleFunc("/api/v1/app/{name}/deploy", deployApplication).Methods("POST")
}
