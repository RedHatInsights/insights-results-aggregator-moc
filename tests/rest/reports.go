/*
Copyright © 2023 Red Hat, Inc.

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

// Package tests contains REST API tests for following endpoints:
//
// apiPrefix
// apiPrefix + "clusters"
// apiPrefix + "groups"
// apiPrefix + "organizations"
package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/verdverm/frisby"

	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

// FullReportResponse represents response containing report for given cluster
type FullReportResponse struct {
	Report types.ReportResponse `json:"report"`
	Status string               `json:"status"`
}

// AllReportsForOrganizationResponse represents response with all reports for
// given organization
type AllReportsForOrganizationResponse struct {
	Clusters    []string    `json:"clusters"`
	Errors      []string    `json:"erors"`
	Reports     interface{} `json:"reports"`
	GeneratedAt time.Time   `json:"generated_at"` // timestamp
}

// reportEndpointForCluster helper function constructs URL for accessing endpoint to
// retrieve report for given cluster (w/o organization ID)
func reportEndpointForCluster(clusterName string) string {
	return fmt.Sprintf("%sreport/%s", apiURL, clusterName)
}

// reportEndpointForOrgAndCluster helper function constructs URL for accessing endpoint to
// retrieve report for given organization and cluster
func reportEndpointForOrgAndCluster(orgID int, clusterName string) string {
	return fmt.Sprintf("%sreport/%d/%s", apiURL, orgID, clusterName)
}

// reportEndpointForAllReportsForOrg helper function constructs URL for accessing endpoint to
// retrieve all reports for given organization
func reportEndpointForAllReportsForOrg(orgID int) string {
	return fmt.Sprintf("%sclusters/%d", apiURL, orgID)
}

// checkReportForKnownOrganizationKnownCluster checks if proper report is returned for
// known organization ID and known cluster name
func checkReportForKnownOrganizationKnownCluster() {
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with known cluster").Get(reportEndpointForOrgAndCluster(organization1, cluster1ForOrg1))
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := FullReportResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError(statusShouldBeSetToOK)
		}
	}
	f.PrintReport()
}

// checkReportForUknownOrganization checks how uknown organization ID is
// checked by REST API handler
func checkReportForUnknownOrganization() {
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with unknown organization").Get(reportEndpointForOrgAndCluster(1234, unknownCluster))
	f.Send()
	f.ExpectStatus(http.StatusNotFound)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	f.PrintReport()
}

// checkReportForImproperOrganization checks how improperly entered
// organization ID is checked by REST API handler
func checkReportForImproperOrganization() {
	url := fmt.Sprintf("%sreport/foobar/%s", apiURL, cluster1ForOrg1)
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with improper organization ID").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusBadRequest)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	f.PrintReport()
}

// checkReportForKnownOrganizationUnknownCluster checks how unknown cluster
// name is checked by REST API handler
func checkReportForKnownOrganizationUnknownCluster() {
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with unknown cluster").Get(reportEndpointForOrgAndCluster(organization1, unknownCluster))
	f.Send()
	f.ExpectStatus(http.StatusNotFound)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	f.PrintReport()
}

// checkReportForKnownOrganizationWrongCluster checks how improper cluster name
// is checked by REST API handler
func checkReportForKnownOrganizationWrongCluster() {
	clusterName := "abcdefghijklmnopqrstuvwyz"
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with improper cluster name").Get(reportEndpointForOrgAndCluster(organization1, clusterName))
	f.Send()
	f.ExpectStatus(http.StatusBadRequest)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	f.PrintReport()
}

// checkWrongMethodsForReportForOrgAndClusterEndpoint checks whether other HTTP methods are
// rejected correctly for the REST API 'report' point
func checkWrongMethodsForReportForOrgAndClusterEndpoint() {
	// known organizations
	checkGetEndpointByOtherMethods(reportEndpointForOrgAndCluster(organization1, cluster1ForOrg1), false)
	checkGetEndpointByOtherMethods(reportEndpointForOrgAndCluster(organization2, cluster1ForOrg1), false)

	// unknown organizations
	checkGetEndpointByOtherMethods(reportEndpointForOrgAndCluster(1, ""), false)
	checkGetEndpointByOtherMethods(reportEndpointForOrgAndCluster(2, ""), false)
}

// checkReportForKnownCluster checks if proper report is returned for
// known cluster name (w/o organization ID)
func checkReportForKnownCluster() {
	url := reportEndpointForCluster(cluster1ForOrg1)
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with known cluster w/o org").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := FullReportResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError(statusShouldBeSetToOK)
		}
	}
	f.PrintReport()
}

// checkReportForUnknownCluster checks how unknown cluster
// name is checked by REST API handler
func checkReportForUnknownCluster() {
	url := reportEndpointForCluster(unknownCluster)
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with unknown cluster w/o org").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusNotFound)
	f.PrintReport()
}

// checkReportForImproperCluster checks how improper cluster
// name is checked by REST API handler
func checkReportForImproperCluster() {
	url := reportEndpointForCluster("foobarbaz")
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with improper cluster w/o org").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusBadRequest)
	f.PrintReport()
}

// checkReportForFailedCluster checks how cluster with special failures setup
// is checked by REST API handler (expected HTTP code is 200)
func checkReportForFailedCluster200() {
	// expected HTTP code is encoded in last three characters in cluster name
	// in this case it should be 200 OK
	url := reportEndpointForCluster("ffffffff-ffff-ffff-ffff-000000000200")
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with expected HTTP code 200").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.PrintReport()
}

// checkReportForFailedCluster checks how cluster with special failures setup
// is checked by REST API handler (expected HTTP code is 400)
func checkReportForFailedCluster400() {
	// expected HTTP code is encoded in last three characters in cluster name
	// in this case it should be 400 Bad Request
	url := reportEndpointForCluster("ffffffff-ffff-ffff-ffff-000000000400")
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with expected HTTP code 400").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusBadRequest)
	f.PrintReport()
}

// checkReportForFailedCluster checks how cluster with special failures setup
// is checked by REST API handler
func checkReportForFailedClusterNegativeTestCase() {
	// expected HTTP code is encoded in last three characters in cluster name
	// parsing of "fff" should raise an error on server side that is caught there
	url := reportEndpointForCluster("ffffffff-ffff-ffff-ffff-000000000fff")
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with expected HTTP code 200").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.PrintReport()
}

// checkReportForFailedCluster checks how cluster with special failures setup
// is checked by REST API handler (expected HTTP code is 500)
func checkReportForFailedCluster500() {
	// expected HTTP code is encoded in last three characters in cluster name
	// in this case it should be 500 Internal Server Error
	url := reportEndpointForCluster("ffffffff-ffff-ffff-ffff-000000000500")
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method with expected HTTP code 500").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusInternalServerError)
	f.PrintReport()
}

// checkWrongMethodsForClusterReportEndpoint checks whether other HTTP methods are
// rejected correctly for the REST API 'report' point
func checkWrongMethodsForClusterReportEndpoint() {
	checkGetEndpointByOtherMethods(reportEndpointForCluster(cluster1ForOrg1), false)
	checkGetEndpointByOtherMethods(reportEndpointForCluster(cluster1ForOrg1), false)
}

// checkReportsForAllClustersInOrganizationPositiveTestCase check the REST API endpoint
// to retrieve reports for all clusters in given organization
func checkReportsForAllClustersInOrganizationPositiveTestCase() {
	url := reportEndpointForAllReportsForOrg(1) // proper org ID
	f := frisby.Create("Check the 'reports for all clusters' REST API point using HTTP GET method with proper organization").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		// try to parse the payload
		response := AllReportsForOrganizationResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		// parsing was ok, so check response content
		// ATM no clusters and no errors should be returned from server
		if len(response.Clusters) != 0 {
			f.AddError("Expecting empty list of clusters")
		}
		if len(response.Errors) != 0 {
			f.AddError("Expecting empty list of errors")
		}
	}
	f.PrintReport()
}
