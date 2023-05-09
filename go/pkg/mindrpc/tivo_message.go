package mindrpc

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

const (
	crlf = "\r\n"
)

type TivoMessageHeaders map[string]string

func (t TivoMessageHeaders) Set(key, val string) {
	t[key] = val
}

func (t TivoMessageHeaders) String() string {
	keys := maps.Keys(t)
	sort.Strings(keys)
	sb := strings.Builder{}

	for _, key := range keys {
		sb.WriteString(fmt.Sprintf("%s: %s%s", key, t[key], crlf))
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

func (t *TivoMessage) WithAuthPayload(mediaAccessKey string) *TivoMessage {
	bodyID := ""
	t.Headers.Set("BodyId", bodyID)

	requestType := bodyAuthenticate
	t.Headers.Set("RequestType", string(requestType))
	t.Payload.Type = requestType

	t.Payload.Credential.Type = string(makCredential)
	t.Payload.Credential.Key = mediaAccessKey

	t.Headers.Set("ResponseCount", string(single))

	return t
}

func (t *TivoMessage) AsRequest() *TivoMessage {
	t.Headers.Set("Type", "request")
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

	return strings.Join([]string{preamble, headers, payloadJSON}, crlf), nil
}

type RequestType string

const (
	bodyAuthenticate          RequestType = "bodyAuthenticate"
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

type ResultType string

const (
	channel             ResultType = "channel"
	recordingFolderItem ResultType = "recordingFolderItem"
	recording           ResultType = "recording"
	offer               ResultType = "offer"
	content             ResultType = "content"
	collection          ResultType = "collection"
	category            ResultType = "category"
	whatsOn             ResultType = "whatsOn"
	state               ResultType = "state"
)

type ResponseCount string

const (
	single   ResponseCount = "single"
	multiple ResponseCount = "multiple"
)

type CredentialType string

const (
	makCredential CredentialType = "makCredential"
)

type Payload struct {
	Type       RequestType `json:"type"`
	BodyID     string      `json:"bodyId,omitempty"`
	Credential struct {
		Type string `json:"type"`
		Key  string `json:"key"`
	} `json:"credential,omitempty"`
	Offset int `json:"offset,omitempty"`
	Count  int `json:"count,omitempty"`
}
