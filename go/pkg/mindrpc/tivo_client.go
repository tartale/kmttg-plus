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
