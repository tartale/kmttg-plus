package message

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/tartale/go/pkg/jsontime"
	"github.com/tartale/go/pkg/primitive"
	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
)

const (
	crlf               = "\r\n"
	schemaVersion      = "17"
	applicationName    = "Quicksilver"
	applicationVersion = "1.2"
)

type TivoMessage struct {
	Headers TivoMessageHeaders
	Body    any
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

func (t *TivoMessage) WithBody(body any) *TivoMessage {
	t.Body = body

	return t
}

func (t *TivoMessage) WithAuthRequest(mediaAccessKey string) *TivoMessage {

	t = t.WithStandardHeaders()
	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(TypeBodyAuthenticate))
	t.Headers.Set("ResponseCount", string(ResponseCountSingle))
	t.Headers.Set("BodyId", "")

	body := &AuthResponseBody{
		Type: TypeBodyAuthenticate,
		Credential: &Credential{
			Type: string(CredentialTypeMak),
			Key:  mediaAccessKey,
		},
	}
	t.Body = body

	return t
}

func (t *TivoMessage) WithGetShowsRequest(ctx context.Context, bodyID string) *TivoMessage {

	t = t.WithStandardHeaders()
	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(TypeRecordingFolderItemSearch))
	t.Headers.Set("ResponseCount", string(ResponseCountSingle))
	t.Headers.Set("BodyId", bodyID)

	body := &RecordingFolderItemSearchRequestBody{
		Type:    TypeRecordingFolderItemSearch,
		BodyID:  bodyID,
		Offset:  primitive.Ref(apicontext.Offset(ctx)),
		Count:   primitive.Ref(apicontext.Limit(ctx)),
		Flatten: primitive.Ref(true),
	}
	t.Body = body

	return t
}

func (t *TivoMessage) WithGetRecordingListRequest(ctx context.Context, bodyID string) *TivoMessage {

	t = t.WithStandardHeaders()
	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(TypeRecordingFolderItemSearch))
	t.Headers.Set("ResponseCount", string(ResponseCountSingle))
	t.Headers.Set("BodyId", bodyID)

	body := &RecordingFolderItemSearchRequestBody{
		Type:    TypeRecordingFolderItemSearch,
		BodyID:  bodyID,
		Offset:  primitive.Ref(apicontext.Offset(ctx)),
		Count:   primitive.Ref(apicontext.Limit(ctx)),
		Flatten: primitive.Ref(true),
	}
	t.Body = body

	return t
}

func (t *TivoMessage) WithGetRecordingRequest(ctx context.Context, bodyID, recordingID string) *TivoMessage {

	t = t.WithStandardHeaders()
	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(TypeRecordingSearch))
	t.Headers.Set("ResponseCount", string(ResponseCountSingle))
	t.Headers.Set("BodyId", bodyID)

	body := &RecordingSearchRequestBody{
		Type:          TypeRecordingSearch,
		BodyID:        bodyID,
		LevelOfDetail: LevelOfDetailMedium,
		RecordingID:   recordingID,
	}
	t.Body = body

	return t
}

func (t *TivoMessage) WithGetCollectionRequest(ctx context.Context, collectionIDs []string) *TivoMessage {

	t = t.WithStandardHeaders()
	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(TypeCollectionSearch))
	t.Headers.Set("ResponseCount", string(ResponseCountSingle))

	body := &CollectionSearchRequestBody{
		Type:          TypeCollectionSearch,
		CollectionIDs: collectionIDs,
		LevelOfDetail: LevelOfDetailMedium,
	}
	t.Body = body

	return t
}

func (t *TivoMessage) WithGetEpisodesRequest(ctx context.Context, bodyID string) *TivoMessage {

	t = t.WithStandardHeaders()
	t.Headers.Set("Type", "request")
	t.Headers.Set("RequestType", string(TypeCollectionSearch))
	t.Headers.Set("ResponseCount", string(ResponseCountSingle))
	t.Headers.Set("BodyId", bodyID)

	body := &RecordingFolderItemSearchRequestBody{
		Type:    TypeRecordingFolderItemSearch,
		BodyID:  bodyID,
		Offset:  primitive.Ref(apicontext.Offset(ctx)),
		Count:   primitive.Ref(apicontext.Limit(ctx)),
		Flatten: primitive.Ref(true),
	}
	t.Body = body

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

	logz.DebugPayload(bodyBytes, (fmt.Sprintf("%03d-response-payload.json", t.Headers.RpcID())))
	err = jsontime.UnmarshalJSON(bodyBytes, &t.Body)
	if err != nil {
		return -1, err
	}

	return 0, nil
}

func (t *TivoMessage) WriteTo(w io.Writer) (n int64, err error) {

	headers := t.Headers.String()
	bodyBytes, err := jsontime.MarshalJSON(&t.Body)
	if err != nil {
		return -1, err
	}
	logz.DebugPayload(bodyBytes, (fmt.Sprintf("%03d-request-payload.json", t.Headers.RpcID())))
	body := string(bodyBytes) + "\n"
	preamble := fmt.Sprintf("MRPC/2 %d %d", len(headers), len(body))
	message := preamble + crlf + headers + body + "\n"

	_, err = w.Write([]byte(message))
	if err != nil {
		return -1, err
	}

	return 0, nil
}
