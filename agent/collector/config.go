package collector

type CollectorConfig struct {
	Enabled bool
	Host    string
	TLS     struct {
		ClientCert string
		ClientKey  string
		CACert     string
	}
}
