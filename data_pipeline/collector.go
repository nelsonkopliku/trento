package data_pipeline

// import (
// 	"crypto/tls"
// 	"crypto/x509"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gin-contrib/sessions"
// 	"github.com/spf13/viper"
// 	"github.com/swaggo/gin-swagger/swaggerFiles"
// )

// type Collector struct {
// 	projectorRegistry string
// 	eventStore        string
// }

// func NewCollector() (*Collector, error) {
// 	app := &App{
// 		Dependencies: deps,
// 		host:         host,
// 		port:         port,
// 	}

// 	InitAlerts()
// 	engine := deps.engine
// 	engine.HTMLRender = NewLayoutRender(templatesFS, "templates/*.tmpl")
// 	engine.Use(ErrorHandler)
// 	engine.Use(sessions.Sessions("session", deps.store))
// 	engine.StaticFS("/static", http.FS(assetsFS))
// 	engine.GET("/", HomeHandler)
// 	engine.GET("/about", NewAboutHandler(deps.subscriptionsService))
// 	engine.GET("/hosts", NewHostListHandler(deps.consul, deps.tagsService))
// 	engine.GET("/hosts/:name", NewHostHandler(deps.consul, deps.subscriptionsService))
// 	engine.GET("/catalog", NewChecksCatalogHandler(deps.checksService))
// 	engine.GET("/clusters", NewClusterListHandler(deps.consul, deps.checksService, deps.tagsService))
// 	engine.GET("/clusters/:id", NewClusterHandler(deps.consul, deps.checksService))
// 	engine.POST("/clusters/:id/settings", NewSaveClusterSettingsHandler(deps.consul))
// 	engine.GET("/sapsystems", NewSAPSystemListHandler(deps.consul, deps.hostsService, deps.sapSystemsService, deps.tagsService))
// 	engine.GET("/sapsystems/:sid", NewSAPResourceHandler(deps.hostsService, deps.sapSystemsService))
// 	engine.GET("/databases", NewHanaDatabaseListHandler(deps.consul, deps.hostsService, deps.sapSystemsService, deps.tagsService))
// 	engine.GET("/databases/:sid", NewSAPResourceHandler(deps.hostsService, deps.sapSystemsService))
// 	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	apiGroup := engine.Group("/api")
// 	{
// 		apiGroup.GET("/ping", ApiPingHandler)

// 		apiGroup.GET("/tags", ApiListTag(deps.tagsService))
// 		apiGroup.POST("/hosts/:name/tags", ApiHostCreateTagHandler(deps.consul, deps.tagsService))
// 		apiGroup.DELETE("/hosts/:name/tags/:tag", ApiHostDeleteTagHandler(deps.consul, deps.tagsService))
// 		apiGroup.POST("/clusters/:id/tags", ApiClusterCreateTagHandler(deps.consul, deps.tagsService))
// 		apiGroup.DELETE("/clusters/:id/tags/:tag", ApiClusterDeleteTagHandler(deps.consul, deps.tagsService))
// 		apiGroup.GET("/clusters/:cluster_id/results", ApiClusterCheckResultsHandler(deps.consul, deps.checksService))
// 		apiGroup.POST("/sapsystems/:sid/tags", ApiSAPSystemCreateTagHandler(deps.consul, deps.tagsService))
// 		apiGroup.DELETE("/sapsystems/:sid/tags/:tag", ApiSAPSystemDeleteTagHandler(deps.consul, deps.tagsService))
// 	}

// 	return app, nil
// }

// func (c *Collector) Start() error {
// 	// Create a CA certificate pool and add cert.pem to it
// 	caCert, err := ioutil.ReadFile("cert.pem")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)

// 	// Create the TLS Config with the CA pool and enable Client certificate validation
// 	tlsConfig := &tls.Config{
// 		ClientCAs:  caCertPool,
// 		ClientAuth: tls.RequireAndVerifyClientCert,
// 	}
// 	// tlsConfig.BuildNameToCertificate()
// 	s2 := &http.Server{
// 		Addr:           fmt.Sprintf("%s:%d", a.host, 8443),
// 		Handler:        a,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 		TLSConfig:      tlsConfig,
// 	}

// 	if viper.GetBool("disable-mtls") {
// 		err = s2.ListenAndServe()
// 	} else {
// 		err = s2.ListenAndServeTLS("cert.pem", "key.pem")
// 	}

// 	return err
// }
