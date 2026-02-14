package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/tartale/go/pkg/retry"
	"github.com/tartale/kmttg-plus/go/pkg/certificate"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
)

const (
	tivoRPCPort       = "1413"
	heartbeatInterval = 5 * time.Second
)

// We need to specify older cipher suites
// https://github.com/golang/go/issues/66512
var cipherSuites = []uint16{
	// TLS 1.0 - 1.2 cipher suites.
	tls.TLS_RSA_WITH_RC4_128_SHA,
	tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
	// TLS 1.3 cipher suites.
	tls.TLS_AES_128_GCM_SHA256,
	tls.TLS_AES_256_GCM_SHA384,
	tls.TLS_CHACHA20_POLY1305_SHA256,
}

type TivoClient struct {
	tivo          *model.Tivo
	address       string
	tlsConfig     *tls.Config
	connection    *tls.Conn
	sessionID     string
	rpcID         int
	tsn           string
	lastHeartbeat time.Time
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
		rpcID:     1,
	}
	err = tivoClient.Connect()
	if err != nil {
		return nil, err
	}
	tivoClient.tivo = tivo

	return tivoClient, nil
}

func NewTLSConfig(tivo *model.Tivo) (*tls.Config, error) {
	cert, certPool, err := certificate.Read()
	if err != nil {
		return nil, err
	}
	var KeyLogWriter io.Writer
	if config.Values.LogMessages {
		keyLog, err := os.Create(path.Join(logz.MustGetDebugDir(), "keys.log"))
		if err != nil {
			return nil, err
		}
		KeyLogWriter = keyLog
	}

	return &tls.Config{
		GetClientCertificate: func(cri *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return cert, nil
		},
		CipherSuites:       cipherSuites,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          certPool,
		InsecureSkipVerify: true,
		KeyLogWriter:       KeyLogWriter,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
		RootCAs:            certPool,
		ServerName:         tivo.ServerName(),
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
	t.lastHeartbeat = time.Now()

	return nil
}

func (t *TivoClient) Reconnect(ctx context.Context, cause error) error {
	logz.Logger.Warn("reconnecting client due to error", zap.Error(cause))
	t.connection.Close()
	err := t.Connect()
	if err != nil {
		return err
	}
	err = t.Authenticate(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t *TivoClient) BodyID() string {
	return "tsn:" + t.tsn
}

func (t *TivoClient) Send(ctx context.Context, tivoMessage *message.TivoMessage) error {
	tivoRequestMessage := tivoMessage.
		WithSessionID(t.sessionID).
		WithRpcID(t.rpcID)

	err := t.ensureConnection(ctx)
	if err != nil {
		return err
	}
	_, err = tivoRequestMessage.WriteTo(t.connection)
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
	if err := t.ensureConnection(ctx); err != nil {
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

	err := t.ensureConnection(ctx)
	if err != nil {
		return err
	}
	_, err = tivoMessage.ReadFrom(responseReader)
	if err != nil {
		return err
	}

	return nil
}

func (t *TivoClient) Write(b []byte) (int, error) {
	err := t.ensureConnection(context.Background())
	if err != nil {
		return 0, err
	}

	return t.connection.Write(b)
}

func (t *TivoClient) Read(b []byte) (int, error) {
	err := t.ensureConnection(context.Background())
	if err != nil {
		return 0, err
	}

	return t.connection.Read(b)
}

func (t *TivoClient) ensureConnection(ctx context.Context) error {
	var err error
	if time.Now().After(t.lastHeartbeat.Add(heartbeatInterval)) {
		err = t.testConnection()
		if shouldReconnect(err) {
			err = t.Reconnect(ctx, err)
			if err != nil {
				return err
			}
			err = errorz.ErrReconnected
		}
	}
	t.lastHeartbeat = time.Now()

	return err
}

func (t *TivoClient) testConnection() error {
	_, err := t.connection.Write([]byte{})
	if err != nil {
		return err
	}

	return nil
}

func shouldReconnect(err error) bool {
	if _, ok := err.(net.Error); ok {
		return true
	}
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
