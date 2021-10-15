package collector

type CollectorConfig struct {
	Enabled bool
	Host    string
	TLS     TlsConfig
}

type TlsConfig struct {
	ClientCert string
	ClientKey  string
	CACert     string
}

func NewCollectorConfig(host string, tls TlsConfig) CollectorConfig {
	collectorConfig := CollectorConfig{}
	collectorConfig.Enabled = host != ""
	collectorConfig.Host = host
	collectorConfig.TLS = tls
	return collectorConfig
}
