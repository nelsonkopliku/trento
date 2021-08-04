package services

import (
	"github.com/mitchellh/mapstructure"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services/ara"
)

//go:generate mockery --name=ChecksService

type ChecksService interface {
	GetChecksCatalog() (map[string]*models.Check, error)
	GetChecksCatalogByGroup() (map[string]map[string]*models.Check, error)
}

type checksService struct {
	araService ara.AraService
}

func NewChecksService(araService ara.AraService) ChecksService {
	return &checksService{araService: araService}
}

func (c *checksService) GetChecksCatalog() (map[string]*models.Check, error) {
	checkList := make(map[string]*models.Check)

	records, err := c.araService.GetRecordList("key=trento-metadata&order=-id")
	if err != nil {
		return checkList, err
	}

	if len(records.Results) == 0 {
		return checkList, nil
	}

	record, err := c.araService.GetRecord(records.Results[0].ID)
	if err != nil {
		return checkList, err
	}

	mapstructure.Decode(record.Value, &checkList)

	return checkList, nil
}

func (c *checksService) GetChecksCatalogByGroup() (map[string]map[string]*models.Check, error) {
	groupedCheckList := make(map[string]map[string]*models.Check)

	checkList, err := c.GetChecksCatalog()
	if err != nil {
		return groupedCheckList, err
	}

	for cId, c := range checkList {
		extendedGroup := c.ExtendedGroupName()
		if _, ok := groupedCheckList[extendedGroup]; !ok {
			groupedCheckList[extendedGroup] = make(map[string]*models.Check)
		}
		groupedCheckList[extendedGroup][cId] = c
	}

	return groupedCheckList, nil
}