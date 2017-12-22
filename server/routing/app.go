package routing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/server/models"
)

func createApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app, err := models.CreateApplication(os.Stdout, params["name"])
	if err != nil {
		fmt.Fprint(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app, err := models.GetApplication(params["name"])
	if err != nil {
		fmt.Fprint(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getApplications(w http.ResponseWriter, r *http.Request) {
	apps, err := models.GetApplications()
	if err != nil {
		fmt.Fprint(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(apps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func deployApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app, err := models.DeployApplication(os.Stdout, params["name"])
	if err != nil {
		models.PrintErr(os.Stdout, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func startApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app, err := models.StartApplication(os.Stdout, params["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func stopApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app, err := models.StopApplication(os.Stdout, params["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func deleteApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	status, err := models.DeleteApplication(os.Stdout, params["name"])
	if err != nil {
		fmt.Fprint(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, status)
}
