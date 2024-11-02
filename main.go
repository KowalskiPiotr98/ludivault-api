package main

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/database"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	log.SetOutput(os.Stdout)
	//todo: make level configurable
	log.SetLevel(log.InfoLevel)

	operations.Errors.RegisterDefaultPostgresHandlers()
}

func main() {
	log.Debugln("Initialising database connection")
	//todo: make this configurable
	if err := gotabase.InitialiseConnection("user=postgres dbname=ludivault password=postgres sslmode=disable", "postgres"); err != nil {
		log.Panicf("Failed to initialise database connection: %v", err)
	}
	defer gotabase.CloseConnection()
	log.Debugln("Applying database migrations...")
	if err := database.RunMigrations(gotabase.GetConnection()); err != nil {
		log.Panicf("Failed to apply database migrations: %v", err)
	}

	router := setupEngine()

	log.Infoln("Starting server...")
	//todo: configurable address
	if err := router.Run("localhost:5500"); err != nil {
		log.Panicf("Server failed while listening: %v", err)
	}
}

func setupEngine() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(getLogger())
	//todo: trusted proxies

	//todo: base path config
	//basePath := ""

	return router
}
