// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

package obfuscate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type execPlanTestCase struct {
	// Rdbms is mysql or postgres
	Rdbms          string                 `json:"rdbms"`
	TestPlan       map[string]interface{} `json:"test_plan"`
	ObfuscatedPlan map[string]interface{} `json:"obfuscated_plan"`
	NormalizedPlan map[string]interface{} `json:"normalized_plan"`
}

const mysqlTestCasesFile = "./testdata/msyql_execution_plan_test_cases.json"
const postgresTestCasesFile = "./testdata/postgres_execution_plan_test_cases.json"

// loadTests loads all XML tests from ./testdata/obfuscate.xml
func loadExecutionPlanTestCases(path string) ([]execPlanTestCase, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var testPlans []execPlanTestCase
	err = json.Unmarshal(raw, &testPlans)
	if err != nil {
		return nil, err
	}
	for _, plan := range testPlans {
		if plan.TestPlan == nil {
			return nil, fmt.Errorf("nil testPlan found %v", plan)
		}
		if plan.ObfuscatedPlan == nil {
			return nil, fmt.Errorf("nil obfuscatedPlan found")
		}
		if plan.NormalizedPlan == nil {
			return nil, fmt.Errorf("nil noramlizedPlan found")
		}
	}
	return testPlans, err
}

func loadAllTestCases() (result []execPlanTestCase, err error) {
	mysqlTestCases, err := loadExecutionPlanTestCases(mysqlTestCasesFile)
	if err != nil {
		return nil, err
	}
	result = append(result, mysqlTestCases...)
	postgresTestCases, err := loadExecutionPlanTestCases(postgresTestCasesFile)
	if err != nil {
		return nil, err
	}
	result = append(result, postgresTestCases...)
	return result, nil
}

func TestPlanObfuscation(t *testing.T) {
	testCases, err := loadAllTestCases()
	if err != nil {
		assert.FailNowf(t, "failed to load test cases", err.Error())
	}

	assert.NotEmpty(t, testCases)
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test_%d_%s", i, testCase.Rdbms), func(t *testing.T) {
			result := NewObfuscator(nil).ObfuscateSQLPlan(testCase.TestPlan)
			assert.Equal(t, testCase.ObfuscatedPlan, result)
		})
	}
}
