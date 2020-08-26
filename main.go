package main

import (
	"context"
	"log"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/dataprovider"
	"github.com/tidepool-org/summary/server"
	"github.com/tidepool-org/summary/store"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// ServerTimeoutAmount is the amount of time before we time out the server
	ServerTimeoutAmount = 20
	_                   = openapi3filter.Options{}
)

//ServiceConfig the configuration for the summary service
type ServiceConfig struct {
	ServiceAuth string `envconfig:"TIDEPOOL_SUMMARY_SERVICE_SECRET" required:"false"`
	Address     string `envconfig:"TIDEPOOL_SUMMARY_SERVICE_SERVER_ADDRESS" default:":8080"`
}

//NewServiceConfigFromEnv creates a service config
func NewServiceConfigFromEnv() *ServiceConfig {
	var config ServiceConfig
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &config
}

//main is the main loop
func main() {

	config := NewServiceConfigFromEnv()

	// Echo instance
	e := echo.New()
	e.Logger.Print("Starting Main Loop")
	swagger, err := api.GetSwagger()
	if err != nil {
		e.Logger.Fatal("Cound not get spec")
	}

	// Middleware
	//authClient := AuthClient{store: dbstore}
	//filterOptions := openapi3filter.Options{AuthenticationFunc: authClient.AuthenticationFunc}
	//options := Options{Options: filterOptions}

	loggerConfig := middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/status")
		},
	}
	e.Use(middleware.LoggerWithConfig(loggerConfig))

	e.GET("/status", hello)

	e.Use(middleware.Recover())

	options := api.Options{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/status")
		},
	}
	e.Use(api.OapiRequestValidatorWithOptions(swagger, &options))

	uriProvider := store.NewMongoURIProviderFromEnv()
	client, err := mongo.NewClient(mongoOptions.Client().ApplyURI(uriProvider.URI()))
	if err != nil {
		log.Fatalln("NewMongoStoreClient: cannot create client:", err)
	}

	client.Connect(context.Background())
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("could not ping client %v", err)
	}
	mongoProvider := dataprovider.NewMongoProvider(client)
	sharerProvider := dataprovider.NewMongoShareProvider(client)
	summaryServer := server.NewSummaryServer(mongoProvider, sharerProvider)

	// Register Handler
	api.RegisterHandlers(e, summaryServer)

	// Start server
	e.Logger.Printf("Starting Server at: %s\n", config.Address)
	go func() {
		if err := e.Start(config.Address); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ServerTimeoutAmount)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
