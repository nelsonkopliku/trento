package projectors

import (
	"bytes"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/datapipeline
	"github.com/trento-project/trento/datapipelinereadmodels"
	"github.com/trento-project/trento/internal/cluster"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ClusterListHandler(event *datapipeline.DataCollectedEvent, db *gorm.DB) error {
	data, _ := event.Payload.MarshalJSON()
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	var cluster cluster.Cluster
	if err := dec.Decode(&cluster); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	clusterListReadModel, err := transformClusterListData(&cluster)
	if err != nil {
		log.Errorf("can't transform data: %s", err)
		return err
	}

	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(clusterListReadModel).Error
}

// pure function!!
func transformClusterListData(cluster *cluster.Cluster) (*readmodels.Cluster, error) {
	return &readmodels.Cluster{
		ID:              cluster.Id,
		Name:            cluster.Name,
		ClusterType:     detectClusterType(cluster),
		SID:             getHanaSID(cluster),
		ResourcesNumber: cluster.Crmmon.Summary.Resources.Number,
		HostsNumber:     cluster.Crmmon.Summary.Nodes.Number,
	}, nil
}

func detectClusterType(cluster *cluster.Cluster) string {
	var hasSapHanaTopology, hasSAPHanaController, hasSAPHana bool

	for _, c := range cluster.Crmmon.Clones {
		for _, r := range c.Resources {
			switch r.Agent {
			case "ocf::suse:SAPHanaTopology":
				hasSapHanaTopology = true
			case "ocf::suse:SAPHana":
				hasSAPHana = true
			case "ocf::suse:SAPHanaController":
				hasSAPHanaController = true
			}
		}
	}

	switch {
	case hasSapHanaTopology && hasSAPHana:
		return readmodels.ClusterTypeScaleUp
	case hasSapHanaTopology && hasSAPHanaController:
		return readmodels.ClusterTypeScaleOut
	default:
		return readmodels.ClusterTypeUnknown
	}
}

func getHanaSID(c *cluster.Cluster) string {
	for _, r := range c.Cib.Configuration.Resources.Clones {
		if r.Primitive.Type == "SAPHanaTopology" {
			for _, a := range r.Primitive.InstanceAttributes {
				if a.Name == "SID" {
					return a.Value
				}
			}
		}
	}

	return ""
}
