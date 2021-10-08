package projectors

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/web/datapipeline"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StartProjectorsWorkerPool(workersNumber int, db *gorm.DB) chan *datapipeline.DataCollectedEvent {
	// trovare un modo per farla diventare non-blocking
	ch := make(chan *datapipeline.DataCollectedEvent, workersNumber)

	for i := 0; i < workersNumber; {
		go Worker(ch, db)
		i++
	}

	return ch
}

func Worker(ch chan *datapipeline.DataCollectedEvent, db *gorm.DB) {
	for event := range ch {
		switch event.DiscoveryType {
		case datapipeline.ClusterDiscovery:
			Project(event, db, ClusterListHandler)
		default:
			log.Errorf("unknown discovery type: %s", event.DiscoveryType)
		}
		fmt.Println("Received event:", event.DiscoveryType)
	}
}

func Project(event *datapipeline.DataCollectedEvent, db *gorm.DB, handler func(*datapipeline.DataCollectedEvent, *gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&Subscription{
			DiscoveryType: event.DiscoveryType,
			AgentID:       event.AgentID,
			EventID:       event.ID,
		})

		err := handler(event, tx)
		if err != nil {
			return err
		}

		return nil
	})
}
