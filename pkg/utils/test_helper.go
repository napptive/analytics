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
package utils

import (
	"context"
	"github.com/napptive/analytics/pkg/config"
	"github.com/napptive/analytics/pkg/entities"
	"github.com/onsi/gomega"
	"github.com/rs/xid"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

// RunIntegrationTests checks whether integration tests should be executed.
func RunIntegrationTests(id string) bool {
	var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")
	if runIntegration == "all" {
		return true
	}
	return runIntegration == id
}

// GenerateLoginData returns a random LoginData
func GenerateTestLoginData() entities.LoginData {
	return entities.LoginData{
		Timestamp: time.Now(),
		UserID:    xid.New().String(),
		Method:    "CLI",
	}
}

// GenerateOperationData returns a random OperationData
func GenerateTestOperationData() entities.OperationData {
	return entities.OperationData{
		Timestamp: time.Now(),
		UserID:    xid.New().String(),
		Operation: "PEPE",
	}
}

func GenerateTestOperation() entities.Operation {
	return entities.Operation{
		Timestamp: time.Now(),
		UserID:    xid.New().String(),
		Operation: "test/operation",
	}
}

func GenerateTestFullContext() context.Context {
	md := metadata.New(map[string]string{"user_id": xid.New().String()})
	return metadata.NewOutgoingContext(context.Background(), md)
}

func GetBigQueryConfig() *config.BigQueryConfig {
	var credentialPath = os.Getenv("CREDENTIALS_PATH")
	gomega.Expect(credentialPath).ShouldNot(gomega.BeEmpty())

	var projectID = os.Getenv("PROJECT_ID")
	gomega.Expect(projectID).ShouldNot(gomega.BeEmpty())

	var schema = os.Getenv("SCHEMA")
	gomega.Expect(schema).ShouldNot(gomega.BeEmpty())

	var table = os.Getenv("TABLE")
	gomega.Expect(table).ShouldNot(gomega.BeEmpty())


	bqConfig := config.NewBigQueryConfig(projectID, schema, table, credentialPath, time.Second)

	return &bqConfig
}

