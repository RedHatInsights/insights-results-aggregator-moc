/*
Copyright © 2020, 2023 Red Hat, Inc.

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

// Implementation of REST API tests that checks all REST API endpoints of
// Insights aggregator mock service.
//
// These test should be started by using one of following commands in order to be configured properly:
//
//	./run_on_ci.sh
//	./test.sh
//
// It is also possible to run REST API tests with code coverage detection:
//
//	./rest-api-tests.sh
//
// REST API endpoints that are tested:
//
// URL                                         handler                                     HTTP methods
// apiPrefix+MainEndpoint                      server.mainEndpoint                         GET
// apiPrefix+GroupsEndpoint                    server.listOfGroups                         GET   OPTIONS
// apiPrefix+ContentEndpoint                   server.serveContentWithGroups               GET   OPTIONS
// apiPrefix+RuleClusterDetailEndpoint         server.ruleClusterDetailEndpoint            GET
// apiPrefix+OrganizationsEndpoint             server.listOfOrganizations                  GET
// apiPrefix+ClustersForOrganizationEndpoint   server.listOfClustersForOrganization        GET
// apiPrefix+ClustersEndpoint                  server.readReportForClusters                GET   POST   OPTIONS
// apiPrefix+ReportForClusterEndpoint          server.readReportForCluster                 GET   OPTIONS
// apiPrefix+ReportForClusterEndpoint2         server.readReportForCluster                 GET   OPTIONS
// apiPrefix+ReportEndpoint                    server.readReportForOrganizationAndCluster  GET   OPTIONS
// apiPrefix+ClustersInOrgEndpoint             server.readReportForAllClustersInOrg        GET
// apiPrefix+AckListEndpoint                   server.readAckList                          GET
// apiPrefix+AckAcknowledgePostEndpoint        server.acknowledgePost                      POST
// apiPrefix+AckGetEndpoint                    server.getAcknowledge                       GET
// apiPrefix+AckUpdateEndpoint                 server.updateAcknowledge                    PUT
// apiPrefix+AckDeleteEndpoint                 server.deleteAcknowledge                    DELETE
// apiPrefix+ListAllRequestIDs                 server.readListOfRequestIDs                 GET
// apiPrefix+ListAllRequestIDs                 server.readListOfRequestIDsPostVariant      POST
// apiPrefix+StatusOfRequestID                 server.readStatusOfRequestID                GET
// apiPrefix+RuleHitsForRequestID              server.readRuleHitsForRequestID             GET
// apiPrefix+AllDVONamespaces                  server.allDVONamespaces                     GET
// apiPrefix+UpgradeRisksPredictionEndpoint    server.upgradeRisksPrediction               GET
// openAPIURL                                  server.serveAPISpecFile                     GET
package main

import (
	"os"

	"github.com/verdverm/frisby"

	tests "github.com/RedHatInsights/insights-results-aggregator-mock/tests/rest"
)

func main() {
	tests.ServerTests()
	frisby.Global.PrintReport()
	os.Exit(frisby.Global.NumErrored)
}
