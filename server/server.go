package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/server/routing"
	"github.com/takama/daemon"
)

// StartServer : start the PiaaS server
func StartServer() {
	service, err := daemon.New("PiaaS", "server")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	status, err := service.Install()
	if err != nil {
		log.Fatal(status, "\nError: ", err)
	}
	fmt.Println(status)
	router := mux.NewRouter()
	routing.SetupRouting(router)
	log.Fatal(http.ListenAndServe(":5000", router))
}
