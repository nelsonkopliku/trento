package discovery

import (
	"fmt"

	"github.com/trento-project/trento/agent/collector"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
)

const ClusterDiscoveryId string = "ha_cluster_discovery"

// This Discover handles any Pacemaker Cluster type
type ClusterDiscovery struct {
	BaseDiscovery
	dicoveryTools cluster.DiscoveryTools
	Cluster       cluster.Cluster
}

type ClusterDiscoveryOptions struct {
	CollectorConfig       collector.CollectorConfig
	ClusterDiscoverytools cluster.DiscoveryTools
}

func NewClusterDiscovery(client consul.Client, options ClusterDiscoveryOptions) ClusterDiscovery {
	discovery := ClusterDiscovery{}
	discovery.dicoveryTools = options.ClusterDiscoverytools
	discovery.withDataCollectionAndLegacyConsulSupport(ClusterDiscoveryId, options.CollectorConfig, client)
	return discovery
}

func (c ClusterDiscovery) GetId() string {
	return c.id
}

// Execute one iteration of a discovery and store the result in the Consul KVStore.
func (d ClusterDiscovery) Discover() (string, error) {
	cluster, err := cluster.NewCluster(d.dicoveryTools)

	if err != nil {
		return "No HA cluster discovered on this host", nil
	}

	d.Cluster = cluster

	err = d.Cluster.Store(d.client)
	if err != nil {
		return "", err
	}

	err = storeClusterMetadata(d.client, cluster.Name, cluster.Id)
	if err != nil {
		return "", err
	}

	d.publishDiscoveredData(cluster)

	return fmt.Sprintf("Cluster with name: %s successfully discovered", cluster.Name), nil
}

func storeClusterMetadata(client consul.Client, clusterName string, clusterId string) error {
	metadata := hosts.Metadata{
		Cluster:   clusterName,
		ClusterId: clusterId,
	}
	err := metadata.Store(client)
	if err != nil {
		return err
	}

	return nil
}
