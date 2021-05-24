package analytics

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"time"
)

var _ = ginkgo.Describe("Config test", func() {

	ginkgo.It("config must be full filled", func() {
		cfg := BigQueryConfig{
			projectID:       "project",
			credentialsPath: "/Users/credentials.json",
			loopTime:        time.Second * 2,
		}
		err := cfg.IsValid()
		gomega.Expect(err).Should(gomega.Succeed())
	})
	ginkgo.It("projectID mus be filled", func() {
		cfg := BigQueryConfig{
			projectID:       "",
			credentialsPath: "/Users/credentials.json",
			loopTime:        time.Second * 2,
		}
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
	ginkgo.It("credentialsPath must be filled", func() {
		cfg := BigQueryConfig{
			projectID:       "project",
			credentialsPath: "",
			loopTime:        time.Second * 2,
		}
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
	ginkgo.It("looptime must be filled", func() {
		cfg := BigQueryConfig{
			projectID:       "project",
			credentialsPath: "/Users/credentials.json",
		}
		err := cfg.IsValid()
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})
})

