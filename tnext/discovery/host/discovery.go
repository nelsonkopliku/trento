package host

import "github.com/trento-project/trento/tnext/discovery"

type hostDiscovery struct{}

func NewDiscovery() discovery.Discovery {
	return &hostDiscovery{}
}

func (hostDiscovery) Type() discovery.DiscoveryType {
	return discovery.Host
}

func (hostDiscovery) Discover() error {
	return nil
}
