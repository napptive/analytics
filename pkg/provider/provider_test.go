package provider

import (
	"github.com/napptive/analytics/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"time"
)

/*
1. Create a service account and get credentials
2. Create a bigquery schema in a project
3. Create login table (timestamp TIMESTAMP, userID STRING, operation STRING)
4. Create operation table (timestamp TIMESTAMP, userID STRING, method STRING)
5. Allow access to serviceAccount

Environment variables:
RUN_IT_TEST: provider
CREDENTIALS_PATH: <path_to_credentials_file>
PROJECT_ID: project
SCHEMA: <schema_name>
TABLE: <operation>
*/


var _ = ginkgo.Describe("Provider test", func() {
	var provider Provider
	var proError error

	if !utils.RunIntegrationTests("provider") {
		log.Warn().Msg("provider integration tests are skipped")
		return
	}

	// Create provider
	cfg := utils.GetBigQueryConfig()

	provider, proError = NewBigQueryProvider(*cfg)
	gomega.Expect(proError).Should(gomega.Succeed())

	ginkgo.It("should be able to add an operation", func() {

		op := utils.GenerateTestOperation()
		for i := 0; i <= 10; i++ {
			err := provider.Send(op)
			gomega.Expect(err).To(gomega.Succeed())
			time.Sleep(cfg.SendingTime/3)
		}
	})
	ginkgo.It("should be able to add an operation", func() {

		op := utils.GenerateTestOperation()
		for i := 0; i <= 10; i++ {
			err := provider.Send(op)
			gomega.Expect(err).To(gomega.Succeed())
		}
		err := provider.Flush()
		gomega.Expect(err).To(gomega.Succeed())
	})
})