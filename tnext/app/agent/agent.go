package agent

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/trento-project/trento/tnext/app/agent/collector"
)

type Agent struct {
	ID        uuid.UUID
	config    *Config
	ctx       context.Context
	ctxCancel context.CancelFunc
}

type Config struct {
	DiscoveryPeriod time.Duration
	CollectorConfig *collector.Config
}

func NewAgent(config *Config) (*Agent, error) {
	agentID, err := loadIdentifier()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load agent identifier")
	}

	// should this be here? or passed to the Start function?
	ctx, cancel := context.WithCancel(context.Background())

	return &Agent{
		ID:        agentID,
		config:    config,
		ctx:       ctx,
		ctxCancel: cancel,
	}, nil
}

func (a *Agent) Start() error {
	discoveryRunner := NewConfiguredDiscoveryRunner()

	// Start a Ticker loop that will iterate over the hardcoded list of Discovery backends
	// and execute them. The initial run will happen relatively quickly after Agent launch
	// subsequent runs are done with a 15 minute delay. The effectiveness of the discoveries
	// is reported back in the "discover_cluster" Service in consul under a TTL of 60 minutes

	interval := a.config.DiscoveryPeriod
	repeat(discoveryRunner.Run, interval, a.ctx)

	return nil
}

func (a *Agent) Stop() {
	a.ctxCancel()
}

func repeat(tick func() error, interval time.Duration, ctx context.Context) {
	// run the first tick immediately
	tick()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			tick()
		case <-ctx.Done():
			return
		}
	}
}
