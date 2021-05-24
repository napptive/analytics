/**
 * Copyright 2020 Napptive
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package analytics

import (
	"github.com/napptive/analytics/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var credentialPath string // = "/Users/cdelope/tmp/carmentest-d56ccf2df548.json"
var projectID string // = "carmentest"

var _ = ginkgo.Describe("Provider test", func() {

	if !utils.RunIntegrationTests("provider") {
		log.Warn().Msg("provider integration tests are skipped")
		return
	}

	credentialPath = os.Getenv("CREDENTIALS_PATH")
	if credentialPath == "" {
		log.Fatal().Msg("CREDENTIALS_PATH not found")
	}

	projectID = os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal().Msg("PROJECT_ID not found")
	}


	ginkgo.It("should be able to add a invitation", func() {
		loopTime := time.Second * 2
		provider, err := NewBigQueryProvider(projectID, credentialPath, loopTime)
		gomega.Expect(provider).ShouldNot(gomega.BeNil())
		gomega.Expect(err).To(gomega.Succeed())

		for i:= 0; i<= 10; i ++ {
			err = provider.SendLoginData(utils.GenerateLoginData())
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.SendOperationData(utils.GenerateOperationData())
			gomega.Expect(err).To(gomega.Succeed())

			time.Sleep(loopTime / 3)
		}

	})
})
