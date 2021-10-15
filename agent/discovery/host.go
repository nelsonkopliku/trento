package discovery

import (
	"fmt"
	"net"
	"strings"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
)

const HostDiscoveryId string = "host_discovery"

type HostDiscovery struct {
	BaseDiscovery
}

func NewHostDiscovery(client consul.Client) HostDiscovery {
	d := HostDiscovery{}
	d.withLegacyConsulSupport(HostDiscoveryId, client)
	return d
}

func (h HostDiscovery) GetId() string {
	return h.id
}

// Execute one iteration of a discovery and store the result in the Consul KVStore.
func (h HostDiscovery) Discover() (string, error) {
	ipAddresses, err := getHostIpAddresses()
	if err != nil {
		return "", err
	}

	metadata := hosts.Metadata{
		HostIpAddresses: ipAddresses,
	}
	err = metadata.Store(h.client)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Host with name: %s successfully discovered", h.host), nil
}

func getHostIpAddresses() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	ipAddrList := make([]string, 0)

	for _, inter := range interfaces {
		addrs, err := inter.Addrs()
		if err != nil {
			continue
		}

		for _, ipaddr := range addrs {
			ipv4Addr, _, _ := net.ParseCIDR(ipaddr.String())
			ipAddrList = append(ipAddrList, ipv4Addr.String())
		}
	}

	return strings.Join(ipAddrList, ","), nil
}
