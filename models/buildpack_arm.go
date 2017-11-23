// +build linux,arm

package models

// Buildpack : stores the different buildpacks
type Buildpack struct {
	Name string
}

// InitBuildpacks : Initializes the database for buildpacks
func InitBuildpacks() {
	records, _ := db.ReadAll("buildpack")
	if len(records) != 0 {
		return
	}
	if err := db.Write("buildpack", "python3", Buildpack{Name: "arm32v6/python:3"}); err != nil {
		return
	}
	if err := db.Write("buildpack", "python2", Buildpack{Name: "arm32v6/python:2"}); err != nil {
		return
	}
}
