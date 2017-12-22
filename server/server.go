package main

import (
	"log"
	"net/http"

	"github.com/magrandera/PiaaS/server/models"

	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/server/routing"
)

// StartServer : start the PiaaS server
func StartServer() {
	router := mux.NewRouter()
	routing.SetupRouting(router)
	models.InitDB()
	log.Fatal(http.ListenAndServe(":5000", router))
}
