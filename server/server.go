package server

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/server/config"
	"github.com/magrandera/PiaaS/server/models"
	"github.com/magrandera/PiaaS/server/routing"
)

// StartServer : start the PiaaS server
func StartServer() {
	values, err := config.ReadConfig("config.json")
	var port *string

	if err != nil {
		port = flag.String("port", "", "IP address")
		flag.Parse()

		//User is expected to give :8080 like input, if they give 8080
		//we'll append the required ':'
		if !strings.HasPrefix(*port, ":") {
			*port = ":" + *port
			log.Println("port is " + *port)
		}

		values.ServerPort = *port
	}
	models.InitDB()
	router := mux.NewRouter()
	routing.SetupRouting(router)
	log.Fatal(http.ListenAndServe(values.ServerPort, router))
}
