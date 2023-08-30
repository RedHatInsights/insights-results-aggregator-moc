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

package tests

import (
	"encoding/json"
	"fmt"

	"github.com/verdverm/frisby"
)

const (
	knownClusterForOrganization1   = "34c3ecc5-624a-49a5-bab8-4fdc5e51a266"
	knownCluster2ForOrganization1  = "74ae54aa-6577-4e80-85e7-697cb646ff37"
	knownCluster3ForOrganization1  = "a7467445-8d6a-43cc-b82c-7007664bdf69"
	unknownClusterForOrganization1 = "bbbbbbbb-bbbb-bbbb-bbbb-cccccccccccc"
)

// MultipleReportsResponse represents response from the server that contains
// results for multiple clusters together with overall status
type MultipleReportsResponse struct {
	Clusters    []string               `json:"clusters"`
	Errors      []string               `json:"errors"`
	Reports     map[string]interface{} `json:"reports"`
	GeneratedAt string                 `json:"generated_at"`
	Status      string                 `json:"status"`
}

// ClusterListInRequest represents request body containing list of clusters
type ClusterListInRequest struct {
	Clusters []string `json:"clusters"`
}

// constructURLForReportForOrgClustersPostMethod function construct an URL to
// access the endpoint to retrieve results for given list of clusters using
// POST method
func constructURLForReportForOrgClustersPostMethod() string {
	return fmt.Sprintf("%sclusters", apiURL)
}

// sendClusterListInPayload function sends the cluster list in request payload
// to server
func sendClusterListInPayload(f *frisby.Frisby, clusterList []string) {
	var payload = ClusterListInRequest{
		Clusters: clusterList,
	}
	// create payload
	f.SetJson(payload)

	// and perform send
	f.Send()
}

// readMultipleReportsResponse reads and parses response body that should
// contains reports for multiple clusters
func readMultipleReportsResponse(f *frisby.Frisby) MultipleReportsResponse {
	response := MultipleReportsResponse{}

	// try to read response body
	text, err := f.Resp.Content()

	if err != nil {
		f.AddError(err.Error())
	} else {
		// try to deserialize response body
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
	}
	return response
}

// expectNumberOfClusters utility function checks if server response contains
// expected number of clusters
func expectNumberOfClusters(f *frisby.Frisby, response MultipleReportsResponse, expected int) {
	clusters := response.Clusters
	actual := len(clusters)
	if actual != expected {
		f.AddError(fmt.Sprintf("expected %d clusters in server response, but got %d instead", expected, actual))
	}
}

// expectNumberOfClusters utility function checks if server response contains
// expected number of errors
func expectNumberOfErrors(f *frisby.Frisby, response MultipleReportsResponse, expected int) {
	clusters := response.Errors
	actual := len(clusters)
	if actual != expected {
		f.AddError(fmt.Sprintf("expected %d errors in server response, but got %d instead", expected, actual))
	}
}

// expectNumberOfReports utility function checks if server response contains
// expected number of errors
func expectNumberOfReports(f *frisby.Frisby, response MultipleReportsResponse, expected int) {
	clusters := response.Reports
	actual := len(clusters)
	if actual != expected {
		f.AddError(fmt.Sprintf("expected %d reports in server response, but got %d instead", expected, actual))
	}
}

// expectClusterInResponse utility function checks if server response contains
// expected cluster name
func expectClusterInResponse(f *frisby.Frisby, response MultipleReportsResponse, clusterName string) {
	clusters := response.Clusters
	for _, cluster := range clusters {
		// cluster has been found
		if cluster == clusterName {
			return
		}
	}
	// cluster was not found
	f.AddError(fmt.Sprintf("Cluster %s can not be found in server response in the cluster list", clusterName))
}

// expectClusterInResponse utility function checks if server response contains
// expected report for specified cluster
func expectReportInResponse(f *frisby.Frisby, response MultipleReportsResponse, clusterName string) {
	reports := response.Reports
	for cluster := range reports {
		// cluster has been found
		if cluster == clusterName {
			return
		}
	}
	// report for cluster was not found
	f.AddError(fmt.Sprintf("Cluster %s can not be found in server response in reports map", clusterName))
}

// expectErrorClusterInResponse utility function checks if server response
// contains expected error
func expectErrorClusterInResponse(f *frisby.Frisby, response MultipleReportsResponse, clusterName string) {
	errors := response.Errors
	for _, cluster := range errors {
		// cluster has been found
		if cluster == clusterName {
			return
		}
	}
	// error for cluster was not found
	f.AddError(fmt.Sprintf("Cluster %s can not be found in server response in the errors list", clusterName))
}

// checkMultipleReportsForKnownOrganizationAnd1KnownClusterUsingPostMethod check the endpoint that returns multiple results
func checkMultipleReportsForKnownOrganizationAnd1KnownClusterUsingPostMethod() {
	clusterList := []string{
		knownClusterForOrganization1,
	}

	// send request to the endpoint
	url := constructURLForReportForOrgClustersPostMethod()
	f := frisby.Create("Check the endpoint to return report for existing organization and one cluster ID (POST variant)").Post(url)
	sendClusterListInPayload(f, clusterList)

	// check the response from server
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the payload returned from server
	response := readMultipleReportsResponse(f)
	expectNumberOfClusters(f, response, 1)
	expectNumberOfErrors(f, response, 0)
	expectNumberOfReports(f, response, 1)
	expectClusterInResponse(f, response, knownClusterForOrganization1)
	expectReportInResponse(f, response, knownClusterForOrganization1)

	f.PrintReport()
}

// checkMultipleReportsForKnownOrganizationAnd2KnownClustersUsingPostMethod check the endpoint that returns multiple results
func checkMultipleReportsForKnownOrganizationAnd2KnownClustersUsingPostMethod() {
	clusterList := []string{
		knownClusterForOrganization1,
		knownCluster2ForOrganization1,
	}

	// send request to the endpoint
	url := constructURLForReportForOrgClustersPostMethod()
	f := frisby.Create("Check the endpoint to return report for existing organization and two cluster IDs (POST variant)").Post(url)
	sendClusterListInPayload(f, clusterList)

	// check the response from server
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the payload returned from server
	response := readMultipleReportsResponse(f)
	expectNumberOfClusters(f, response, 2)
	expectNumberOfErrors(f, response, 0)
	expectNumberOfReports(f, response, 2)
	expectClusterInResponse(f, response, knownClusterForOrganization1)
	expectClusterInResponse(f, response, knownCluster2ForOrganization1)
	expectReportInResponse(f, response, knownClusterForOrganization1)
	expectReportInResponse(f, response, knownCluster2ForOrganization1)

	f.PrintReport()
}

// checkMultipleReportsForKnownOrganizationAnd3KnownClustersUsingPostMethod check the endpoint that returns multiple results
func checkMultipleReportsForKnownOrganizationAnd3KnownClustersUsingPostMethod() {
	clusterList := []string{
		knownClusterForOrganization1,
		knownCluster2ForOrganization1,
		knownCluster3ForOrganization1,
	}

	// send request to the endpoint
	url := constructURLForReportForOrgClustersPostMethod()
	f := frisby.Create("Check the endpoint to return report for existing organization and three cluster IDs (POST variant)").Post(url)
	sendClusterListInPayload(f, clusterList)

	// check the response from server
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the payload returned from server
	response := readMultipleReportsResponse(f)
	expectNumberOfClusters(f, response, 3)
	expectNumberOfErrors(f, response, 0)
	expectNumberOfReports(f, response, 3)
	expectClusterInResponse(f, response, knownClusterForOrganization1)
	expectClusterInResponse(f, response, knownCluster2ForOrganization1)
	expectClusterInResponse(f, response, knownCluster3ForOrganization1)
	expectReportInResponse(f, response, knownClusterForOrganization1)
	expectReportInResponse(f, response, knownCluster2ForOrganization1)
	expectReportInResponse(f, response, knownCluster3ForOrganization1)

	f.PrintReport()
}

// checkMultipleReportsForKnownOrganizationAndUnknownClusterUsingPostMethod check the endpoint that returns multiple results
func checkMultipleReportsForKnownOrganizationAndUnknownClusterUsingPostMethod() {
	clusterList := []string{
		unknownClusterForOrganization1,
	}

	// send request to the endpoint
	url := constructURLForReportForOrgClustersPostMethod()
	f := frisby.Create("Check the endpoint to return report for existing organization and one unknown cluster ID (POST variant)").Post(url)
	sendClusterListInPayload(f, clusterList)

	// check the response from server
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the payload returned from server
	response := readMultipleReportsResponse(f)
	expectNumberOfClusters(f, response, 0)
	expectNumberOfErrors(f, response, 1)
	expectNumberOfReports(f, response, 0)
	expectErrorClusterInResponse(f, response, unknownClusterForOrganization1)

	f.PrintReport()
}
