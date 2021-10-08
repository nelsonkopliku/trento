package web

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
	"github.com/trento-project/trento/internal/sapsystem/sapcontrol"

	consulMocks "github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/web/models"
	servicesMocks "github.com/trento-project/trento/web/services/mocks"
)

var sapSystemsList = sapsystem.SAPSystemsList{
	&sapsystem.SAPSystem{
		SID:  "HA1",
		Type: sapsystem.Application,
		Instances: map[string]*sapsystem.SAPInstance{
			"ASCS00": &sapsystem.SAPInstance{
				Host: "netweaver01",
				SAPControl: &sapsystem.SAPControl{
					Properties: map[string]*sapcontrol.InstanceProperty{
						"SAPSYSTEM": &sapcontrol.InstanceProperty{
							Property:     "SAPSYSTEM",
							Propertytype: "string",
							Value:        "00",
						},
					},
					Instances: map[string]*sapcontrol.SAPInstance{
						"netweaver01": &sapcontrol.SAPInstance{
							Hostname:      "netweaver01",
							InstanceNr:    0,
							Features:      "MESSAGESERVER|ENQUEENQREP",
							HttpPort:      50013,
							HttpsPort:     50014,
							StartPriority: "0.5",
							Dispstatus:    "SAPControl-GREEN",
						},
						"netweaver02": &sapcontrol.SAPInstance{
							Hostname:   "netweaver02",
							InstanceNr: 10,
							Features:   "ENQREP",
						},
					},
				},
			},
		},
	},
	&sapsystem.SAPSystem{
		SID:  "HA1",
		Type: sapsystem.Application,
		Instances: map[string]*sapsystem.SAPInstance{
			"ERS10": &sapsystem.SAPInstance{
				Host: "netweaver02",
				SAPControl: &sapsystem.SAPControl{
					Properties: map[string]*sapcontrol.InstanceProperty{
						"SAPSYSTEM": &sapcontrol.InstanceProperty{
							Property:     "SAPSYSTEM",
							Propertytype: "string",
							Value:        "10",
						},
					},
					Instances: map[string]*sapcontrol.SAPInstance{
						"netweaver01": &sapcontrol.SAPInstance{
							Hostname:   "netweaver01",
							InstanceNr: 0,
							Features:   "MESSAGESERVER|ENQUE",
						},
						"netweaver02": &sapcontrol.SAPInstance{
							Hostname:   "netweaver02",
							InstanceNr: 10,
							Features:   "ENQREP",
						},
					},
				},
			},
		},
	},
}

var sapDatabasesList = sapsystem.SAPSystemsList{
	&sapsystem.SAPSystem{
		SID:  "PRD",
		Type: sapsystem.Database,
		Instances: map[string]*sapsystem.SAPInstance{
			"HDB00": &sapsystem.SAPInstance{
				Host: "hana01",
				SAPControl: &sapsystem.SAPControl{
					Properties: map[string]*sapcontrol.InstanceProperty{
						"SAPSYSTEM": &sapcontrol.InstanceProperty{
							Property:     "SAPSYSTEM",
							Propertytype: "string",
							Value:        "00",
						},
					},
					Instances: map[string]*sapcontrol.SAPInstance{
						"hana01": &sapcontrol.SAPInstance{
							Hostname:   "hana01",
							InstanceNr: 0,
							Features:   "HDB_WORKER",
						},
					},
				},
			},
		},
	},
}

func TestSAPSystemsListHandler(t *testing.T) {
	consulInst := new(consulMocks.Client)
	kv := new(consulMocks.KV)
	sapSystemsService := new(servicesMocks.SAPSystemsService)
	hostsService := new(servicesMocks.HostsService)
	sapSystemsService.On("GetSAPSystemsByType", sapsystem.Application).Return(sapSystemsList, nil)
	tagsService := new(servicesMocks.TagsService)
	tagsService.On("GetAllByResource", models.TagSAPSystemResourceType, "HA1").Return([]string{"tag1"}, nil)

	hostsService.On("GetHostMetadata", "netweaver01").Return(map[string]string{
		"trento-ha-cluster":    "netweaver_cluster",
		"trento-ha-cluster-id": "e2f2eb50aef748e586a7baa85e0162cf",
	}, nil)

	hostsService.On("GetHostMetadata", "netweaver02").Return(map[string]string{
		"trento-ha-cluster":    "netweaver_cluster",
		"trento-ha-cluster-id": "e2f2eb50aef748e586a7baa85e0162cf",
	}, nil)

	deps := testDependencies()
	deps.consul = consulInst
	deps.hostsService = hostsService
	deps.sapSystemsService = sapSystemsService
	deps.tagsService = tagsService

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sapsystems", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.webEngine.ServeHTTP(resp, req)

	kv.AssertExpectations(t)
	hostsService.AssertExpectations(t)
	sapSystemsService.AssertExpectations(t)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, responseBody, "SAP Systems")
	assert.Regexp(t, regexp.MustCompile("<td><a href=/sapsystems/HA1>HA1</a></td><td></td><td>.*<input.*value=tag1.*>.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<td>HA1</td><td>MESSAGESERVER|ENQUEENQREP</td><td>00</td><td><a href=/clusters/e2f2eb50aef748e586a7baa85e0162cf>netweaver_cluster</a></td><td><a href=/hosts/netweaver01>netweaver01</a></td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<td>HA1</td><td>ENQREP</td><td>10</td><td><a href=/clusters/e2f2eb50aef748e586a7baa85e0162cf>netweaver_cluster</a></td><td><a href=/hosts/netweaver02>netweaver02</a></td>"), responseBody)
}

func TestSAPDatabaseListHandler(t *testing.T) {
	consulInst := new(consulMocks.Client)
	kv := new(consulMocks.KV)
	sapSystemsService := new(servicesMocks.SAPSystemsService)
	hostsService := new(servicesMocks.HostsService)
	tagsService := new(servicesMocks.TagsService)
	tagsService.On("GetAllByResource", models.TagDatabaseResourceType, "PRD").Return([]string{"tag1"}, nil)

	sapSystemsService.On("GetSAPSystemsByType", sapsystem.Database).Return(sapDatabasesList, nil)

	hostsService.On("GetHostMetadata", "hana01").Return(map[string]string{
		"trento-ha-cluster":    "hana_cluster",
		"trento-ha-cluster-id": "e2f2eb50aef748e586a7baa85e0162cf",
	}, nil)

	deps := testDependencies()
	deps.consul = consulInst
	deps.hostsService = hostsService
	deps.sapSystemsService = sapSystemsService
	deps.tagsService = tagsService

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/databases", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.webEngine.ServeHTTP(resp, req)

	kv.AssertExpectations(t)
	hostsService.AssertExpectations(t)
	sapSystemsService.AssertExpectations(t)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, responseBody, "HANA Databases")
	assert.Regexp(t, regexp.MustCompile("<td><a href=/databases/PRD>PRD</a></td><td></td><td>.*<input.*value=tag1.*>.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<td>PRD</td><td>HDB_WORKER</td><td>00</td><td><a href=/clusters/e2f2eb50aef748e586a7baa85e0162cf>hana_cluster</a></td><td><a href=/hosts/hana01>hana01</a></td>"), responseBody)
}

func TestSAPResourceHandler(t *testing.T) {
	consulInst := new(consulMocks.Client)
	health := new(consulMocks.Health)
	consulInst.On("Health").Return(health)
	sapSystemsService := new(servicesMocks.SAPSystemsService)
	hostsService := new(servicesMocks.HostsService)

	deps := testDependencies()
	deps.consul = consulInst
	deps.sapSystemsService = sapSystemsService
	deps.hostsService = hostsService

	host := hosts.NewHost(consulApi.Node{
		Node:    "netweaver01",
		Address: "192.168.10.10",
		Meta: map[string]string{
			"trento-sap-systems":      "foobar",
			"trento-sap-systems-type": "Application",
			"trento-cloud-provider":   "azure",
			"trento-agent-version":    "0",
			"trento-ha-cluster-id":    "e2f2eb50aef748e586a7baa85e0162cf",
			"trento-ha-cluster":       "banana",
		},
	},
		consulInst)
	hostList := hosts.HostList{
		&host,
	}

	passHealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	health.On("Node", "netweaver01", (*consulApi.QueryOptions)(nil)).Return(passHealthChecks, nil, nil)
	sapSystemsService.On("GetSAPSystemsBySid", "foobar").Return(sapSystemsList, nil)
	hostsService.On("GetHostsBySid", "foobar").Return(hostList, nil)

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sapsystems/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.webEngine.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	responseBody := minifyHtml(resp.Body.String())

	sapSystemsService.AssertExpectations(t)
	hostsService.AssertExpectations(t)
	consulInst.AssertExpectations(t)

	assert.Contains(t, responseBody, "SAP System details")
	assert.Contains(t, responseBody, "foobar")
	// Layout
	assert.Regexp(t, regexp.MustCompile("<tr><td>netweaver01</td><td>00</td><td>MESSAGESERVER|ENQUEENQREP</td><td>50013</td><td>50014</td><td>0.5</td><td><span.*primary.*>SAPControl-GREEN</span></td></tr>"), responseBody)
	// Host
	assert.Regexp(t, regexp.MustCompile("<tr><td>.*check_circle.*</td><td><a href=/hosts/netweaver01>netweaver01</a></td><td>192.168.10.10</td><td>azure</td><td><a href=/clusters/e2f2eb50aef748e586a7baa85e0162cf>banana</a></td><td><a href=/sapsystems/foobar>foobar</a></td><td>v0</td></tr>"), responseBody)
}

func TestSAPResourceHandler404Error(t *testing.T) {
	sapSystemsService := new(servicesMocks.SAPSystemsService)

	deps := testDependencies()
	deps.sapSystemsService = sapSystemsService

	sapSystemsService.On("GetSAPSystemsBySid", "foobar").Return(sapsystem.SAPSystemsList{}, nil)

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sapsystems/foobar", nil)
	req.Header.Set("Accept", "text/html")

	app.webEngine.ServeHTTP(resp, req)

	sapSystemsService.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, resp.Body.String(), "Not Found")
}

func minifyHtml(input string) string {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", input)
	if err != nil {
		panic(err)
	}
	return minified
}
