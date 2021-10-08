package projectors

import "time"

type Subscription struct {
	EventID       int64
	AgentID       string `gorm:"primaryKey"`
	DiscoveryType string `gorm:"primaryKey"`
	SeenAt        time.Time
}
