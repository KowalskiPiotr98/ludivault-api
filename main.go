package main

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/auth"
	"github.com/KowalskiPiotr98/ludivault/controllers"
	"github.com/KowalskiPiotr98/ludivault/database"
	"github.com/KowalskiPiotr98/ludivault/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	log.SetOutput(os.Stdout)
	//todo: make level configurable
	if gin.Mode() == gin.ReleaseMode {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	operations.Errors.RegisterDefaultPostgresHandlers()
}

func main() {
	log.Debugln("Initialising database connection")
	if err := gotabase.InitialiseConnection(utils.GetRequiredConfig("db"), "postgres"); err != nil {
		log.Panicf("Failed to initialise database connection: %v", err)
	}
	defer gotabase.CloseConnection()
	log.Debugln("Applying database migrations...")
	if err := database.RunMigrations(gotabase.GetConnection()); err != nil {
		log.Panicf("Failed to apply database migrations: %v", err)
	}

	if err := runEngine(); err != nil {
		log.Panicf("Server failed while listening: %v", err)
	}
}

func runEngine() error {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(getLogger())
	router.Use(auth.GetUserMiddleware())
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"192.168.0.0/16", "10.0.0.0/8", "172.16.0.0/12"})

	listenAddress := utils.GetOptionalConfig("listen", "localhost:5500")
	baseAddress, err := url.Parse(utils.GetRequiredConfig("base_address"))
	if err != nil {
		log.Panicf("Failed to parse base address: %v", err)
	}
	listenDomain := baseAddress.Hostname()

	controllers.SetRoutes(router.Group("/api/v1"))

	auth.InitSessionStore(listenDomain)
	if err = auth.SetupProviders(baseAddress.String()); err != nil {
		log.Panicf("Failed to setup login providers: %v", err)
	}

	log.Infoln("Starting server...")
	return router.Run(listenAddress)
}
