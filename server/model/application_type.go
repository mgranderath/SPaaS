package model

type ApplicationType int

const (
	Undefined ApplicationType = iota
	Node      ApplicationType = iota
	Ruby      ApplicationType = iota
	Python    ApplicationType = iota
)

func (appType ApplicationType) ToString() string {
	return [...]string{"node", "ruby", "python"}[appType]
}
