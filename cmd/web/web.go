package web

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/trento-project/trento/web"
)

var host string
var port int
var araAddr string

var dbhost string
var dbport string
var dbuser string
var dbpassword string
var dbname string

func NewWebCmd() *cobra.Command {
	webCmd := &cobra.Command{
		Use:   "web",
		Short: "Command tree related to the web application component",
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the web application",
		Run:   serve,
	}

	serveCmd.Flags().StringVar(&host, "host", "0.0.0.0", "The host to bind the HTTP service to")
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port for the HTTP service to listen at")
	serveCmd.Flags().StringVar(&araAddr, "ara-addr", "127.0.0.1:8000", "Address where ARA is running (ex: localhost:80)")

	serveCmd.Flags().StringVar(&dbhost, "dbhost", "localhost", "The database host")
	serveCmd.Flags().StringVar(&dbport, "dbport", "5432", "The database port to connect to")
	serveCmd.Flags().StringVar(&dbuser, "dbuser", "postgres", "The database user")
	serveCmd.Flags().StringVar(&dbpassword, "dbpassword", "postgres", "The database password")
	serveCmd.Flags().StringVar(&dbname, "dbname", "trento", "The database name that the application will usee")

	// Bind the flags to viper and make them available to the application
	serveCmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Needed if the flags contain dashes
		if strings.Contains(f.Name, "-") {
			viper.BindEnv(f.Name, strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_")))
		}

		viper.BindPFlag(f.Name, f)
	})

	webCmd.AddCommand(serveCmd)

	return webCmd
}

func serve(cmd *cobra.Command, args []string) {
	var err error

	deps := web.DefaultDependencies()

	app, err := web.NewAppWithDeps(host, port, deps)
	if err != nil {
		log.Fatal("Failed to create the web application instance: ", err)
	}

	err = app.Start()
	if err != nil {
		log.Fatal("Failed to start the web application service: ", err)
	}
}
