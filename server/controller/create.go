package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/config"
	"github.com/mgranderath/SPaaS/server/auth"
	"github.com/mgranderath/SPaaS/server/hook"
	git "gopkg.in/src-d/go-git.v4"
)

func create(name string, messages chan<- Application) {
	appPath := filepath.Join(basePath, "applications", name)
	repoPath := filepath.Join(appPath, "repo")
	externalRepoPath := filepath.Join(config.Cfg.Config.GetString("HOST_CONFIG_FOLDER"), "applications", name, "repo")
	// Check if app already exists
	if common.Exists(appPath) {
		messages <- Application{
			Type:    "error",
			Message: "Already exists",
		}
		close(messages)
		return
	}
	// Create Directories
	messages <- Application{
		Type:    "info",
		Message: "Creating directories",
	}
	err := os.MkdirAll(repoPath, os.ModePerm)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Creating directories",
	}
	// Initialize the git repository
	messages <- Application{
		Type:    "info",
		Message: "Creating git repo",
	}
	if _, err := git.PlainInit(repoPath, true); err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Creating git repo",
	}
	// Create git post-receive hook
	messages <- Application{
		Type:    "info",
		Message: "Creating git receive hook",
	}
	err = os.MkdirAll(filepath.Join(repoPath, "hooks"), os.ModePerm)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	file, err := os.Create(filepath.Join(repoPath, "hooks", "post-receive"))
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	defer file.Close()
	token, err := auth.GetToken()
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	prefix := "http://"
	if config.Cfg.Config.GetBool("letsencrypt") {
		prefix = "https://"
	}
	postReceive, err := hook.CreatePostReceive(name, token, "spaas."+config.Cfg.Config.GetString("domain"), prefix)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	_, err = file.WriteString(postReceive)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	// Make the hook executable
	err = os.Chmod(filepath.Join(repoPath, "hooks", "post-receive"), 0755)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Creating git receive hook",
	}
	messages <- Application{
		Type:    "success",
		Message: "Creating app",
		Extended: []KeyValue{
			{Key: "RepoPath", Value: externalRepoPath},
		},
	}
	close(messages)
}

// CreateApplication creates a new application
func CreateApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	messages := make(chan Application)
	go create(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			return c.JSON(http.StatusInternalServerError, Application{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}
