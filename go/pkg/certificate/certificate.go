package certificate

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"

	"github.com/tartale/kmttg-plus/go/assets"
	"golang.org/x/crypto/pkcs12"
)

const (
	certificateFilename         = "cdata.p12"
	certificatePasswordFilename = "cdata.password"
)

func Read() (*tls.Certificate, *x509.CertPool, error) {
	return ReadFrom(bytes.NewReader(assets.CertificateZipBytes))
}

func ReadFrom(certificateReader *bytes.Reader) (*tls.Certificate, *x509.CertPool, error) {
	zipReader, err := zip.NewReader(certificateReader, certificateReader.Size())
	if err != nil {
		return nil, nil, err
	}
	certificateFile, err := zipReader.Open(certificateFilename)
	if err != nil {
		return nil, nil, err
	}
	certificateBytes, err := io.ReadAll(certificateFile)
	if err != nil {
		return nil, nil, err
	}
	certificatePasswordFile, err := zipReader.Open(certificatePasswordFilename)
	if err != nil {
		return nil, nil, err
	}
	certificatePasswordBytes, err := io.ReadAll(certificatePasswordFile)
	if err != nil {
		return nil, nil, err
	}
	certificatePassword := string(certificatePasswordBytes)

	// the following is copied from this example:
	// https://pkg.go.dev/golang.org/x/crypto/pkcs12#example-ToPEM
	pemBlocks, err := pkcs12.ToPEM(certificateBytes, certificatePassword)
	if err != nil {
		return nil, nil, err
	}

	var pemData []byte
	for _, b := range pemBlocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		return nil, nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(pemData)

	return &cert, certPool, nil
}
