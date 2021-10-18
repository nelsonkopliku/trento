package collector

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type CollectorClient interface {
	Publish(discoveryType string, payload interface{}) error
}

type collectorClient struct {
	cfg        Config
	machineId  string
	httpClient *http.Client
}

type Config struct {
	CollectorHost string
	CollectorPort int
	EnablemTLS    bool
	Cert          string
	Key           string
	CA            string
}

func NewCollectorClient(cfg Config) (*collectorClient, error) {
	var tlsConfig *tls.Config
	var err error

	if cfg.EnablemTLS {
		tlsConfig, err = getTLSConfig(cfg.Cert, cfg.Key, cfg.CA)
		if err != nil {
			return nil, err
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	machineId, err := os.ReadFile("/etc/machine-id") // what if it breaks? can it actually break?
	if err != nil {
		return nil, err
	}

	return &collectorClient{
		cfg:        cfg,
		httpClient: client,
		machineId:  string(machineId),
	}, nil
}

func (c *collectorClient) Publish(discoveryType string, payload interface{}) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"agent_id":       c.machineId,
		"discovery_type": discoveryType,
		"payload":        payload,
	})
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/api/collect_data", c.cfg.CollectorHost)
	resp, err := c.httpClient.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("something wrong happened while publishing data to the collector. Agent: %s, discovery: %s", c.machineId, discoveryType)
	}

	return nil
}

func getTLSConfig(cert, key, ca string) (*tls.Config, error) {
	var err error
	if cert == "" {
		err = fmt.Errorf("you must provide a server ssl certificate")
	}
	if key == "" {
		err = errors.Wrap(err, "you must provide a key to enable mTLS")
	}
	if ca == "" {
		err = errors.Wrap(err, "you must provide a CA ssl certificate")
	}
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            caCertPool,
		Certificates:       []tls.Certificate{certificate},
	}, nil
}
