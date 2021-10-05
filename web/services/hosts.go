package services

import (
	"fmt"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
)

//go:generate mockery --name=HostsService

type HostsService interface {
	GetHostMetadata(host string) (map[string]string, error)
	GetHostsBySid(sid string) (hosts.HostList, error)
}

type hostsService struct {
	consul consul.Client
}

func NewHostsService(client consul.Client) HostsService {
	return &hostsService{consul: client}
}

func (h *hostsService) GetHostMetadata(host string) (map[string]string, error) {
	hostList, err := hosts.Load(h.consul, fmt.Sprintf("Node == %s", host), nil)
	if err != nil {
		return nil, err
	}

	if len(hostList) == 0 {
		return nil, fmt.Errorf("host with name %s not found", host)
	}

	return hostList[0].TrentoMeta(), nil
}

func (h *hostsService) GetHostsBySid(sid string) (hosts.HostList, error) {
	hostList, err := hosts.Load(h.consul, fmt.Sprintf("Meta[\"trento-sap-systems\"] == %s", sid), nil)
	if err != nil {
		return nil, err
	}

	return hostList, nil
}
