package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dataPipelineServices "github.com/trento-project/trento/data_pipeline/services"
)

func ApiCollectDataHandler(collectorService dataPipelineServices.CollectorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var e dataPipelineServices.DataCollectedEvent

		err := c.BindJSON(&e)
		if err != nil {
			_ = c.Error(err)
			return
		}

		err = collectorService.StoreEvent(&e)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"stored": "ok"})
	}
}
