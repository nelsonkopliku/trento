package cloud

import (
	"github.com/trento-project/trento/tnext/discovery"
	"github.com/trento-project/trento/tnext/internal/cloud"
)

type cloudDiscovery struct{
	cloudDetector Detector
}

func NewDiscovery() discovery.Discovery {
	return &cloudDiscovery{}
}

func (cloudDiscovery) Type() discovery.DiscoveryType {
	return discovery.Cloud
}

func (cloudDiscovery) Discover() (discovery.DiscoveryResult, error) {
	cloudData, err := cloud.NewCloudInstance()

	return cloudData, err
}
