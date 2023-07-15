package mindrpc

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	crlf = "\r\n"
)

var headerOrder = []string{"Type", "RpcId", "SchemaVersion", "Content-Type", "RequestType", "ResponseCount"}

type TivoMessageHeaders map[string]string

func (t TivoMessageHeaders) Set(key, val string) {
	t[key] = val
}

func (t TivoMessageHeaders) String() string {
	sb := strings.Builder{}

	for _, key := range headerOrder {
		if val, ok := t[key]; ok {
			sb.WriteString(fmt.Sprintf("%s: %s%s", key, val, crlf))
		}
	}
	sb.WriteString(crlf)

	return sb.String()
}

type TivoMessage struct {
	Headers TivoMessageHeaders
	Payload Payload
}

func NewTivoMessage() *TivoMessage {
	var tivoMessage = TivoMessage{
		Headers: make(map[string]string),
	}

	tivoMessage.Headers.Set("SchemaVersion", schemaVersion)
	tivoMessage.Headers.Set("Content-Type", "application/json")
	tivoMessage.Headers.Set("X-ApplicationName", applicationName)
	tivoMessage.Headers.Set("X-ApplicationVersion", applicationVersion)

	return &tivoMessage
}

func (t *TivoMessage) WithAuthRequest(mediaAccessKey string) *TivoMessage {

	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(bodyAuthenticate))
	t.Headers.Set("ResponseCount", string(single))

	t.Payload.Type = bodyAuthenticate
	t.Payload.Credential = &Credential{
		Type: string(makCredential),
		Key:  mediaAccessKey,
	}

	return t
}

func (t *TivoMessage) WithBodyConfigSearch() *TivoMessage {

	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(bodyConfigSearch))
	t.Headers.Set("ResponseCount", string(single))

	t.Payload.Type = bodyConfigSearch
	t.Payload.BodyID = "-"

	return t
}

func (t *TivoMessage) WithSessionID(sessionID string) *TivoMessage {
	t.Headers.Set("X-ApplicationSessionId", sessionID)

	return t
}

func (t *TivoMessage) WithRpcID(rpcID int) *TivoMessage {
	t.Headers.Set("RpcId", fmt.Sprintf("%d", rpcID))

	return t
}

func (t *TivoMessage) PayloadJSON() (string, error) {

	payloadJSON, err := json.Marshal(t.Payload)
	if err != nil {
		return "", err
	}

	return string(payloadJSON) + "\n", nil
}

func (t *TivoMessage) ToMindRpcMessage() (string, error) {

	headers := t.Headers.String()
	payloadJSON, err := t.PayloadJSON()
	if err != nil {
		return "", err
	}
	preamble := fmt.Sprintf("MRPC/2 %d %d", len(headers), len(payloadJSON))
	message := preamble + crlf + headers + payloadJSON

	return message, nil
}

type RequestType string

const (
	bodyAuthenticate          RequestType = "bodyAuthenticate"
	bodyConfigSearch          RequestType = "bodyConfigSearch"
	channelSearch             RequestType = "channelSearch"
	recordingFolderItemSearch RequestType = "recordingFolderItemSearch"
	recordingSearch           RequestType = "recordingSearch"
	offerSearch               RequestType = "offerSearch"
	contentSearch             RequestType = "contentSearch"
	collectionSearch          RequestType = "collectionSearch"
	categorySearch            RequestType = "categorySearch"
	whatsOnSearch             RequestType = "whatsOnSearch"
	tunerStateEventRegister   RequestType = "tunerStateEventRegister"
)

// type ResponseType string

// const (
// 	channel             ResponseType = "channel"
// 	recordingFolderItem ResponseType = "recordingFolderItem"
// 	recording           ResponseType = "recording"
// 	offer               ResponseType = "offer"
// 	content             ResponseType = "content"
// 	collection          ResponseType = "collection"
// 	category            ResponseType = "category"
// 	whatsOn             ResponseType = "whatsOn"
// 	state               ResponseType = "state"
// )

type ResponseCount string

const (
	single   ResponseCount = "single"
	multiple ResponseCount = "multiple"
)

type CredentialType string

const (
	makCredential CredentialType = "makCredential"
)

type Credential struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type Payload struct {
	Type       RequestType `json:"type"`
	BodyID     string      `json:"bodyId,omitempty"`
	Credential *Credential `json:"credential,omitempty"`
	Offset     int         `json:"offset,omitempty"`
	Count      int         `json:"count,omitempty"`
}
