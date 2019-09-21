package controller

import (
	"github.com/labstack/gommon/log"
	"github.com/mgranderath/SPaaS/server/model"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/config"
	"github.com/mgranderath/SPaaS/server/auth"
	"github.com/mgranderath/SPaaS/server/hook"
	"gopkg.in/src-d/go-git.v4"
)

func create(name string, messages model.StatusChannel) {
	app := model.NewApplication(name)
	externalRepoPath := filepath.Join(config.Cfg.Config.GetString("HOST_CONFIG_FOLDER"), "applications", name, "repo")
	// Check if app already exists
	if app.Exists() {
		messages.SendError(errors.New("Already exists"))
		close(messages)
		return
	}
	// Create Directories
	messages.SendInfo("Creating directories")
	err := os.MkdirAll(app.RepositoryPath, os.ModePerm)
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Creating directories")
	// Initialize the git repository
	messages.SendInfo("Creating repository")
	if _, err := git.PlainInit(app.RepositoryPath, true); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Creating repository")
	// Create git post-receive hook
	messages.SendInfo("Creating receive hook")
	err = os.MkdirAll(filepath.Join(app.RepositoryPath, "hooks"), os.ModePerm)
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	file, err := os.Create(filepath.Join(app.RepositoryPath, "hooks", "post-receive"))
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	defer file.Close()
	token, err := auth.GetToken()
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	prefix := "http://"
	if config.Cfg.Config.GetBool("letsencrypt") {
		prefix = "https://"
	}
	postReceive, err := hook.CreatePostReceive(name, token, "spaas."+config.Cfg.Config.GetString("domain"), prefix)
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	_, err = file.WriteString(postReceive)
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	written, err := common.Copy("./util/post-receive", filepath.Join(app.RepositoryPath, "hooks", "post-receive-deploy"))
	if err != nil || written == 0 {
		messages.SendError(err)
		close(messages)
		return
	}
	// Make the hook executable
	err = os.Chmod(filepath.Join(app.RepositoryPath, "hooks", "post-receive"), 0755)
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	err = os.Chmod(filepath.Join(app.RepositoryPath, "hooks", "post-receive-deploy"), 0755)
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Creating git receive hook")
	messages <- model.Status{
		Type:    "success",
		Message: "Creating app",
		Extended: []model.KeyValue{
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
	log.Infof("application '%s' is being created\n", name)
	messages := make(chan model.Status)
	go create(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' creation failed with: %v\n", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}
