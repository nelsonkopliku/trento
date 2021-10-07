package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/trento-project/trento/internal/cluster"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	clusterDiscovery   = "clusterDiscovery"
	sapsystemDiscovery = "sapsystemDiscovery"
)

type DataCollectedEvent struct {
	ID            int64
	CreatedAt     time.Time
	AgentID       string         `json:"agent_id" binding:"required"`
	DiscoveryType string         `json:"discovery_type" binding:"required"`
	Payload       datatypes.JSON `json:"payload" binding:"required"`
}

func (d *DataCollectedEvent) Validate() error {
	switch d.DiscoveryType {
	case clusterDiscovery:
		b, _ := d.Payload.MarshalJSON()
		return json.Unmarshal(b, &cluster.Cluster{})
	default:
		return fmt.Errorf("unknown DataCollectedEvent type: %s", d.DiscoveryType)
	}
}

type CollectorService interface {
	StoreEvent(dataCollected *DataCollectedEvent) error
}

type collectorService struct {
	db                *gorm.DB
	projectorRegistry interface{}
}

func NewCollectorService(db *gorm.DB, projectorRegistry interface{}) *collectorService {
	return &collectorService{db: db, projectorRegistry: projectorRegistry}
}

func (c *collectorService) StoreEvent(collectedData *DataCollectedEvent) error {
	if err := collectedData.Validate(); err != nil {
		return err
	}

	if err := c.db.Create(collectedData).Error; err != nil {
		return err
	}

	//c.projectorRegistry.Project(collectedData)
	return nil
}
