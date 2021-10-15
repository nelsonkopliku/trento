package discovery

import (
	"fmt"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/subscription"
)

const SubscriptionDiscoveryId string = "subscription_discovery"

type SubscriptionDiscovery struct {
	BaseDiscovery
}

func NewSubscriptionDiscovery(client consul.Client) SubscriptionDiscovery {
	discovery := SubscriptionDiscovery{}
	discovery.withLegacyConsulSupport(SubscriptionDiscoveryId, client)
	return discovery
}

func (d SubscriptionDiscovery) GetId() string {
	return d.id
}

func (d SubscriptionDiscovery) Discover() (string, error) {
	subsData, err := subscription.NewSubscriptions()
	if err != nil {
		return "", err
	}

	err = subsData.Store(d.client)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Subscription (%d entries) discovered", len(subsData)), nil
}
