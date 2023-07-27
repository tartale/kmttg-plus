package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/tartale/go/pkg/errorx"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/test"
)

const tivoRPCPort = "1413"

type TivoClient struct {
	address    string
	tlsConfig  *tls.Config
	connection *tls.Conn
	sessionID  string
	rpcID      int
	tsn        string
}

func NewTivoClient(tivo *model.Tivo) (*TivoClient, error) {

	tlsConfig, err := NewTLSConfig(tivo)
	if err != nil {
		return nil, err
	}

	tivoClient := &TivoClient{
		address:   tivo.Address,
		tlsConfig: tlsConfig,
		tsn:       tivo.Tsn,
	}
	err = tivoClient.Connect()
	if err != nil {
		return nil, err
	}

	return tivoClient, nil
}

func NewTLSConfig(tivo *model.Tivo) (*tls.Config, error) {

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
		GetClientCertificate: func(cri *tls.CertificateRequestInfo) (*tls.Certificate, error) {
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

func (t *TivoClient) Close() error {

	return t.connection.Close()
}

func (t *TivoClient) Connect() error {

	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, config.Values.Timeout)
	defer cancelFunc()
	dialer := tls.Dialer{
		NetDialer: new(net.Dialer),
		Config:    t.tlsConfig,
	}

	conn, err := dialer.DialContext(ctx, "tcp", t.address+":"+tivoRPCPort)
	if err != nil {
		return err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sessionID := fmt.Sprintf("0x%x", r.Int())[:9]

	t.connection = conn.(*tls.Conn)
	t.sessionID = sessionID
	t.rpcID = 1

	return nil
}

func (t *TivoClient) Reconnect() error {

	t.connection.Close()
	return t.Connect()
}

func (t *TivoClient) BodyID() string {
	return "tsn:" + t.tsn
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
		return errorz.ErrNotAuthenticated(authResponseBody.Message)
	}

	return nil
}

func (t *TivoClient) GetAllRecordings(ctx context.Context) ([]model.Show, error) {

	request := message.NewTivoMessage().WithGetAllRecordingsRequest(ctx, t.BodyID())
	err := t.Send(ctx, request)
	if err != nil {
		return nil, err
	}

	responseBody := &message.RecordingFolderItemSearchResponseBody{}
	response := message.NewTivoMessage().WithBody(responseBody)
	err = t.Receive(ctx, response)
	if err != nil {
		return nil, err
	}
	if responseBody.Type != message.TypeRecordingFolderItemList {
		return nil, errorz.ErrResponse(responseBody.Message)
	}
	var result []model.Show
	var errs errorx.Errors
	for _, recording := range responseBody.RecordingFolderItem {
		show, err := t.GetRecording(ctx, recording.ChildRecordingID)
		if err != nil {
			errs = append(errs, err)
		} else {
			result = append(result, show)
		}
	}

	return result, errs.Combine("", "; ")
}

func (t *TivoClient) GetRecording(ctx context.Context, recordingID string) (model.Show, error) {

	request := message.NewTivoMessage().WithGetRecordingRequest(ctx, t.BodyID(), recordingID)
	err := t.Send(ctx, request)
	if err != nil {
		return nil, err
	}

	responseBody := &message.RecordingSearchResponseBody{}
	response := message.NewTivoMessage().WithBody(responseBody)
	err = t.Receive(ctx, response)
	if err != nil {
		return nil, err
	}
	if responseBody.Type != message.TypeRecordingList {
		return nil, errorz.ErrResponse(responseBody.Message)
	}
	recordingCount := len(responseBody.Recording)
	if recordingCount != 1 {
		return nil, fmt.Errorf("unexpected number of recordings in response: %d", recordingCount)
	}
	show, err := model.NewShow(responseBody.Recording[0])
	if err != nil {
		return nil, err
	}

	return show, nil
}

func (t *TivoClient) Send(ctx context.Context, tivoMessage *message.TivoMessage) error {

	tivoRequestMessage := tivoMessage.
		WithSessionID(t.sessionID).
		WithRpcID(t.rpcID)

	test.Debug(tivoRequestMessage, (fmt.Sprintf("%03d-request.txt", t.rpcID)))

	err := t.Retry(func() error {
		_, err := tivoRequestMessage.WriteTo(t.connection)
		return err
	})
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

	err = t.Retry(func() error {
		_, err := t.connection.Write(message)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *TivoClient) Receive(ctx context.Context, tivoMessage *message.TivoMessage) error {

	responseReader := bufio.NewReader(t.connection)

	err := t.Retry(func() error {
		_, err := tivoMessage.ReadFrom(responseReader)
		return err
	})
	if err != nil {
		return err
	}
	test.Debug(tivoMessage, fmt.Sprintf("%03d-response.txt", tivoMessage.Headers.RpcID()))

	return nil
}

func (t *TivoClient) Retry(fn func() error) error {

	t.connection.SetDeadline(time.Now().Add(config.Values.Timeout))
	defer t.connection.SetDeadline(time.Time{})
	for err := fn(); err != nil; {
		if errors.Is(err, syscall.EPIPE) {
			t.Reconnect()
			continue
		}
		return err
	}

	return nil
}
