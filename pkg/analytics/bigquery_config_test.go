package analytics

import (
	"github.com/napptive/analytics/pkg/config"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"time"
)

var _ = ginkgo.Describe("Config test", func() {

	ginkgo.It("config must be full filled", func() {
		cfg := config.NewBigQueryConfig("project", "/Users/credentials.json", time.Second*2)
		err := cfg.IsValid()
		gomega.Expect(err).Should(gomega.Succeed())
	})
	ginkgo.It("projectID mus be filled", func() {
		cfg := config.NewBigQueryConfig("", "/Users/credentials.json", time.Second*2)
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
	ginkgo.It("credentialsPath must be filled", func() {
		cfg := config.NewBigQueryConfig("project", "", time.Second*2)
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
	ginkgo.It("looptime must be filled", func() {
		cfg := config.NewBigQueryConfig("project", "/Users/credentials.json", 0)
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
})
