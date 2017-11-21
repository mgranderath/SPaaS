package models

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	sh "github.com/codeskyblue/go-sh"
	git "gopkg.in/src-d/go-git.v4"
)

// Application : stores information about the applications
type Application struct {
	Name       string `json:"name"`
	Running    bool   `json:"running"`
	Path       string `json:"path"`
	Repository string `json:"repo"`
	Type       string `json:"type"`
}

// FileExists returns if a path exists and is a file
func fileExists(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}

func tarit(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}

// CreateApplication : Creates a new Application
func CreateApplication(w *os.File, name string) (Application, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return Application{}, err
	}
	path := filepath.Join(dir, "Applications", name, "repo")
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintln(w, "----->ERROR: App already exists.")
		return Application{}, err
	}

	fmt.Fprintf(w, "----->Creating Application '%s'.\n", name)
	fmt.Fprintln(w, "----->Creating directories.")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return Application{}, err
	}
	fmt.Fprintln(w, "----->Success creating directories.")

	fmt.Fprintln(w, "----->Initializing git repository.")
	if _, err := git.PlainInit(path, true); err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return Application{}, err
	}
	err = os.MkdirAll(filepath.Join(path, "hooks"), os.ModePerm)
	if err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return Application{}, err
	}
	file, err := os.Create(filepath.Join(path, "hooks", "post-receive"))
	if err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return Application{}, err
	}
	defer file.Close()
	fmt.Fprintf(file, "#!/usr/bin/env bash\n%s/PiaaS app deploy %s\n", dir, name)
	err = os.Chmod(filepath.Join(path, "hooks", "post-receive"), 0755)
	if err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return Application{}, err
	}
	fmt.Fprintln(w, "----->Success initializing git repository.")
	fmt.Fprintf(w, "----->Repository path: %s\n", path)

	app := Application{}
	app.Name = name
	app.Path = filepath.Join(dir, "Applications", name)
	app.Repository = path
	fmt.Fprintln(w, "----->Creating database record.")
	if err := db.Write("app", name, app); err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return Application{}, err
	}
	fmt.Fprintln(w, "----->Success creating database record.")
	fmt.Fprintf(w, "----->Application '%s' successfully created.\n", name)
	return app, nil
}

// DeleteApplication : Deletes existing Application
func DeleteApplication(w *os.File, name string) (bool, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return false, err
	}
	appPath := filepath.Join("Applications", name)
	path := filepath.Join(dir, appPath)
	if _, err := os.Stat(path); err != nil {
		fmt.Fprintln(w, "----->ERROR: App does not exist.")
		return false, err
	}
	fmt.Fprintf(w, "----->Removing Application '%s'.\n", name)
	fmt.Fprintln(w, "----->Removing Directories.")
	if err = os.RemoveAll(path); err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return false, err
	}
	fmt.Fprintln(w, "----->Success removing directories.")
	fmt.Fprintln(w, "----->Deleting Database Record.")
	if err := db.Delete("app", name); err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return false, err
	}
	fmt.Fprintln(w, "----->Success deleting database record.")
	fmt.Fprintf(w, "----->Application '%s' successfully removed.\n", name)
	return true, nil
}

// GetApplications : Get a list of all applications
func GetApplications() ([]Application, error) {
	records, err := db.ReadAll("app")
	if err != nil {
		fmt.Println("Error", err)
		return []Application{}, err
	}
	applications := []Application{}
	for _, f := range records {
		appFound := Application{}
		if err := json.Unmarshal([]byte(f), &appFound); err != nil {
			fmt.Println("Error", err)
			return []Application{}, err
		}
		applications = append(applications, appFound)
	}
	return applications, nil
}

// UpdateApplication : Update the state of an application
func UpdateApplication(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement update function
}

// DeployApplication : Deploys the application
func DeployApplication(w *os.File, name string) {
	fmt.Fprintf(w, "----->Deploying %s.\n", name)
	app, err := GetApplication(name)
	if err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return
	}
	path := filepath.Join(app.Path, "deploy")
	if err = os.RemoveAll(path); err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return
	}
	fmt.Fprintln(w, "----->Creating seperate directory for deployment.")
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL: app.Repository,
	})
	if err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
		return
	}
	fmt.Fprintln(w, "----->Success.")
	if fileExists(filepath.Join(path, "requirements.txt")) {
		fmt.Fprintln(w, "----->Python was detected.")
		app.Type = "python"
		if err := db.Write("app", name, app); err != nil {
			fmt.Fprintln(w, "----->ERROR:", err)
			return
		}
	} else {
		fmt.Fprintln(w, "----->No type detected.")
	}
	err = CreateDockerfile(app)
	if err != nil {
		fmt.Fprintln(w, "----->ERROR:", err)
	}
	session := sh.NewSession()
	session.SetDir(path + "/")
	session.Command("docker", "build", "-t", name, ".").Run()
	session.Command("docker", "run", "-d", "--rm", "--name", name, name).Run()
	session.ShowCMD = true
}

// GetApplication : Get specific application
func GetApplication(name string) (Application, error) {
	app := Application{}
	if err := db.Read("app", name, &app); err != nil {
		fmt.Println("Error", err)
		return Application{}, err
	}
	return app, nil
}
