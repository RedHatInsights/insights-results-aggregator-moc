/*
Copyright © 2020, 2021, 2022, 2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package server contains implementation of REST API server (HTTPServer) for the
// Insights content service. In current version, the following
// REST API endpoints are available:
package server

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	// we just have to import this package in order to expose pprof interface in debug mode
	// disable "G108 (CWE-): Profiling endpoint is automatically exposed on /debug/pprof"
	_ "net/http/pprof" // #nosec G108
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/insights-results-aggregator-mock/content"
	"github.com/RedHatInsights/insights-results-aggregator-mock/groups"
	"github.com/RedHatInsights/insights-results-aggregator-mock/storage"
)

// HTTPServer in an implementation of Server interface
type HTTPServer struct {
	Config     Configuration
	Storage    storage.Storage
	Groups     map[string]groups.Group
	Serv       *http.Server
	groupsList []groups.Group
	Content    []content.RuleContent
}

// New constructs new implementation of Server interface
func New(config Configuration,
	storageInstance storage.Storage,
	ruleGroups map[string]groups.Group,
	ruleContents []content.RuleContent) *HTTPServer {
	return &HTTPServer{
		Config:  config,
		Storage: storageInstance,
		Groups:  ruleGroups,
		Content: ruleContents,
	}
}

// Start starts server
func (server *HTTPServer) Start() error {
	address := server.Config.Address
	log.Info().Msgf("Starting HTTP server at '%s'", address)
	router := server.Initialize(address)
	server.Serv = &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	server.printAccessInfo()

	err := server.Serv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Error().Err(err).Msg("Unable to start HTTP/S server")
		return err
	}

	return nil
}

func (server *HTTPServer) printAccessInfo() {
	// access command should look like:
	// curl localhost:8080/api/insights-results-aggregator/v2/
	address := server.Config.Address
	apiPrefix := server.Config.APIPrefix

	hostname, err := os.Hostname()
	if err != nil {
		log.Error().Err(err).Msg("Cannot retrieve hostname")
		hostname = "localhost"
	}

	log.Info().Msgf("Access REST API via: curl %s%s%s", hostname, address, apiPrefix)
}

// Stop stops server's execution
func (server *HTTPServer) Stop(ctx context.Context) error {
	return server.Serv.Shutdown(ctx)
}

// Initialize perform the server initialization
func (server *HTTPServer) Initialize(address string) http.Handler {
	log.Info().Msgf("Initializing HTTP server at '%s'", address)

	router := mux.NewRouter().StrictSlash(true)

	server.addEndpointsToRouter(router)

	// Endpoints enabled in Debug mode only
	if server.Config.Debug {
		log.Info().Msg("Debug endpoints enabled")
		server.addDebugEndpointsToRouter(router)
	}

	log.Info().Msgf("Server has been initiliazed")

	return router
}

func (server *HTTPServer) addEndpointsToRouter(router *mux.Router) {
	apiPrefix := server.Config.APIPrefix
	if !strings.HasSuffix(apiPrefix, "/") {
		apiPrefix += "/"
	}
	log.Info().Msgf("API prefix is set to '%s'", apiPrefix)

	openAPIURL := apiPrefix + filepath.Base(server.Config.APISpecFile)

	// common REST API endpoints
	router.HandleFunc(apiPrefix+MainEndpoint, server.mainEndpoint).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+GroupsEndpoint, server.listOfGroups).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(apiPrefix+ContentEndpoint, server.serveContentWithGroups).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(apiPrefix+InfoEndpoint, server.serviceInfo).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc(apiPrefix+OrganizationsEndpoint, server.listOfOrganizations).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+ClustersForOrganizationEndpoint, server.listOfClustersForOrganization).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+ReportEndpoint, server.readReportForOrganizationAndCluster).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(apiPrefix+ReportForClusterEndpoint, server.readReportForCluster).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(apiPrefix+ReportForClusterEndpoint2, server.readReportForCluster).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(apiPrefix+ClustersEndpoint, server.readReportForClusters).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(apiPrefix+ClustersInOrgEndpoint, server.readReportForAllClustersInOrg).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+RuleClusterDetailEndpoint, server.ruleClusterDetailEndpoint).Methods(http.MethodGet)

	// Endpoints to manipulate with simplified rule results stored
	// independently under "tracker_id" identifier
	router.HandleFunc(apiPrefix+ListAllRequestIDs, server.readListOfRequestIDs).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+ListAllRequestIDs, server.readListOfRequestIDsPostVariant).Methods(http.MethodPost)
	router.HandleFunc(apiPrefix+StatusOfRequestID, server.readStatusOfRequestID).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+RuleHitsForRequestID, server.readRuleHitsForRequestID).Methods(http.MethodGet)

	// Acknowledgement-related endpoints. Please look into acks_handlers.go
	// and acks_utils.go for more information about these endpoints
	// prepared to be compatible with RHEL Insights Advisor.
	router.HandleFunc(apiPrefix+AckListEndpoint, server.readAckList).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+AckAcknowledgePostEndpoint, server.acknowledgePost).Methods(http.MethodPost)
	router.HandleFunc(apiPrefix+AckGetEndpoint, server.getAcknowledge).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+AckUpdateEndpoint, server.updateAcknowledge).Methods(http.MethodPut)
	router.HandleFunc(apiPrefix+AckDeleteEndpoint, server.deleteAcknowledge).Methods(http.MethodDelete)

	// Upgrade risks prediction endpoints. Please look into upgrade_risks_prediction.go
	// for more information about this endpoint
	router.HandleFunc(apiPrefix+UpgradeRisksPredictionEndpoint, server.upgradeRisksPrediction).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+UpgradeRisksPredictionMultiClusterEndpoint, server.upgradeRisksPredictionMultiCluster).Methods(http.MethodPost)

	// DVO-related endpoints:
	//
	// AllDVONamespaces = "namespaces/dvo"
	// DVONamespaceForCluster1 = "cluster/{cluster_name}/namespaces/dvo"
	// DVONamespaceForCluster2 = "namespaces/dvo/cluster/{cluster_name}"
	// DVONamespaceInfo = "namespaces/dvo/{namespace_id}/info"
	// DVONamespaceReports = "namespaces/dvo/{namespace_id}/reports"
	router.HandleFunc(apiPrefix+AllDVONamespaces, server.allDVONamespaces).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+DVONamespaceForCluster1, server.dvoNamespaceForCluster).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+DVONamespaceForCluster2, server.dvoNamespaceForCluster).Methods(http.MethodGet)

	// OpenAPI specs
	router.HandleFunc(openAPIURL, server.serveAPISpecFile).Methods(http.MethodGet)
}

func (server *HTTPServer) addDebugEndpointsToRouter(router *mux.Router) {
	apiPrefix := server.Config.APIPrefix
	if !strings.HasSuffix(apiPrefix, "/") {
		apiPrefix += "/"
	}

	router.HandleFunc(apiPrefix+ExitEndpoint, server.exit).Methods(http.MethodPut)
}

/*
// addCORSHeaders - middleware for adding headers that should be in any response
func (server *HTTPServer) addCORSHeaders(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			nextHandler.ServeHTTP(w, r)
		})
}
*/

/*
// handleOptionsMethod - middleware for handling OPTIONS method
func (server *HTTPServer) handleOptionsMethod(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
			} else {
				nextHandler.ServeHTTP(w, r)
			}
		})
}
*/
