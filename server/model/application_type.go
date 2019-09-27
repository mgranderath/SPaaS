package model

type ApplicationType int

const (
	Undefined ApplicationType = iota
	Node      ApplicationType = iota
	Ruby      ApplicationType = iota
	Python    ApplicationType = iota
	Docker    ApplicationType = iota
)

func (appType ApplicationType) ToString() string {
	return [...]string{"undefined", "node", "ruby", "python", "docker"}[appType]
}
