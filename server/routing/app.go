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
	status := models.DeployApplication(os.Stdout, params["name"])
	fmt.Fprint(w, status)
}

func deleteApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	status, err := models.DeleteApplication(os.Stdout, params["name"])
	if err != nil {
		fmt.Fprint(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, status)
}
