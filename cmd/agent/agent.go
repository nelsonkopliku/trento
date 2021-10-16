package agent

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/trento-project/trento/agent"
	"github.com/trento-project/trento/agent/collector"
	"github.com/trento-project/trento/internal/cluster"
)

var consulConfigDir string
var discoveryPeriod int

var collectorHost string
var cert string
var key string
var ca string

func NewAgentCmd() *cobra.Command {

	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Command tree related to the agent component",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the agent",
		Run:   start,
	}
	startCmd.Flags().StringVarP(&consulConfigDir, "consul-config-dir", "", "consul.d", "Consul configuration directory used to store node meta-data")
	startCmd.Flags().IntVarP(&discoveryPeriod, "discovery-period", "", 2, "Discovery mechanism loop period on minutes")

	startCmd.Flags().StringVar(&collectorHost, "collector-host", "https://localhost:8443", "Data Collector endpoint")
	startCmd.Flags().StringVar(&cert, "collector-client-cert", "", "mTLS client certificate")
	startCmd.Flags().StringVar(&key, "collector-client-key", "", "mTLS client key")
	startCmd.Flags().StringVar(&ca, "collector-ca", "", "mTLS Certificate Authority")

	// Bind the flags to viper and make them available to the application
	startCmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	agentCmd.AddCommand(startCmd)

	return agentCmd
}

func start(cmd *cobra.Command, args []string) {
	var err error

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := agent.DefaultConfig()
	if err != nil {
		log.Fatal("Failed to create the agent configuration: ", err)
	}

	cfg.ConsulConfigDir = consulConfigDir
	cfg.DiscoveryPeriod = time.Duration(discoveryPeriod) * time.Minute
	cfg.ClusterDiscoveryOptions.CollectorConfig = extractCollectorConnectionOptions()
	cfg.ClusterDiscoveryOptions.ClusterDiscoverytools = extractClusterDiscoveryTools()

	a, err := agent.NewWithConfig(cfg)
	if err != nil {
		log.Fatal("Failed to create the agent: ", err)
	}

	go func() {
		quit := <-signals
		log.Printf("Caught %s signal!", quit)

		log.Println("Stopping the agent...")
		a.Stop()
	}()

	log.Println("Starting the Console Agent...")
	err = a.Start()
	if err != nil {
		log.Fatal("Failed to start the agent: ", err)
	}
}

func extractCollectorConnectionOptions() collector.CollectorConfig {
	return collector.NewCollectorConfig(
		viper.GetString("collector-host"),
		collector.TlsConfig{
			ClientCert: viper.GetString("collector-client-cert"),
			ClientKey:  viper.GetString("collector-client-key"),
			CACert:     viper.GetString("collector-ca"),
		},
	)
}

func extractClusterDiscoveryTools() cluster.DiscoveryTools {
	return cluster.DiscoveryTools{
		CibAdmPath:      viper.GetString("ClusterDiscovery.Tools.cibAdminPath"),
		CrmmonAdmPath:   viper.GetString("ClusterDiscovery.Tools.crmmonAdmPath"),
		CorosyncKeyPath: viper.GetString("ClusterDiscovery.Tools.corosyncKeyPath"),
		SBDPath:         viper.GetString("ClusterDiscovery.Tools.SBDPath"),
		SBDConfigPath:   viper.GetString("ClusterDiscovery.Tools.SBDConfigPath"),
	}
}
