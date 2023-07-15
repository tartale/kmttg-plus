package mindrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

const (
	crlf = "\r\n"
)

type TivoMessage struct {
	Headers TivoMessageHeaders
	Body    TivoMessageBody
}

func NewTivoMessage() *TivoMessage {
	var tivoMessage = TivoMessage{
		Headers: make(map[string]string),
	}

	return &tivoMessage
}

func (t *TivoMessage) WithStandardHeaders() *TivoMessage {
	t.Headers.Set("SchemaVersion", schemaVersion)
	t.Headers.Set("Content-Type", "application/json")
	t.Headers.Set("X-ApplicationName", applicationName)
	t.Headers.Set("X-ApplicationVersion", applicationVersion)

	return t
}

func (t *TivoMessage) WithAuthRequest(mediaAccessKey string) *TivoMessage {

	t = t.WithStandardHeaders()
	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(bodyAuthenticate))
	t.Headers.Set("ResponseCount", string(single))

	t.Body.Type = bodyAuthenticate
	t.Body.Credential = &Credential{
		Type: string(makCredential),
		Key:  mediaAccessKey,
	}

	return t
}

func (t *TivoMessage) WithBodyConfigSearch() *TivoMessage {

	t = t.WithStandardHeaders()
	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(bodyConfigSearch))
	t.Headers.Set("ResponseCount", string(single))

	t.Body.Type = bodyConfigSearch
	t.Body.BodyID = "-"

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

	payloadJSON, err := json.Marshal(t.Body)
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

func (t *TivoMessage) ReadFrom(r io.Reader) (n int64, err error) {
	responseReader := bufio.NewReader(r)

	preamble, err := responseReader.ReadString('\n')
	if err != nil {
		return -1, err
	}

	var headerLength, bodyLength int
	_, err = fmt.Sscanf(preamble, "MRPC/2 %d %d \n", &headerLength, &bodyLength)
	if err != nil {
		return -1, err
	}

	headers, err := NewTivoMessageHeaders(responseReader, headerLength)
	if err != nil {
		return -1, err
	}
	_ = headers

	bodyBytes := make([]byte, bodyLength)
	_, err = io.ReadFull(responseReader, bodyBytes)
	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(bodyBytes, &t.Body)
	if err != nil {
		return -1, err
	}

	return 0, nil
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

type StatusType string

const (
	success StatusType = "success"
	failure StatusType = "failure"
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

type TivoMessageBody struct {
	Type       RequestType `json:"type"`
	BodyID     string      `json:"bodyId,omitempty"`
	Credential *Credential `json:"credential,omitempty"`
	Offset     int         `json:"offset,omitempty"`
	Count      int         `json:"count,omitempty"`
	Message    string      `json:"message,omitempty"`
	Status     StatusType  `json:"status,omitempty"`
}
