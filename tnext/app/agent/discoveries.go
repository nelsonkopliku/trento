package agent

import (
	"github.com/trento-project/trento/tnext/discovery"
	"github.com/trento-project/trento/tnext/discovery/cloud"
	"github.com/trento-project/trento/tnext/discovery/cluster"
	"github.com/trento-project/trento/tnext/discovery/host"
)

func NewConfiguredDiscoveryRunner() discovery.DiscoveryRunner {
	discoveryRunner := discovery.NewDiscoveryRunner()

	discoveryRunner.RegisterDiscovery(host.NewDiscovery())
	discoveryRunner.RegisterDiscovery(cloud.NewDiscovery())
	discoveryRunner.RegisterDiscovery(cluster.NewDiscovery())

	return discoveryRunner
}
