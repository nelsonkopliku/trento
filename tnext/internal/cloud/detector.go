package cloud

type CloudInstance struct {
	Provider string      `mapstructure:"provider,omitempty"`
	Metadata interface{} `mapstructure:"metadata,omitempty"`
}

const (
	Azure   = "azure"
	AWS     = "aws"
	GCP     = "gcp"
	UNKNOWN = "unknown"
)

type Detector interface {
	Detect() (*CloudInstance, error)
}

type detector struct {
	detection []Detector
}

func NewCloudDetector(detectors ...Detector) Detector {
	return &detector{detectors}
}

func (detector *detector) Detect() (*CloudInstance, error) {
	for _, detector := range detector.detection {
		if csp, err := detector.Detect(); err == nil {
			return csp, nil
		}
	}
	return &CloudInstance{
		Provider: string(UNKNOWN),
	}, nil
}

// identification and metadata loading

type Identifier interface {
	Identify() (string, error)
}

type MetadataLoader interface {
	Load(*CloudInstance) error
}

type MetadataLoadingDetector struct {
	identifier     Identifier
	metadataLoader MetadataLoader
}

func NewMetadataLoadingDetector(identifier Identifier, metadataLoader MetadataLoader) Detector {
	return &MetadataLoadingDetector{
		identifier:     identifier,
		metadataLoader: metadataLoader,
	}
}

func (detector *MetadataLoadingDetector) Detect() (*CloudInstance, error) {
	var detectedCloudInstance *CloudInstance

	csp, err := detector.identifier.Identify()
	if err != nil {
		return detectedCloudInstance, err
	}

	detectedCloudInstance.Provider = csp

	if err := detector.metadataLoader.Load(detectedCloudInstance); err != nil {
		return detectedCloudInstance, err
	}

	return detectedCloudInstance, nil
}
