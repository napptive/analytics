package provider

import (
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/napptive/analytics/pkg/config"
	"github.com/napptive/analytics/pkg/entities"
	"github.com/napptive/nerrors/pkg/nerrors"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"sync"
	"time"
)

const databaseTimeout = time.Second * 10

// Provider with an interface that defines the monitoring provider methods
type Provider interface {
	// Send inserts the data in the database
	Send(data entities.Operation) error
	// Flush flush the data
	Flush() error
}

// BigQueryProvider with a provider that stores the operation information in bigquery
type BigQueryProvider struct {
	// Client with the client to connect to bigquery
	Client *bigquery.Client

	// BigQueryConfig with all the configuration (schema, table, etc.)
	config.BigQueryConfig

	// Cache with a cache of operation data
	// in this cache we store the operations to avoid send msg by msg
	Cache []entities.Operation
	// Mutex for managing login cache access.
	sync.Mutex
}

// NewBigQueryProvider
func NewBigQueryProvider(cfg config.BigQueryConfig) (Provider, error) {

	// validate config
	if err := cfg.IsValid(); err != nil {
		return nil, err
	}

	client, err := bigquery.NewClient(context.Background(), cfg.ProjectID,
		option.WithCredentialsFile(cfg.CredentialsPath))
	if err != nil {
		return nil, err
	}
	provider := &BigQueryProvider{
		Client:         client,
		BigQueryConfig: cfg,
		Cache:          []entities.Operation{},
		Mutex:          sync.Mutex{},
	}

	// start the loop
	go provider.LaunchSendingLoop()

	return provider, nil
}

// Send stores the operation information in a cache
func (bq *BigQueryProvider) Send(operation entities.Operation) error {
	bq.Lock()
	defer bq.Unlock()

	// store the login in the cache
	bq.Cache = append(bq.Cache, operation)

	return nil
}

func (bq *BigQueryProvider) Flush() error {
	return bq.SendCacheOperation()
}

// LaunchSendingLoop launch a loop to insert the cache data in the database
func (bq *BigQueryProvider) LaunchSendingLoop () {
	for range time.Tick(bq.SendingTime) {
		if err := bq.SendCacheOperation();
		err != nil{
			log.Error().Str("error", err.Error()).Str("trace", nerrors.FromError(err).StackTraceToString()).
				Msg("error sending operation data")
		}
	}
}

// SendCacheOperation inserts all the operations in the database
func (bq *BigQueryProvider) SendCacheOperation() error {
	bq.Lock()
	toSend := bq.Cache
	bq.Cache = []entities.Operation{}
	bq.Unlock()

	if len(toSend) == 0 {
		return nil
	}
	i := bq.Client.Dataset(bq.Schema).Table(bq.Table).Inserter()

	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	err := i.Put(ctx, toSend)
	if err != nil {
		return nerrors.NewInternalErrorFrom(err, "error sending operation")
	}

		return nil

}