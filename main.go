package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/config"
	"github.com/magrandera/PiaaS/routing"
)

func main() {
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

	router := mux.NewRouter()
	routing.SetupRouting(router)
	log.Fatal(http.ListenAndServe(values.ServerPort, router))
}
