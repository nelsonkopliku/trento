package discovery

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
	log "github.com/sirupsen/logrus"

	"github.com/trento-project/trento/agent/collector"
	"github.com/trento-project/trento/internal/consul"
)

type Discovery interface {
	// Returns an arbitrary unique string identifier of the discovery, so that we can associate it to a Consul check ID
	GetId() string
	// Execute the discovery mechanism
	Discover() (string, error)
}

type BaseDiscovery struct {
	id              string
	client          consul.Client
	host            string
	collectorConfig *collector.CollectorConfig
	machineId       string
}

func (d BaseDiscovery) GetId() string {
	return d.id
}

// Execute one iteration of a discovery and store the result in the Consul KVStore.
func (d BaseDiscovery) Discover() (string, error) {
	d.host, _ = os.Hostname()
	return "Basic discovery example", nil
}

func (d *BaseDiscovery) withDataCollectionAndLegacyConsulSupport(discoveryId string, collectorConfig collector.CollectorConfig, client consul.Client) {
	d.id = discoveryId
	d.collectorConfig = &collectorConfig
	d.client = client
	d.initialize()

	// if err := checkDataCollectorConnectionOptions(*d.collectorConfig); err != nil {
	// 	return errors.Wrap(err, "Not enough options provided to initialize connection to Data Collector")
	// }
}

func (d *BaseDiscovery) withLegacyConsulSupport(discoveryId string, client consul.Client) {
	d.id = discoveryId
	d.client = client
	d.initialize()
}

func (d *BaseDiscovery) initialize() {
	d.host, _ = os.Hostname()
	machineId, _ := os.ReadFile("/etc/machine-id") // what if it breaks? can it actually break?
	d.machineId = string(machineId)
}

func (d BaseDiscovery) publishDiscoveredData(discoveredData interface{}) error {
	collectorConfig := d.collectorConfig

	if err := checkDataCollectorConnectionOptions(*collectorConfig); err != nil {
		return errors.Wrap(err, "Not enough options provided to initialize connection to Data Collector")
	}

	if !collectorConfig.Enabled {
		return nil
	}

	cert, err := ioutil.ReadFile(d.collectorConfig.TLS.CACert)
	if err != nil {
		log.Fatalf("could not open CA certificate file: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certificate, err := tls.LoadX509KeyPair(collectorConfig.TLS.ClientCert, collectorConfig.TLS.ClientKey)
	if err != nil {
		log.Fatalf("could not load either client certificate or client key: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // for development purposes
				RootCAs:            caCertPool,
				Certificates:       []tls.Certificate{certificate},
			},
		},
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"agent_id":       d.machineId,
		"discovery_type": d.GetId(),
		"payload":        discoveredData,
	})
	if err != nil {
		log.Errorf("unable to decode data: %s", err)
	}

	endpoint := fmt.Sprintf("https://%s:%d/api/collect_data", collectorConfig.Host, collectorConfig.Port)

	resp, err := client.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil
}

func checkDataCollectorConnectionOptions(collectorConfig collector.CollectorConfig) error {
	var err error

	if !collectorConfig.Enabled {
		return nil
	}
	if collectorConfig.Host == "" {
		err = fmt.Errorf("you must provide the host of the data collector")
	}
	if collectorConfig.TLS.CACert == "" {
		err = errors.Wrap(err, "you must provide a CA certificate")
	}
	if collectorConfig.TLS.ClientCert == "" {
		err = errors.Wrap(err, "you must provide a Client Certificate")
	}
	if collectorConfig.TLS.ClientKey == "" {
		err = errors.Wrap(err, "you must provide a Client Key")
	}

	return err
}
