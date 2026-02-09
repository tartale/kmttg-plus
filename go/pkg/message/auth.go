package message

type Credential struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type AuthResponseBody struct {
	Type       Type        `json:"type,omitempty"`
	BodyID     string      `json:"bodyId,omitempty"`
	Status     StatusType  `json:"status,omitempty"`
	Message    string      `json:"message,omitempty"`
	Credential *Credential `json:"credential,omitempty"`
}
