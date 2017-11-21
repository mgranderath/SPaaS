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

func DeployApplication(name string) {
	models.DeployApplication(os.Stdout, name)
}
