package model

// Status stores information about the application
type Status struct {
	Type     string     `json:"type"`
	Message  string     `json:"message"`
	Extended []KeyValue `json:"extended,omitempty"`
}

// KeyValue holds extra information of a model
type KeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type StatusChannel chan<- Status

func (channel StatusChannel) SendInfo(message string) {
	channel <- Status{
		Type:    "info",
		Message: message,
	}
}

func (channel StatusChannel) SendSuccess(message string) {
	channel <- Status{
		Type:    "success",
		Message: message,
	}
}

func (channel StatusChannel) SendError(err error) {
	channel <- Status{
		Type:    "error",
		Message: err.Error(),
	}
}
