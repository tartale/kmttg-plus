package mindrpc

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path"
	"time"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/test"
	"go.uber.org/zap"
)

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
	tsn        string
}

func NewTivoClient(tivo *model.Tivo) (*TivoClient, error) {

	tlsConfig, err := newTLSConfig(tivo)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, config.Values.Timeout)
	dialer := tls.Dialer{
		NetDialer: new(net.Dialer),
		Config:    tlsConfig,
	}

	conn, err := dialer.DialContext(ctx, "tcp", tivo.Address+":"+tivoRPCPort)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	sessionID := fmt.Sprintf("0x%x", rand.Int())

	return &TivoClient{
		connection: conn.(*tls.Conn),
		sessionID:  sessionID,
		rpcID:      0,
		tsn:        tivo.Tsn,
	}, nil
}

func (t *TivoClient) Close() error {
	return t.connection.Close()
}

func (t *TivoClient) Authenticate(ctx context.Context) error {

	authRequest := NewTivoMessage().WithAuthRequest(config.Values.MediaAccessKey)
	err := t.Send(ctx, *authRequest)
	if err != nil {
		return err
	}

	authResponse, err := t.Receive(context.Background())
	if err != nil {
		return err
	}
	if authResponse.Body.Status != success {
		return ErrNotAuthenticated(authResponse.Body.Message)
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

	bodyConfigRequest := NewTivoMessage().WithBodyConfigSearch()
	err = t.Send(ctx, *bodyConfigRequest)
	if err != nil {
		return err
	}

	_, err = t.Receive(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (t *TivoClient) GetRecordings(ctx context.Context) ([]*model.Show, error) {

	getRecordingsRequest := NewTivoMessage().WithGetRecordingsRequest()
	err := t.Send(ctx, *getRecordingsRequest)
	if err != nil {
		return nil, err
	}

	_, err = t.Receive(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *TivoClient) Send(ctx context.Context, tivoMessage TivoMessage) error {

	tivoRequestMessage := tivoMessage.
		WithSessionID(t.sessionID).
		WithRpcID(t.rpcID)

	if _, ok := tivoRequestMessage.Headers["BodyId"]; !ok && tivoRequestMessage.Body.BodyID == "" {
		tivoRequestMessage = tivoRequestMessage.WithBodyId("tsn:" + t.tsn)
	}

	if logz.Logger.Level() == zap.DebugLevel {
		debugDir, err := test.DebugDir()
		if err == nil {
			debugFile, err := os.Create(path.Join(debugDir, fmt.Sprintf("request-%d.txt", t.rpcID)))
			if err == nil {
				tivoRequestMessage.WriteTo(debugFile)
			}
		}
	}

	t.connection.SetDeadline(time.Now().Add(config.Values.Timeout))
	defer t.connection.SetDeadline(time.Time{})
	_, err := tivoRequestMessage.WriteTo(t.connection)
	if err != nil {
		return err
	}

	t.rpcID++

	return nil
}

func (t *TivoClient) Receive(ctx context.Context) (*TivoMessage, error) {

	responseReader := bufio.NewReader(t.connection)
	tivoResponseMessage := NewTivoMessage()

	t.connection.SetDeadline(time.Now().Add(config.Values.Timeout))
	defer t.connection.SetDeadline(time.Time{})
	_, err := tivoResponseMessage.ReadFrom(responseReader)
	if err != nil {
		return nil, err
	}

	if logz.Logger.Level() == zap.DebugLevel {
		debugDir, err := test.DebugDir()
		if err == nil {
			rpcId := tivoResponseMessage.Headers["RpcId"]
			debugFile, err := os.Create(path.Join(debugDir, fmt.Sprintf("response-%s.txt", rpcId)))
			if err == nil {
				tivoResponseMessage.WriteTo(debugFile)
			}
		}
	}

	return tivoResponseMessage, nil
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
