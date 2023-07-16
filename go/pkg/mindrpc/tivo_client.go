package mindrpc

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
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
	}, nil
}

func (t *TivoClient) Close() error {
	return t.connection.Close()
}

func (t *TivoClient) Authorize(ctx context.Context) error {
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
		return ErrUnauthorized(authResponse.Body.Message)
	}

	return nil
}

func (t *TivoClient) GetRecordings(ctx context.Context) ([]*model.Show, error) {

	return nil, nil
}

func (t *TivoClient) Send(ctx context.Context, tivoMessage TivoMessage) error {

	tivoRequestMessage := tivoMessage.
		WithSessionID(t.sessionID).
		WithRpcID(t.rpcID)

	if ce := logz.Logger.Check(zap.DebugLevel, "debugging"); ce != nil {
		logz.Logger.Info("sending mRPC message:")
		tivoRequestMessage.WriteTo(os.Stdout)
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
	tivoMessage := NewTivoMessage()

	t.connection.SetDeadline(time.Now().Add(config.Values.Timeout))
	defer t.connection.SetDeadline(time.Time{})
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
