package services

import (
	"github.com/trento-project/trento/web/datapipeline"
	"gorm.io/gorm"
)

type CollectorService interface {
	StoreEvent(dataCollected *datapipeline.DataCollectedEvent) error
}

type collectorService struct {
	db                *gorm.DB
	projectorsChannel chan *datapipeline.DataCollectedEvent
}

func NewCollectorService(db *gorm.DB, projectorsChannel chan *datapipeline.DataCollectedEvent) *collectorService {
	return &collectorService{db: db, projectorsChannel: projectorsChannel}
}

func (c *collectorService) StoreEvent(collectedData *datapipeline.DataCollectedEvent) error {
	if err := c.db.Create(collectedData).Error; err != nil {
		return err
	}
	c.projectorsChannel <- collectedData

	return nil
}
