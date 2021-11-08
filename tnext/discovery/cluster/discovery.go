package cluster

import "github.com/trento-project/trento/tnext/discovery"

type clusterDiscovery struct{}

func NewDiscovery() discovery.Discovery {
	return &clusterDiscovery{}
}

func (clusterDiscovery) Type() discovery.DiscoveryType {
	return discovery.Cluster
}

func (clusterDiscovery) Discover() error {
	return nil
}
