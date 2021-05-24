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
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/napptive/analytics/pkg/entities"
	"github.com/napptive/nerrors/pkg/nerrors"
	"google.golang.org/api/option"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	schema         = "analytics"
	loginTable     = "login"
	operationTable = "operation"
	databaseTimeout = time.Second * 5
)

type BigQueryConfig struct {
	projectID string
	credentialsPath string
	loopTime time.Duration
}

func (bqc *BigQueryConfig) IsValid() error  {
	if bqc.projectID == "" {
		return nerrors.NewFailedPreconditionError("projectID mus be filled")
	}
	if bqc.credentialsPath == "" {
		return nerrors.NewFailedPreconditionError("credentials path must be informed")
	}
	if bqc.loopTime <= 0 {
		return nerrors.NewFailedPreconditionError("loopTime mus be filled")
	}
	return nil
}

type BigQueryProvider struct {
	Client *bigquery.Client

	// sendTime is the time between database inserts
	sendTime time.Duration

	// loginCache with a cache of login operation
	// in this cache we store the login operations to avoid send msg by msg
	loginCache []entities.LoginData
	// Mutex for managing login cache access.
	loginMutex sync.Mutex

	// opCache with a cache of operation operation
	// in this cache we store the operation messages to avoid send msg by msg
	opCache []entities.OperationData
	// Mutex for managing operation cache access.
	opMutex sync.Mutex
}

// NewBigQueryProvider
func NewBigQueryProvider(cfg BigQueryConfig) (Provider, error) {

	// validate config
	if err := cfg.IsValid(); err != nil {
		return nil, err
	}

	client, err := bigquery.NewClient(context.Background(), cfg.projectID,
		option.WithCredentialsFile(cfg.credentialsPath))
	if err != nil {
		return nil, err
	}
	provider := &BigQueryProvider{
		Client:     client,
		loginCache: []entities.LoginData{},
		sendTime:   cfg.loopTime,
	}

	// start the loop
	go provider.SendDataLoop()

	return provider, nil
}

// SendLoginData stores the data in the login cache
func (bq *BigQueryProvider) SendLoginData(data entities.LoginData) error {

	bq.loginMutex.Lock()
	defer bq.loginMutex.Unlock()

	// store the login in the cache
	bq.loginCache = append(bq.loginCache, data)

	return nil
}

// SendOperationData stores the data in the operation cache
func (bq *BigQueryProvider) SendOperationData(data entities.OperationData) error {

	bq.opMutex.Lock()
	defer bq.opMutex.Unlock()

	// store the operation in the cache
	bq.opCache = append(bq.opCache, data)

	return nil
}

// SendLoginCache inserts all the data stored in login cache in the database
func (bq *BigQueryProvider) SendLoginCache() error {
	bq.loginMutex.Lock()
	toSend := bq.loginCache
	bq.loginCache = []entities.LoginData{}
	bq.loginMutex.Unlock()

	if len(toSend) == 0 {
		return nil
	}

	i := bq.Client.Dataset(schema).Table(loginTable).Inserter()

	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	err := i.Put(ctx, toSend)
	if err != nil {
		return nerrors.NewInternalErrorFrom(err, "error sending login data")
	}
	// empty cache
	bq.loginCache = []entities.LoginData{}
	return nil
}

// SendOperationCache inserts all the data stored in operation cache in the database
func (bq *BigQueryProvider) SendOperationCache() error {
	bq.opMutex.Lock()
	toSend := bq.opCache
	bq.opCache = []entities.OperationData{}
	bq.opMutex.Unlock()

	if len(toSend) == 0 {
		return nil
	}
	i := bq.Client.Dataset(schema).Table(operationTable).Inserter()

	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	err := i.Put(ctx, toSend)
	if err != nil {
		return nerrors.NewInternalErrorFrom(err, "error sending operation")
	}

	// empty cache
	bq.opCache = []entities.OperationData{}

	return nil
}

// SendDataLoop is a loop to send the data to the database.
func (bq *BigQueryProvider) SendDataLoop() {

	for range time.Tick(bq.sendTime) {
		if err := bq.SendLoginCache(); err != nil {
			log.Error().Str("error", err.Error()).Str("trace", nerrors.FromError(err).StackTraceToString()).
				Msg("error sending login data")
		}

		if err := bq.SendOperationCache(); err != nil {
			log.Error().Str("error", err.Error()).Str("trace", nerrors.FromError(err).StackTraceToString()).
				Msg("error sending operation data")
		}

	}
}