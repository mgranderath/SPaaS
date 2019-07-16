package config

// Dockerfile houses the information to generate the Dockerfile
type Dockerfile struct {
	Command []string
	Length  int
}
