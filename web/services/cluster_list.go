package services

import (
	"fmt"

	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/web/datapipeline/readmodels"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

type ClusterListService interface {
	GetAll(filters map[string][]string) (readmodels.ClusterList, error)
}

type clusterListService struct {
	db            *gorm.DB
	checksService ChecksService
	tagsService   TagsService
}

func NewClusterList(db *gorm.DB, checksService ChecksService, tagsService TagsService) *clusterListService {
	return &clusterListService{
		db:            db,
		checksService: checksService,
		tagsService:   tagsService,
	}
}

func (s *clusterListService) GetAll(filters map[string][]string) (readmodels.ClusterList, error) {
	var clusterList readmodels.ClusterList
	db := s.db

	for _, f := range []string{"name", "sid", "cluster_type"} {
		if v, ok := filters[f]; ok {
			if len(v) > 0 {
				q := fmt.Sprintf("%s IN (?)", f)
				db = s.db.Where(q, v)
			}
		}
	}

	err := db.Find(&clusterList).Error
	if err != nil {
		return nil, err
	}

	err = s.enrichClusterData(clusterList)
	if err != nil {
		return nil, err
	}

	if tagsFilter, ok := filters["tags"]; ok {
		clusterList = filterByTags(clusterList, tagsFilter)
	}

	if healthFilter, ok := filters["health"]; ok {
		clusterList = filterByHealth(clusterList, healthFilter)
	}

	//db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count) > 0
	return clusterList, nil
}

func (s *clusterListService) enrichClusterData(clusterList readmodels.ClusterList) error {
	names := make(map[string]int)
	for _, c := range clusterList {
		names[c.Name] += 1
	}

	for _, c := range clusterList {
		if names[c.Name] > 1 {
			c.HasDuplicatedName = true
		}
		health, _ := s.checksService.GetAggregatedChecksResultByCluster(c.ID)
		c.Health = health.String()

		tags, err := s.tagsService.GetAllByResource(models.TagClusterResourceType, c.ID)
		if err != nil {
			return err
		}
		c.Tags = tags
	}

	return nil
}

func filterByTags(clusterList readmodels.ClusterList, tags []string) readmodels.ClusterList {
	var filteredClusterList readmodels.ClusterList

	for _, c := range clusterList {
		for _, t := range tags {
			if internal.Contains(c.Tags, t) {
				filteredClusterList = append(filteredClusterList, c)
				break
			}
		}
	}

	return filteredClusterList
}

func filterByHealth(clusterList readmodels.ClusterList, health []string) readmodels.ClusterList {
	var filteredClusterList readmodels.ClusterList

	for _, c := range clusterList {
		if internal.Contains(health, c.Health) {
			filteredClusterList = append(filteredClusterList, c)
		}
	}

	return filteredClusterList
}
