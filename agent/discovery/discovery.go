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

// type PublishableDiscovery interface {
// 	Discovery
// 	// Enables this Discovery to publish data to the Data Collector
// 	WithDataCollectorConfig(collectorConfig collector.CollectorConfig) Discovery
// }

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

// Return a Host Discover instance
func NewDiscovery(client consul.Client) BaseDiscovery {
	r := BaseDiscovery{}
	r.id = ""
	r.client = client
	r.host, _ = os.Hostname()
	return r
}

// func NewPublishableDiscovery(inner Discovery, collectorConfig collector.CollectorConfig) Discovery {
// 	baseDiscovery, err := inner.(BaseDiscovery)
// 	if err {
// 		return inner
// 	}

// 	return baseDiscovery.withDataCollectorConfig(collectorConfig)
// 	// switch innerDiscovery := inner.(type) {
// 	// case BaseDiscovery:
// 	// 	return innerDiscovery.withDataCollectorConfig(collectorConfig)
// 	// }
// 	// return inner
// }

func (d *BaseDiscovery) init() {
	d.host, _ = os.Hostname()
	machineId, _ := os.ReadFile("/etc/machine-id") // what if it breaks? can it actually break?
	d.machineId = string(machineId)
}

// func (d BaseDiscovery) withDataCollectorConfig(collectorConfig collector.CollectorConfig) Discovery {
// 	d.collectorConfig = &collectorConfig
// 	return d
// }

func (d BaseDiscovery) publishDiscoveredData(discoveredData interface{}) error {
	collectorConfig := d.collectorConfig
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
		log.Error("unable to decode data")
	}

	endpoint := fmt.Sprintf("%s/api/collect_data", collectorConfig.Host)

	resp, err := client.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil
}
