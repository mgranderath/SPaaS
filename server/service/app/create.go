package app

import (
	"github.com/labstack/gommon/log"
	"github.com/mgranderath/SPaaS/server/hook"
	"github.com/mgranderath/SPaaS/server/model"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/mgranderath/SPaaS/common"
	"gopkg.in/src-d/go-git.v4"
)

func createPostReceiveHook(name string, app *model.Application) error {
	hooksPath := filepath.Join(app.RepositoryPath, "hooks")
	err := os.MkdirAll(hooksPath, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(hooksPath, "post-receive"))
	if err != nil {
		return err
	}
	defer file.Close()
	postReceive, err := hook.GetPostReceiveHookHelperString(name)
	if err != nil {
		return err
	}
	_, err = file.WriteString(postReceive)
	if err != nil {
		return err
	}
	written, err := common.Copy("./util/post-receive", filepath.Join(hooksPath, "post-receive-deploy"))
	if err != nil || written == 0 {
		return err
	}
	// Make the hook executable
	err = os.Chmod(filepath.Join(hooksPath, "post-receive"), 0755)
	if err != nil {
		return err
	}
	err = os.Chmod(filepath.Join(hooksPath, "post-receive-deploy"), 0755)
	if err != nil {
		return err
	}
	return nil
}

func (appService *AppService) create(name string, messages model.StatusChannel) {
	app := model.NewApplication(name)
	externalRepoPath := filepath.Join(appService.Config.Config.GetString("HOST_CONFIG_FOLDER"), "applications", name, "repo")
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
	err = createPostReceiveHook(name, app)
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
func (app *AppService) CreateApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	log.Infof("application '%s' is being created", name)
	messages := make(chan model.Status)
	go app.create(name, messages)
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
