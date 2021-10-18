package collector

type CollectorConfig struct {
	Enabled bool
	Host    string
	Port    int
	TLS     TlsConfig
}

type TlsConfig struct {
	ClientCert string
	ClientKey  string
	CACert     string
}

func NewCollectorConfig(host string, port int, tls TlsConfig) CollectorConfig {
	collectorConfig := CollectorConfig{}
	collectorConfig.Enabled = host != ""
	collectorConfig.Host = host
	collectorConfig.Port = port
	collectorConfig.TLS = tls
	return collectorConfig
}
