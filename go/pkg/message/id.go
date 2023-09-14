package message

type IdSearchRequestBody struct {
	Type      Type        `json:"type,omitempty"`
	BodyID    string      `json:"bodyId,omitempty"`
	Namespace IdNamespace `json:"namespace,omitempty"`
	ObjectID  string      `json:"objectId,omitempty"`
}

// {"objectId": ["mfs:rc.142248"], "type": "idSet"}
type IdSearchResponseBody struct {
	Type     Type     `json:"type,omitempty"`
	ObjectID []string `json:"objectId,omitempty"`
}
