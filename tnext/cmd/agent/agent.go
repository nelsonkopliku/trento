package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/tnext/app/agent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var discoveryPeriod int

var collectorHost string
var collectorPort int

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
	startCmd.Flags().IntVarP(&discoveryPeriod, "discovery-period", "", 2, "Discovery mechanism loop period on minutes")

	startCmd.Flags().StringVar(&collectorHost, "collector-host", "localhost", "Data Collector host")
	startCmd.Flags().IntVar(&collectorPort, "collector-port", 8081, "Data Collector port")

	// Bind the flags to viper and make them available to the application
	startCmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	agentCmd.AddCommand(startCmd)

	return agentCmd
}

func start(cmd *cobra.Command, args []string) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		log.Infof("Bye!")
		cancel()
	}()

	var err error

	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to create the agent configuration: ", err)
	}

	a, err := agent.NewAgent(config)
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
