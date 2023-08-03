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
	"github.com/tartale/go/pkg/retry"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
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
	keyLog, err := os.Create(path.Join(logz.MustGetDebugDir(), "keys.log"))
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

	timeout := config.Values.Timeout
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, timeout)
	defer cancelFunc()
	dialer := tls.Dialer{
		NetDialer: new(net.Dialer),
		Config:    t.tlsConfig,
	}

	var conn net.Conn
	err := retry.Eventually(func() error {
		var err error
		conn, err = dialer.DialContext(ctx, "tcp", t.address+":"+tivoRPCPort)
		return err
	}, timeout, 1*time.Second)
	if err != nil {
		logz.Logger.Warn("unable to connect to tivo", zap.Error(err))
		return err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sessionID := fmt.Sprintf("0x%x", r.Int())[:9]

	t.connection = conn.(*tls.Conn)
	t.sessionID = sessionID
	t.rpcID = 1

	return nil
}

func (t *TivoClient) Reconnect(cause error) error {

	logz.Logger.Warn("reconnecting client due to error", zap.Error(cause))
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
		logz.Logger.Warn("tivo error response", zap.Any("request", request), zap.Any("response", responseBody))
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected response type: %s", responseBody.Type))
	}
	var result []model.Show
	var errs errorx.Errors
	for _, recording := range responseBody.RecordingFolderItem {
		show, err := t.GetRecording(ctx, recording)
		if err != nil {
			errs = append(errs, err)
		} else {
			result = append(result, show)
		}
	}
	result = model.MergeEpisodes(result)

	return result, errs.Combine("", "; ")
}

func (t *TivoClient) GetRecording(ctx context.Context, parent message.RecordingFolderItem) (model.Show, error) {

	request := message.NewTivoMessage().WithGetRecordingRequest(ctx, t.BodyID(), parent.ChildRecordingID)
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
		logz.Logger.Error("tivo error response", zap.Any("responseBody", responseBody))
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected response type: %s", responseBody.Type))
	}
	recordingCount := len(responseBody.Recording)
	if recordingCount != 1 {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected number of recordings in response: %d", recordingCount))
	}
	show, err := model.NewShow(responseBody.Recording[0], parent)
	if err != nil {
		return nil, err
	}

	return show, nil
}

func (t *TivoClient) Send(ctx context.Context, tivoMessage *message.TivoMessage) error {

	tivoRequestMessage := tivoMessage.
		WithSessionID(t.sessionID).
		WithRpcID(t.rpcID)

	logz.Debug(tivoRequestMessage, (fmt.Sprintf("%03d-request.txt", t.rpcID)))

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

	_, err = t.connection.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func (t *TivoClient) Receive(ctx context.Context, tivoMessage *message.TivoMessage) error {

	responseReader := bufio.NewReader(t.connection)

	_, err := tivoMessage.ReadFrom(responseReader)
	if err != nil {
		return err
	}
	logz.Debug(tivoMessage, fmt.Sprintf("%03d-response.txt", tivoMessage.Headers.RpcID()))

	return nil
}

func (t *TivoClient) Retry(fn func() error) error {

	t.connection.SetDeadline(time.Now().Add(config.Values.Timeout))
	defer t.connection.SetDeadline(time.Time{})
	var err error
	for err = fn(); err != nil; {
		if shouldReconnect(err) {
			t.Reconnect(err)
			continue
		}
		break
	}

	return err
}

func shouldReconnect(err error) bool {
	if errors.Is(err, syscall.EPIPE) {
		return true
	}
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return true
	}
	if opError, ok := err.(*net.OpError); ok {
		if !opError.Temporary() {
			return true
		}
	}

	return false
}
