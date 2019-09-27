package handlers

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/config"
	"github.com/mgranderath/SPaaS/server/docker"
	"github.com/mgranderath/SPaaS/server/model"
	"github.com/mgranderath/SPaaS/server/services"
	"net/http"
)

type ServiceProvider interface {
	CreateApplication(c echo.Context) error
	DeleteApplication(c echo.Context) error
	DeployApplication(c echo.Context) error
	StartApplication(c echo.Context) error
	StopApplication(c echo.Context) error
	GetApplication(c echo.Context) error
	GetApplications(c echo.Context) error
	GetApplicationLogs(c echo.Context) error
	ChangePassword(c echo.Context) error
	Authorize(c echo.Context) error
}

type serviceProvider struct {
	appService  *services.AppService
	authService *services.AuthService
}

func NewServiceProvider(configRepository *config.Store, dockerClient *docker.Docker) ServiceProvider {
	return &serviceProvider{
		services.NewAppService(configRepository, dockerClient),
		services.NewAuthService(configRepository),
	}
}

// CreateApplication creates a new application
func (provider *serviceProvider) CreateApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	log.Infof("application '%s' is being created", name)
	messages := make(chan model.Status)
	go provider.appService.Create(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' creation failed with: %v", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}

// DeleteApplication deletes the application
func (provider *serviceProvider) DeleteApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	log.Infof("application '%s' is being deleted", name)
	messages := make(chan model.Status)
	go provider.appService.Delete(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' deletion failed with: %v", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}

// DeployApplication deploys an application
func (provider *serviceProvider) DeployApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	log.Infof("application '%s' is being deployed", name)
	messages := make(chan model.Status)
	go provider.appService.Deploy(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' deployment failed with: %v", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}

// StartApplication starts an application
func (provider *serviceProvider) StartApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	log.Infof("application '%s' is being started", name)
	messages := make(chan model.Status)
	go provider.appService.Start(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' start failed with: %v", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}

// StopApplication starts an application
func (provider *serviceProvider) StopApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	log.Infof("application '%s' is being stopped", name)
	messages := make(chan model.Status)
	go provider.appService.Stop(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' stop failed with: %v", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}

// GetApplication returns a current application
func (provider *serviceProvider) GetApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	container, err := provider.appService.GetApplicationStats(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Status{
			Type:    "error",
			Message: err.Error(),
		})
	}
	if err := common.EncodeJSONAndFlush(c, container); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Status{
			Type:    "error",
			Message: err.Error(),
		})
	}
	return nil
}

// GetApplications returns a list of all applications
func (provider *serviceProvider) GetApplications(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	files, err := provider.appService.GetApplications()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Status{
			Type:    "error",
			Message: err.Error(),
		})
	}
	applications := make([]map[string]string, len(files))
	for index, f := range files {
		app := model.NewApplication(f)
		appType := app.DetectType()
		applications[index] = map[string]string{
			"name": app.Name,
			"type": appType.ToString(),
		}
	}
	if err := common.EncodeJSONAndFlush(c, applications); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Status{
			Type:    "error",
			Message: err.Error(),
		})
	}
	return nil
}

func (provider *serviceProvider) GetApplicationLogs(c echo.Context) error {
	name := c.Param("name")
	messages := make(chan model.Status)
	go provider.appService.GetLogsStream(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' logs failed with: %v", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}
