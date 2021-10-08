package readmodels

const (
	ClusterTypeScaleUp  = "HANA scale-up"
	ClusterTypeScaleOut = "HANA scale-out"
	ClusterTypeUnknown  = "Unknown"
)

type Cluster struct {
	ID                string
	Name              string
	ClusterType       string
	SID               string `gorm:"column:sid"`
	ResourcesNumber   int
	HostsNumber       int
	Health            string   `gorm:"-"`
	Tags              []string `gorm:"-"`
	HasDuplicatedName bool     `gorm:"-"`
}

type ClusterList []*Cluster

func (clusterList ClusterList) GetAllSIDs() []string {
	var sids []string
	set := make(map[string]struct{})

	for _, c := range clusterList {
		_, ok := set[c.SID]
		if !ok {
			set[c.SID] = struct{}{}
			sids = append(sids, c.SID)
		}
	}

	return sids
}

func (clusterList ClusterList) GetAllTags() []string {
	var tags []string
	set := make(map[string]struct{})

	for _, c := range clusterList {
		for _, tag := range c.Tags {
			_, ok := set[tag]
			if !ok {
				set[tag] = struct{}{}
				tags = append(tags, tag)
			}
		}
	}

	return tags
}

func (clusterList ClusterList) GetAllClusterTypes() []string {
	var clusterTypes []string
	set := make(map[string]struct{})

	for _, c := range clusterList {
		_, ok := set[c.ClusterType]
		if !ok {
			set[c.ClusterType] = struct{}{}
			clusterTypes = append(clusterTypes, c.ClusterType)
		}
	}

	return clusterTypes
}
