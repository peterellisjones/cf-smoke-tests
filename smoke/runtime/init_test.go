package runtime

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
	"github.com/cloudfoundry/cf-smoke-tests/smoke"
)

const (
	// timeout for cf push cli calls
	CF_PUSH_TIMEOUT_IN_SECONDS = 300

	// timeout for cf scale cli calls
	CF_SCALE_TIMEOUT_IN_SECONDS = 120

	// timeout for cf app cli calls
	CF_APP_STATUS_TIMEOUT_IN_SECONDS = 120
)

func TestSmokeTests(t *testing.T) {
	RegisterFailHandler(Fail)

	testConfig := smoke.GetConfig()
	testSetup := workflowhelpers.NewSmokeTestSuiteSetup(testConfig)

	SynchronizedBeforeSuite(func() []byte {
		testSetup.Setup()
		return nil
	}, func(data []byte) {})

	SynchronizedAfterSuite(func() {
	}, func() {
		testSetup.Teardown()
	})

	rs := []Reporter{}

	if testConfig.ArtifactsDirectory != "" {
		os.Setenv("CF_TRACE", traceLogFilePath(testConfig))
		rs = append(rs, reporters.NewJUnitReporter(jUnitReportFilePath(testConfig)))
	}

	if testConfig.Reporter == "TeamCity" {
		rs = append(rs, reporters.NewTeamCityReporter(GinkgoWriter))
	}

	RunSpecsWithDefaultAndCustomReporters(t, "CF-Runtime-Smoke-Tests", rs)
}

func traceLogFilePath(testConfig *smoke.Config) string {
	return filepath.Join(testConfig.ArtifactsDirectory, fmt.Sprintf("CF-TRACE-%s-%d.txt", testConfig.SuiteName, ginkgoNode()))
}

func jUnitReportFilePath(testConfig *smoke.Config) string {
	return filepath.Join(testConfig.ArtifactsDirectory, fmt.Sprintf("junit-%s-%d.xml", testConfig.SuiteName, ginkgoNode()))
}

func ginkgoNode() int {
	return ginkgoconfig.GinkgoConfig.ParallelNode
}