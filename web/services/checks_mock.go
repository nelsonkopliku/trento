// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package services

import (
	mock "github.com/stretchr/testify/mock"
	models "github.com/trento-project/trento/web/models"
)

// MockChecksService is an autogenerated mock type for the ChecksService type
type MockChecksService struct {
	mock.Mock
}

// CreateChecksCatalog provides a mock function with given fields: checkList
func (_m *MockChecksService) CreateChecksCatalog(checkList models.ChecksCatalog) error {
	ret := _m.Called(checkList)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.ChecksCatalog) error); ok {
		r0 = rf(checkList)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateChecksCatalogEntry provides a mock function with given fields: check
func (_m *MockChecksService) CreateChecksCatalogEntry(check *models.Check) error {
	ret := _m.Called(check)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Check) error); ok {
		r0 = rf(check)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateConnectionSettings provides a mock function with given fields: node, cluster, user
func (_m *MockChecksService) CreateConnectionSettings(node string, cluster string, user string) error {
	ret := _m.Called(node, cluster, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(node, cluster, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateSelectedChecks provides a mock function with given fields: id, selectedChecksList
func (_m *MockChecksService) CreateSelectedChecks(id string, selectedChecksList []string) error {
	ret := _m.Called(id, selectedChecksList)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string) error); ok {
		r0 = rf(id, selectedChecksList)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAggregatedChecksResultByCluster provides a mock function with given fields: clusterId
func (_m *MockChecksService) GetAggregatedChecksResultByCluster(clusterId string) (*AggregatedCheckData, error) {
	ret := _m.Called(clusterId)

	var r0 *AggregatedCheckData
	if rf, ok := ret.Get(0).(func(string) *AggregatedCheckData); ok {
		r0 = rf(clusterId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AggregatedCheckData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(clusterId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAggregatedChecksResultByHost provides a mock function with given fields: clusterId
func (_m *MockChecksService) GetAggregatedChecksResultByHost(clusterId string) (map[string]*AggregatedCheckData, error) {
	ret := _m.Called(clusterId)

	var r0 map[string]*AggregatedCheckData
	if rf, ok := ret.Get(0).(func(string) map[string]*AggregatedCheckData); ok {
		r0 = rf(clusterId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*AggregatedCheckData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(clusterId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChecksCatalog provides a mock function with given fields:
func (_m *MockChecksService) GetChecksCatalog() (models.ChecksCatalog, error) {
	ret := _m.Called()

	var r0 models.ChecksCatalog
	if rf, ok := ret.Get(0).(func() models.ChecksCatalog); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(models.ChecksCatalog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChecksCatalogByGroup provides a mock function with given fields:
func (_m *MockChecksService) GetChecksCatalogByGroup() (models.GroupedCheckList, error) {
	ret := _m.Called()

	var r0 models.GroupedCheckList
	if rf, ok := ret.Get(0).(func() models.GroupedCheckList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(models.GroupedCheckList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChecksResult provides a mock function with given fields:
func (_m *MockChecksService) GetChecksResult() (map[string]*models.Results, error) {
	ret := _m.Called()

	var r0 map[string]*models.Results
	if rf, ok := ret.Get(0).(func() map[string]*models.Results); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*models.Results)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChecksResultAndMetadataByCluster provides a mock function with given fields: clusterId
func (_m *MockChecksService) GetChecksResultAndMetadataByCluster(clusterId string) (*models.ClusterCheckResults, error) {
	ret := _m.Called(clusterId)

	var r0 *models.ClusterCheckResults
	if rf, ok := ret.Get(0).(func(string) *models.ClusterCheckResults); ok {
		r0 = rf(clusterId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ClusterCheckResults)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(clusterId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChecksResultByCluster provides a mock function with given fields: clusterId
func (_m *MockChecksService) GetChecksResultByCluster(clusterId string) (*models.Results, error) {
	ret := _m.Called(clusterId)

	var r0 *models.Results
	if rf, ok := ret.Get(0).(func(string) *models.Results); ok {
		r0 = rf(clusterId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Results)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(clusterId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConnectionSettingsById provides a mock function with given fields: id
func (_m *MockChecksService) GetConnectionSettingsById(id string) (map[string]models.ConnectionSettings, error) {
	ret := _m.Called(id)

	var r0 map[string]models.ConnectionSettings
	if rf, ok := ret.Get(0).(func(string) map[string]models.ConnectionSettings); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]models.ConnectionSettings)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConnectionSettingsByNode provides a mock function with given fields: node
func (_m *MockChecksService) GetConnectionSettingsByNode(node string) (models.ConnectionSettings, error) {
	ret := _m.Called(node)

	var r0 models.ConnectionSettings
	if rf, ok := ret.Get(0).(func(string) models.ConnectionSettings); ok {
		r0 = rf(node)
	} else {
		r0 = ret.Get(0).(models.ConnectionSettings)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(node)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSelectedChecksById provides a mock function with given fields: id
func (_m *MockChecksService) GetSelectedChecksById(id string) (models.SelectedChecks, error) {
	ret := _m.Called(id)

	var r0 models.SelectedChecks
	if rf, ok := ret.Get(0).(func(string) models.SelectedChecks); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(models.SelectedChecks)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
