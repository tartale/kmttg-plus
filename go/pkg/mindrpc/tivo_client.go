package mindrpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"go.uber.org/zap"
	"golang.org/x/crypto/pkcs12"
)

const (
	tivoRPCPort         = "1413"
	certificatePassword = "vlZaKoduom"
	schemaVersion       = "17"
	applicationName     = "Quicksilver"
	applicationVersion  = "1.2"
)

type TivoClient struct {
	connection *tls.Conn
	sessionID  string
	rpcID      int
}

func NewTivoClient(address string) (*TivoClient, error) {

	tlsConfig, err := tlsConfigFromCertificates()
	if err != nil {
		return nil, err
	}

	conn, err := tls.Dial("tcp", address+":"+tivoRPCPort, tlsConfig)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	sessionID := fmt.Sprintf("0x%x", rand.Int())

	return &TivoClient{
		connection: conn,
		sessionID:  sessionID,
		rpcID:      1,
	}, nil
}

func (t *TivoClient) Close() error {
	return t.connection.Close()
}

func (t *TivoClient) SendRequest(tivoMessage TivoMessage) error {

	tivoRequestMessage := tivoMessage.
		AsRequest().
		WithSessionID(t.sessionID).
		WithRpcID(t.rpcID)

	mindRpcMessage, err := tivoRequestMessage.ToMindRpcMessage()
	if err != nil {
		return err
	}

	if logz.Logger.Level() == zap.DebugLevel {
		logz.Logger.Info("sending mRpc message:")
		fmt.Println(mindRpcMessage)
	}
	_, err = t.connection.Write([]byte(mindRpcMessage))
	if err != nil {
		return err
	}

	t.rpcID++

	return nil
}

func (t *TivoClient) ReceiveResponse(ctx context.Context) {

	var buf bytes.Buffer
	io.Copy(&buf, t.connection)

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

func tlsConfigFromCertificates() (*tls.Config, error) {

	certificatePath, err := config.CertificatePath()
	if err != nil {
		return nil, err
	}
	certificateBytes, err := os.ReadFile(certificatePath)
	if err != nil {
		return nil, err
	}

	// the following is copied from this example:
	// https://pkg.go.dev/golang.org/x/crypto/pkcs12#example-ToPEM
	pemBlocks, err := pkcs12.ToPEM(certificateBytes, certificatePassword)
	if err != nil {
		return nil, err
	}

	var pemData []byte
	for _, b := range pemBlocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}, nil
}
