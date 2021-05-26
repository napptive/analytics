package config

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"time"
)

func GenerateTestBigQueryConfig() BigQueryConfig {
	return BigQueryConfig{
		ProjectID:       "projectID",
		Schema:          "schema",
		Table:           "operation",
		CredentialsPath: "/Users/tmp/credentials.json",
		SendingTime:     time.Second * 2,
	}
}

var _ = ginkgo.Describe("Config test", func() {

	ginkgo.It("config must be full filled", func() {
		cfg := GenerateTestBigQueryConfig()
		err := cfg.IsValid()
		gomega.Expect(err).Should(gomega.Succeed())
	})
	ginkgo.It("projectID mus be filled", func() {
		cfg := GenerateTestBigQueryConfig()
		cfg.ProjectID = ""
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
	ginkgo.It("schema must be filled", func() {
		cfg := GenerateTestBigQueryConfig()
		cfg.Schema = ""
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
	ginkgo.It("Table must be filled", func() {
		cfg := GenerateTestBigQueryConfig()
		cfg.Table = ""
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
	ginkgo.It("credentialsPath must be filled", func() {
		cfg := GenerateTestBigQueryConfig()
		cfg.CredentialsPath = ""
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
	ginkgo.It("SendingTime must be filled", func() {
		cfg := GenerateTestBigQueryConfig()
		cfg.SendingTime = 0
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
})
