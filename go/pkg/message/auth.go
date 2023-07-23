package message

type Credential struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type AuthResponseBody struct {
	Type       Type        `json:"type,omitempty"`
	BodyID     string      `json:"bodyId,omitempty"`
	Message    string      `json:"message,omitempty"`
	Status     StatusType  `json:"status,omitempty"`
	Credential *Credential `json:"credential,omitempty"`
}
