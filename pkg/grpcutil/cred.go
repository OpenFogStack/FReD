package grpcutil

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/go-errors/errors"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func getSecureCreds(certFile string, keyFile string, caFiles []string, config *tls.Config) (credentials.TransportCredentials, *x509.CertPool, error) {
	if certFile == "" {
		return nil, nil, errors.New("no certificate file given")
	}

	if keyFile == "" {
		return nil, nil, errors.New("no key file given")
	}

	if len(caFiles) == 0 {
		return nil, nil, errors.New("no root certificate file given")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		return nil, nil, err
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	for _, f := range caFiles {
		loaded, err := os.ReadFile(f)

		if err != nil {
			return nil, nil, err
		}

		rootCAs.AppendCertsFromPEM(loaded)
	}

	config.Certificates = []tls.Certificate{serverCert}
	config.ClientCAs = rootCAs
	config.RootCAs = rootCAs
	config.MinVersion = tls.VersionTLS12

	return credentials.NewTLS(config), rootCAs, nil
}

func getInsecureCreds() credentials.TransportCredentials {
	return insecure.NewCredentials()
}

func GetCreds(certFile string, keyFile string, caFiles []string, insecure bool) (credentials.TransportCredentials, *x509.CertPool, error) {
	return GetCredsFromConfig(certFile, keyFile, caFiles, insecure, &tls.Config{})
}

func GetCredsFromConfig(certFile string, keyFile string, caFiles []string, insecure bool, config *tls.Config) (credentials.TransportCredentials, *x509.CertPool, error) {
	if insecure {
		return getInsecureCreds(), nil, nil
	}

	return getSecureCreds(certFile, keyFile, caFiles, config)
}
