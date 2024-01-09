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
*/

// Package storage contains an implementation of interface between Go code and
// (almost any) SQL database like PostgreSQL, SQLite, or MariaDB. An implementation
// named DBStorage is constructed via New function and it is mandatory to call Close
// for any opened connection to database. The storage might be initialized by Init
// method if database schema is empty.
package storage

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

const clusterNotFoundMessage = "Cluster not found"

// Storage represents an interface to almost any database or storage system
type Storage interface {
	Init() error
	Close() error
	ListOfOrgs() ([]types.OrgID, error)
	ListOfClustersForOrg(orgID types.OrgID) ([]types.ClusterName, error)
	ReadReportForCluster(clusterName types.ClusterName) (types.ClusterReport, error)
	ReadReportForOrganizationAndCluster(orgID types.OrgID, clusterName types.ClusterName) (types.ClusterReport, error)
	GetRuleWithContent(ruleID types.RuleID, ruleErrorKey types.ErrorKey) (*types.RuleWithContent, error)
	GetPredictionForCluster(cluster types.ClusterName) (*types.UpgradeRiskPrediction, error)
}

// MemoryStorage data structure represents configuration of memory storage used
// to store mock data.
type MemoryStorage struct {
}

// Special clusters can change results in given time period, for example each
// 10 minutes or so. This is to simulate real world behaviour.
const changingClustersPeriodInMinutes = 15

const noPermissionsForOrg = "You have no permissions to get or change info about this organization"

var reports = make(map[string]string)

func readReport(path, clusterName string) (string, error) {
	absPath, err := filepath.Abs(path + "/report_" + clusterName + ".json")
	if err != nil {
		return "", err
	}

	// disable "G304 (CWE-22): Potential file inclusion via variable"
	report, err := os.ReadFile(absPath) // #nosec G304
	if err != nil {
		return "", err
	}
	return string(report), nil
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func initStorage(path string) error {
	clusters := []string{
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a266",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a267",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a268",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a269",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26a",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26b",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26c",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26d",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26e",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26f",
		"74ae54aa-6577-4e80-85e7-697cb646ff37",
		"a7467445-8d6a-43cc-b82c-7007664bdf69",
		"ee7d2bf4-8933-4a3a-8634-3328fe806e08",
		"eeeeeeee-eeee-eeee-eeee-000000000001",
		"00000001-624a-49a5-bab8-4fdc5e51a266",
		"00000001-624a-49a5-bab8-4fdc5e51a267",
		"00000001-624a-49a5-bab8-4fdc5e51a268",
		"00000001-624a-49a5-bab8-4fdc5e51a269",
		"00000001-624a-49a5-bab8-4fdc5e51a26a",
		"00000001-624a-49a5-bab8-4fdc5e51a26b",
		"00000001-624a-49a5-bab8-4fdc5e51a26c",
		"00000001-624a-49a5-bab8-4fdc5e51a26d",
		"00000001-624a-49a5-bab8-4fdc5e51a26e",
		"00000001-624a-49a5-bab8-4fdc5e51a26f",
		"00000001-6577-4e80-85e7-697cb646ff37",
		"00000001-8933-4a3a-8634-3328fe806e08",
		"00000001-8d6a-43cc-b82c-7007664bdf69",
		"00000001-eeee-eeee-eeee-000000000001",
		"00000002-624a-49a5-bab8-4fdc5e51a266",
		"00000002-6577-4e80-85e7-697cb646ff37",
		"00000002-8933-4a3a-8634-3328fe806e08",
		"00000003-8933-4a3a-8634-3328fe806e08",
		"00000003-8d6a-43cc-b82c-7007664bdf69",
		"00000003-eeee-eeee-eeee-000000000001",
	}
	for _, cluster := range clusters {
		report, err := readReport(path, cluster)
		if err != nil {
			return err
		}
		reports[cluster] = report
	}
	return nil
}

// New function creates and initializes a new instance of Storage interface
func New(path string) (*MemoryStorage, error) {
	err := initStorage(path)
	return &MemoryStorage{}, err
}

// Init performs all database initialization
// tasks necessary for further service operation.
func (storage MemoryStorage) Init() error {
	log.Info().Msg("Initializing connection to data storage")
	return nil
}

// Close method closes the connection to database. Needs to be called at the end of application lifecycle.
func (storage MemoryStorage) Close() error {
	log.Info().Msg("Closing connection to data storage")
	return nil
}

// Report represents one (latest) cluster report.
//
//	Org: organization ID
//	Name: cluster GUID in the following format:
//	    c8590f31-e97e-4b85-b506-c45ce1911a12
type Report struct {
	Org        types.OrgID         `json:"org"`
	Name       types.ClusterName   `json:"cluster"`
	Report     types.ClusterReport `json:"report"`
	ReportedAt types.Timestamp     `json:"reported_at"`
}

// ListOfOrgs reads list of all organizations that have at least one cluster report
func (storage MemoryStorage) ListOfOrgs() ([]types.OrgID, error) {
	orgs := []types.OrgID{
		11789772,
		11940171,
	}
	return orgs, nil
}

func clustersForOrganization11789772() []types.ClusterName {
	clusters := make([]types.ClusterName, 0)
	clusters = append(clusters,
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a266",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a267",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a268",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a269",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26a",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26b",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26c",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26d",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26e",
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a26f",
		"74ae54aa-6577-4e80-85e7-697cb646ff37",
		"a7467445-8d6a-43cc-b82c-7007664bdf69",
		"ee7d2bf4-8933-4a3a-8634-3328fe806e08",
		"eeeeeeee-eeee-eeee-eeee-000000000001")
	return clusters
}

// ListOfClustersForOrg reads list of all clusters fro given organization
func (storage MemoryStorage) ListOfClustersForOrg(orgID types.OrgID) ([]types.ClusterName, error) {
	clusters := make([]types.ClusterName, 0)
	switch orgID {
	case 11940171:
		return clusters, errors.New(noPermissionsForOrg)
	case 11789772:
		return clustersForOrganization11789772(), nil
	case 1:
		clusters = append(clusters,
			"00000001-624a-49a5-bab8-4fdc5e51a266",
			"00000001-624a-49a5-bab8-4fdc5e51a267",
			"00000001-624a-49a5-bab8-4fdc5e51a268",
			"00000001-624a-49a5-bab8-4fdc5e51a269",
			"00000001-624a-49a5-bab8-4fdc5e51a26a",
			"00000001-624a-49a5-bab8-4fdc5e51a26b",
			"00000001-624a-49a5-bab8-4fdc5e51a26c",
			"00000001-624a-49a5-bab8-4fdc5e51a26d",
			"00000001-624a-49a5-bab8-4fdc5e51a26e",
			"00000001-624a-49a5-bab8-4fdc5e51a26f",
			"00000001-6577-4e80-85e7-697cb646ff37",
			"00000001-8933-4a3a-8634-3328fe806e08",
			"00000001-8d6a-43cc-b82c-7007664bdf69",
			"00000001-eeee-eeee-eeee-000000000001")
	case 2:
		clusters = append(clusters,
			"00000002-624a-49a5-bab8-4fdc5e51a266",
			"00000002-6577-4e80-85e7-697cb646ff37",
			"00000002-8933-4a3a-8634-3328fe806e08")
	case 3:
		clusters = append(clusters,
			"00000003-8933-4a3a-8634-3328fe806e08",
			"00000003-8d6a-43cc-b82c-7007664bdf69",
			"00000003-eeee-eeee-eeee-000000000001")
	}

	return clusters, nil
}

func getReportForCluster(clusterName types.ClusterName) (string, bool) {
	report, ok := reports[string(clusterName)]
	if !ok {
		return "", false
	}
	return report, true
}

// ReadReportForCluster reads result (health status) for selected cluster
func (storage MemoryStorage) ReadReportForCluster(
	clusterName types.ClusterName,
) (types.ClusterReport, error) {
	var report string

	// clusters that can change its output (report)
	// please note that these clusters have special name:
	// "cccccccc-cccc-cccc-cccc-{index}"
	//
	// Mnemotechnic: c - changing
	changingClusters := map[string][]string{
		"cccccccc-cccc-cccc-cccc-000000000001": {
			"34c3ecc5-624a-49a5-bab8-4fdc5e51a266",
			"74ae54aa-6577-4e80-85e7-697cb646ff37",
			"a7467445-8d6a-43cc-b82c-7007664bdf69"},
		"cccccccc-cccc-cccc-cccc-000000000002": {
			"74ae54aa-6577-4e80-85e7-697cb646ff37",
			"a7467445-8d6a-43cc-b82c-7007664bdf69",
			"ee7d2bf4-8933-4a3a-8634-3328fe806e08"},
		"cccccccc-cccc-cccc-cccc-000000000003": {
			"ee7d2bf4-8933-4a3a-8634-3328fe806e08",
			"ee7d2bf4-8933-4a3a-8634-3328fe806e08",
			"34c3ecc5-624a-49a5-bab8-4fdc5e51a266"},
		"cccccccc-cccc-cccc-cccc-000000000004": {
			"eeeeeeee-eeee-eeee-eeee-000000000001",
			"eeeeeeee-eeee-eeee-eeee-000000000001",
			"34c3ecc5-624a-49a5-bab8-4fdc5e51a266"},
	}

	reportName := clusterName

	// handling for clusters that can change its report
	if changingCluster, found := changingClusters[string(clusterName)]; found {
		reportName = chooseReport(changingCluster)
	}

	report, found := getReportForCluster(reportName)
	if !found {
		return types.ClusterReport(""), errors.New(clusterNotFoundMessage)
	}

	return types.ClusterReport(report), nil
}

// chooseReport for "changing cluster"
func chooseReport(variants []string) types.ClusterName {
	const operationName = "changingCluster"

	// first we need to get the minute in hour
	currentTime := time.Now()
	minute := currentTime.Minute()
	log.Info().Int("Minute in hour", minute).Msg(operationName)

	// then compute index of report
	i := minute / changingClustersPeriodInMinutes
	i %= len(variants)

	// and choose the report according to the index
	cluster := variants[i]
	log.Info().Int("Index", i).Msg(operationName)
	log.Info().Str("Cluster", cluster).Msg(operationName)
	return types.ClusterName(cluster)
}

// ReadReportForOrganizationAndCluster reads result (health status) for
// selected cluster for given organization
func (storage MemoryStorage) ReadReportForOrganizationAndCluster(
	orgID types.OrgID, clusterName types.ClusterName,
) (types.ClusterReport, error) {
	var report string

	switch orgID {
	case 11940171:
		return types.ClusterReport(report), errors.New(noPermissionsForOrg)
	case 1, 2, 3, 11789772:
		report, found := getReportForCluster(clusterName)
		if found {
			return types.ClusterReport(report), nil
		}
	}

	return types.ClusterReport(report), errors.New(clusterNotFoundMessage)
}

// GetPredictionForCluster gets a prediction for the cluster
func (storage MemoryStorage) GetPredictionForCluster(_ types.ClusterName) (*types.UpgradeRiskPrediction, error) {
	return &types.UpgradeRiskPrediction{
		Recommended: true,
		Predictors: types.UpgradeRisksPredictors{
			Alerts:             []types.Alert{},
			OperatorConditions: []types.OperatorCondition{},
		},
	}, nil
}
