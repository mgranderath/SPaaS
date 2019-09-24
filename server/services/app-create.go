package services

import (
	"github.com/mgranderath/SPaaS/server/hook"
	"github.com/mgranderath/SPaaS/server/model"
	"github.com/pkg/errors"
	"os"
	"path/filepath"

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

func (a *AppService) Create(name string, messages model.StatusChannel) {
	app := model.NewApplication(name)
	externalRepoPath := filepath.Join(a.ConfigRespository.Config.GetString("HOST_CONFIG_FOLDER"), "applications", name, "repo")
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
