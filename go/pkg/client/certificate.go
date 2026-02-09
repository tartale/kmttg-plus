package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"

	"golang.org/x/crypto/pkcs12"
)

const (
	certificatePassword = "SKA2Kvgxvs"
)

func GetCertificates(certificatePath string) (*tls.Certificate, *x509.CertPool, error) {
	certificateBytes, err := os.ReadFile(certificatePath)
	if err != nil {
		return nil, nil, err
	}

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
