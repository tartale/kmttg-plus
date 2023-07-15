package mindrpc

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
)

// https://github.com/lart2150/tivo-scripts/blob/master/lib/tivo.js

const (
	tivoRPCPort        = "1413"
	schemaVersion      = "21"
	applicationName    = "Quicksilver"
	applicationVersion = "1.2"
)

type TivoClient struct {
	connection *tls.Conn
	sessionID  string
	rpcID      int
}

func NewTivoClient(tivo *model.Tivo) (*TivoClient, error) {

	tlsConfig, err := newTLSConfig(tivo)
	if err != nil {
		return nil, err
	}

	conn, err := tls.Dial("tcp", tivo.Address+":"+tivoRPCPort, tlsConfig)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	sessionID := fmt.Sprintf("0x%x", rand.Int())

	return &TivoClient{
		connection: conn,
		sessionID:  sessionID,
		rpcID:      0,
	}, nil
}

func (t *TivoClient) Close() error {
	return t.connection.Close()
}

func (t *TivoClient) SendRequest(tivoMessage TivoMessage) error {

	tivoRequestMessage := tivoMessage.
		WithSessionID(t.sessionID).
		WithRpcID(t.rpcID)

	mindRpcMessage, err := tivoRequestMessage.ToMindRpcMessage()
	if err != nil {
		return err
	}

	if ce := logz.Logger.Check(zap.DebugLevel, "debugging"); ce != nil {
		logz.Logger.Info("sending mRPC message:")
		fmt.Print(mindRpcMessage)
	}
	_, err = t.connection.Write([]byte(strings.ToValidUTF8(mindRpcMessage, "")))
	if err != nil {
		return err
	}

	t.rpcID++

	return nil
}

func (t *TivoClient) ReceiveResponse(ctx context.Context) (*TivoMessage, error) {

	responseReader := bufio.NewReader(t.connection)
	tivoMessage := NewTivoMessage()
	_, err := tivoMessage.ReadFrom(responseReader)
	if err != nil {
		return nil, err
	}

	return tivoMessage, nil

	// buffer := bytes.NewBuffer([]byte{})
	// data := make([]byte, 4096)

	// _, err := m.socket.Read(data)
	// if err != nil {
	// 	fmt.Println("Error reading data: ", err)
	// 	return nil, nil
	// }

	// for bytes.IndexByte(buffer.Bytes(), '\n') < 0 {
	// 	_, err = m.socket.Read(data)
	// 	if err != nil {
	// 		fmt.Println("Error reading data: ", err)
	// 		return nil, nil
	// 	}
	// 	buffer.Write(data)
	// }

	// buf_val := buffer.String()
	// matches := m.proto_pat.FindStringSubmatch(buf_val)
	// h_size, _ := strconv.Atoi(matches[1])
	// b_size, _ := strconv.Atoi(matches[2])
	// h_start := matches[0][0] + len(matches[0])
	// if m.debug {
	// 	fmt.Printf("RPC Response (Offset: %d, H Size: %d, B Size: %d)\n", h_start, h_size, b_size)
	// 	fmt.Printf("RPC Response (Bytes Loaded: %d)\n", buffer.Len()-h_start)
	// }
	// for buffer.Len()-h_start < h_size+b_size {
	// 	_, err = m.socket.Read(data)
	// 	if err != nil {
	// 		fmt.Println("Error reading data: ", err)
	// 		return nil, nil
	// 	}
	// 	buffer.Write(data)
	// 	if m.debug {
	// 		fmt.Printf("RPC Response (Bytes Loaded: %d)\n", buffer.Len()-h_start)
	// 	}
	// }
	// buf_val = buffer.String()
	// headers := m.parse_headers(buf_val[h_start : h_start+h_size])
	// if m.debug {
	// 	fmt.Printf("RPC Response ID: %s\n", headers["RpcId"])
	// }
	// response_json := make(map[string]interface{})
	// err = json.Unmarshal([]byte(buf_val[h_start+h_size:]), &response_json)
	// if err != nil {
	// 	fmt.Println("Error parsing response JSON: ", err)
	// 	return nil, nil
	// }
	// return headers, response_json

}

func (t *TivoClient) WaitForResponse(rpcID string) {

}

func newTLSConfig(tivo *model.Tivo) (*tls.Config, error) {

	certPath, err := config.CertificatePath()
	if err != nil {
		return nil, err
	}
	cert, certPool, err := GetCertificates(certPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		GetClientCertificate: func(cri *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			logz.Logger.Debug("received client certificate request", zap.Any("cri", cri))
			return cert, nil
		},
		RootCAs:            certPool,
		ServerName:         tivo.ServerName(),
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          certPool,
		InsecureSkipVerify: true,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
	}, nil
}

// openssl s_client -crlf  -connect 10.0.1.18:1413 -cert ./cdata.pem -key ./cdata.pem \
//   -CAfile ./cdata.pem -cipher 'DEFAULT:@SECLEVEL=0' -purpose sslclient -servername 846-0001-909E-14AD
