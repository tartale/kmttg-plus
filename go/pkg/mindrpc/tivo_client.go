package mindrpc

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"path"
	"time"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/test"
)

const tivoRPCPort = "1413"

type ReaderFromWriterTo interface {
	io.ReaderFrom
	io.WriterTo
}

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
	ctx, cancelFunc := context.WithTimeout(ctx, config.Values.Timeout)
	defer cancelFunc()
	dialer := tls.Dialer{
		NetDialer: new(net.Dialer),
		Config:    tlsConfig,
	}

	conn, err := dialer.DialContext(ctx, "tcp", tivo.Address+":"+tivoRPCPort)
	if err != nil {
		return nil, err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sessionID := fmt.Sprintf("0x%x", r.Int())[:9]

	return &TivoClient{
		connection: conn.(*tls.Conn),
		sessionID:  sessionID,
		rpcID:      1,
		tsn:        tivo.Tsn,
	}, nil
}

func (t *TivoClient) Close() error {
	return t.connection.Close()
}

func (t *TivoClient) Authenticate(ctx context.Context) error {

	authRequest := message.NewTivoMessage().WithAuthRequest(config.Values.MediaAccessKey)
	err := t.Send(ctx, authRequest)
	if err != nil {
		return err
	}

	authResponseBody := &message.AuthResponseBody{}
	authResponse := message.NewTivoMessage().WithBody(authResponseBody)
	err = t.Receive(context.Background(), authResponse)
	if err != nil {
		return err
	}
	if authResponseBody.Status != message.StatusTypeSuccess {
		return ErrNotAuthenticated(authResponseBody.Message)
	}

	return nil
}

func (t *TivoClient) GetRecordings(ctx context.Context) ([]*model.Show, error) {

	getRecordingsRequest := message.NewTivoMessage().WithGetRecordingsRequest("tsn:" + t.tsn)
	err := t.Send(ctx, getRecordingsRequest)
	if err != nil {
		return nil, err
	}

	getRecordingResponseBody := &message.RecordingFolderItemSearchResponseBody{}
	getRecordingResponse := message.NewTivoMessage().WithBody(getRecordingResponseBody)
	err = t.Receive(ctx, getRecordingResponse)
	if err != nil {
		return nil, err
	}
	_ = getRecordingResponse

	return nil, nil
}

func (t *TivoClient) Send(ctx context.Context, tivoMessage *message.TivoMessage) error {

	tivoRequestMessage := tivoMessage.
		WithSessionID(t.sessionID).
		WithRpcID(t.rpcID)

	test.Debug(tivoRequestMessage, (fmt.Sprintf("%03d-request.txt", t.rpcID)))

	t.connection.SetDeadline(time.Now().Add(config.Values.Timeout))
	defer t.connection.SetDeadline(time.Time{})
	_, err := tivoRequestMessage.WriteTo(t.connection)
	if err != nil {
		return err
	}

	t.rpcID++

	return nil
}

func (t *TivoClient) SendFile(ctx context.Context, filename string) error {

	message, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	t.connection.SetDeadline(time.Now().Add(config.Values.Timeout))
	defer t.connection.SetDeadline(time.Time{})
	_, err = t.connection.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func (t *TivoClient) Receive(ctx context.Context, tivoMessage ReaderFromWriterTo) error {

	responseReader := bufio.NewReader(t.connection)

	t.connection.SetDeadline(time.Now().Add(config.Values.Timeout))
	defer t.connection.SetDeadline(time.Time{})
	_, err := tivoMessage.ReadFrom(responseReader)
	if err != nil {
		return err
	}
	test.Debug(tivoMessage, (fmt.Sprintf("%03d-response.txt", t.rpcID)))

	return nil
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
	keyLog, err := os.Create(path.Join(test.MustGetDebugDir(), "keys.log"))
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
			logz.Logger.Info("certificate request")
			return nil, nil
		},
		GetClientCertificate: func(cri *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			logz.Logger.Info("client certificate request")
			return cert, nil
		},
		RootCAs:            certPool,
		ServerName:         tivo.ServerName(),
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          certPool,
		InsecureSkipVerify: true,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
		KeyLogWriter:       keyLog,
	}, nil
}
