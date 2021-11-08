package discovery

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type DiscoveryType string

const (
	Cloud        DiscoveryType = "cloud_discovery"
	Host         DiscoveryType = "host_discovery"
	Cluster      DiscoveryType = "cluster_discovery"
	SAP          DiscoveryType = "sap_discovery"
	Subscription DiscoveryType = "subscription_discovery"
)

type DiscoveredContent interface{}

type DiscoveryResult struct {
	Error   error
	Content DiscoveredContent
}

type Discovery interface {
	Type() DiscoveryType
	Discover() DiscoveryResult
}

type RegisteredDiscoveries []Discovery

type DiscoveryRunner interface {
	RegisterDiscovery(discovery Discovery)
	Run() error
}

type discoveryRunner struct {
	discoveries RegisteredDiscoveries
}

func (d *discoveryRunner) RegisterDiscovery(discovery Discovery) {
	d.discoveries = append(d.discoveries, discovery)
}

func (d *discoveryRunner) Run() error {
	for _, discovery := range d.discoveries {
		discoveryType, discoveryResult := discovery.Type(), discovery.Discover()

		if err := discoveryResult.Error; err != nil {
			log.Errorln(fmt.Sprintf("Error while running discovery '%s': %s", discoveryType, err))
			continue
		}

		log.Println(fmt.Sprintf("Discovery '%s' was successful", discoveryType))

		// do something with discoveryResult.Content
	}
	return nil
}

func NewDiscoveryRunner() DiscoveryRunner {
	return &discoveryRunner{}
}
