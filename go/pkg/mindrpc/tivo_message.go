package mindrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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
	// t.Headers.Set("X-ApplicationName", applicationName)
	// t.Headers.Set("X-ApplicationVersion", applicationVersion)

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

/*
MRPC/2 132 41
Type: request
RpcId: 1
SchemaVersion: 21
Content-Type: application/json
RequestType: bodyConfigSearch
ResponseCount: single

{"type":"bodyConfigSearch","bodyId":"-"}
*/

func (t *TivoMessage) WithBodyConfigSearch() *TivoMessage {

	t = t.WithStandardHeaders()
	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(bodyConfigSearch))
	t.Headers.Set("ResponseCount", string(single))

	t.Body.Type = bodyConfigSearch
	t.Body.BodyID = "-"

	return t
}

/*
MRPC/2 170 68
Type: request
RpcId: 2
SchemaVersion: 21
Content-Type: application/json
RequestType: recordingFolderItemSearch
BodyId: tsn:8460001909E14AD
ResponseCount: single

{"type":"recordingFolderItemSearch","bodyId":"tsn:8460001909E14AD"}

*/

func (t *TivoMessage) WithGetRecordingsRequest() *TivoMessage {

	t = t.WithStandardHeaders()
	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(recordingFolderItemSearch))
	t.Headers.Set("ResponseCount", string(single))

	t.Body.Type = recordingFolderItemSearch

	return t
}

func (t *TivoMessage) WithSessionID(sessionID string) *TivoMessage {
	// t.Headers.Set("X-ApplicationSessionId", sessionID)

	return t
}

func (t *TivoMessage) WithRpcID(rpcID int) *TivoMessage {
	t.Headers.Set("RpcId", fmt.Sprintf("%d", rpcID))

	return t
}

func (t *TivoMessage) WithBodyId(bodyId string) *TivoMessage {
	t.Headers.Set("BodyId", bodyId)
	t.Body.BodyID = bodyId

	return t
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
	t.Headers = headers

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

func (t *TivoMessage) WriteTo(w io.Writer) (n int64, err error) {

	headers := t.Headers.String()
	bodyBytes, err := json.Marshal(t.Body)
	if err != nil {
		return -1, err
	}
	body := string(bodyBytes) + "\n"
	preamble := fmt.Sprintf("MRPC/2 %d %d", len(headers), len(body))
	message := preamble + crlf + headers + body

	_, err = w.Write([]byte(strings.ToValidUTF8(message, "")))
	if err != nil {
		return -1, err
	}

	return 0, nil
}

type MessageType string

const (
	bodyAuthenticate          MessageType = "bodyAuthenticate"
	bodyConfigSearch          MessageType = "bodyConfigSearch"
	recordingSearch           MessageType = "recordingSearch"
	recordingFolderItemSearch MessageType = "recordingFolderItemSearch"
	errorz                    MessageType = "error"
)

type StatusType string

const (
	success StatusType = "success"
	failure StatusType = "failure"
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

type Credential struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type TivoMessageBody struct {
	Type       MessageType `json:"type,omitempty"`
	BodyID     string      `json:"bodyId,omitempty"`
	Credential *Credential `json:"credential,omitempty"`
	Offset     int         `json:"offset,omitempty"`
	Count      int         `json:"count,omitempty"`
	Message    string      `json:"message,omitempty"`
	Status     StatusType  `json:"status,omitempty"`
}
