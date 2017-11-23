package command

import (
	"fmt"
	"os"

	"github.com/magrandera/PiaaS/models"
)

// ListApplications : Print the applications
func ListApplications() {
	applications, err := models.GetApplications()
	if err != nil {
		fmt.Println(err)
	}
	for _, name := range applications {
		fmt.Println(name.Name)
	}
}

// CreateApplication : Create a new Application
func CreateApplication(name string) {
	models.CreateApplication(os.Stdout, name)
}

// DeleteApplication : Delete existing Application
func DeleteApplication(name string) {
	models.DeleteApplication(os.Stdout, name)
}

// DeployApplication : Deploys the Application
func DeployApplication(name string) {
	models.DeployApplication(os.Stdout, name)
}

// StartApplication : starts the Application
func StartApplication(name string) {
	models.StartApplication(os.Stdout, name)
}

// StopApplication : stops the Application
func StopApplication(name string) {
	models.StopApplication(os.Stdout, name)
}
