package message

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

var headerOrder = []string{
	"Type", "RpcId", "SchemaVersion", "Content-Type", "RequestType",
	"ResponseCount", "BodyId", "X-ApplicationName", "X-ApplicationVersion",
	"X-ApplicationSessionId", "IsFinal",
}

type TivoMessageHeaders map[string]string

func NewTivoMessageHeaders(r io.Reader, headerLength int) (TivoMessageHeaders, error) {

	result := make(map[string]string)
	headerBuffer := make([]byte, headerLength)
	_, err := io.ReadFull(r, headerBuffer)
	if err != nil {
		return nil, err
	}
	header := string(headerBuffer)
	headerFields := strings.Split(header, crlf)
	for _, headerField := range headerFields {
		if headerField == "" {
			continue
		}
		headerKeyVal := strings.SplitN(headerField, ": ", 2)
		if len(headerKeyVal) == 1 {
			result[headerKeyVal[0]] = ""
		} else {
			result[headerKeyVal[0]] = headerKeyVal[1]
		}
	}

	return result, nil
}

func (t TivoMessageHeaders) Set(key, val string) {
	t[key] = val
}

func (t TivoMessageHeaders) RpcID() int {
	if rpcID, ok := t["RpcId"]; ok {
		val, err := strconv.ParseInt(rpcID, 10, 32)
		if err == nil {
			return int(val)
		}
	}

	return 0
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
