package rosetta

import (
	"fmt"
	"net/http"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/services"
	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

var log = logger.GetOrCreate("rosetta")

// CreateServer creates a HTTP server
func CreateServer(elrondFacade api.ElrondProxyHandler, generalConfig *config.Config, port int) (*http.Server, error) {
	elrondClient, err := client.NewElrondClient(elrondFacade)
	if err != nil {
		log.Error("cannot create elrond client", "err", err)
		return nil, err
	}

	networkConfig, err := elrondClient.GetNetworkConfig()
	if err != nil {
		log.Error("cannot get network config", "err", err)
		return nil, err
	}

	cfg := configuration.LoadConfiguration(networkConfig, generalConfig)

	// The asserter automatically rejects incorrectly formatted
	// requests.
	asserterServer, err := asserter.NewServer(
		services.SupportedOperationTypes,
		false,
		[]*types.NetworkIdentifier{
			cfg.Network,
		},
	)
	if err != nil {
		log.Error("cannot create asserter", "err", err)
		return nil, err
	}

	// Create network service
	networkAPIService := services.NewNetworkAPIService(elrondClient, cfg)
	networkAPIController := server.NewNetworkAPIController(
		networkAPIService,
		asserterServer,
	)

	// Create account service
	accountAPIService := services.NewAccountAPIService(elrondClient, cfg)
	accountAPIController := server.NewAccountAPIController(
		accountAPIService,
		asserterServer,
	)

	// Create block service
	blockAPIService := services.NewBlockAPIService(elrondClient, cfg, networkConfig)
	blockAPIController := server.NewBlockAPIController(
		blockAPIService,
		asserterServer,
	)

	// Create construction service
	constructionAPIService := services.NewConstructionAPIService(elrondClient, cfg, networkConfig)
	constructionAPIController := server.NewConstructionAPIController(
		constructionAPIService,
		asserterServer,
	)

	// Create mempool service
	mempoolAPIService := services.NewMempoolApiService(elrondClient, cfg, networkConfig)
	mempoolAPIController := server.NewMempoolAPIController(
		mempoolAPIService,
		asserterServer,
	)

	router := server.NewRouter(
		networkAPIController,
		accountAPIController,
		blockAPIController,
		constructionAPIController,
		mempoolAPIController,
	)

	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: corsRouter,
	}

	return httpServer, nil
}
