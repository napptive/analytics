package config

import (
	"github.com/napptive/nerrors/pkg/nerrors"
	"github.com/rs/zerolog/log"
	"time"
)

// BigQueryConfig with the bigQuery provider configuration
type BigQueryConfig struct {
	// ProjectID with the GKE project identifier
	ProjectID string
	// Schema with the BigQuery Schema
	Schema string
	// Table with the BigQuery table name
	Table string
	// CredentialsPath with the Service Account credentials file
	CredentialsPath string
	// SendingTime is the time between database inserts
	SendingTime time.Duration
}

func NewBigQueryConfig(projectId string, schema string, table string, credentialsPath string, loopTime time.Duration) BigQueryConfig {
	return BigQueryConfig{
		ProjectID:       projectId,
		Schema:          schema,
		Table:           table,
		CredentialsPath: credentialsPath,
		SendingTime:     loopTime,
	}
}

func (bqc *BigQueryConfig) IsValid() error {
	if bqc.ProjectID == "" {
		return nerrors.NewFailedPreconditionError("ProjectID must be filled")
	}
	if bqc.Schema == "" {
		return nerrors.NewFailedPreconditionError("Schema must be filled")
	}
	if bqc.Table == "" {
		return nerrors.NewFailedPreconditionError("Table must be filled")
	}
	if bqc.CredentialsPath == "" {
		return nerrors.NewFailedPreconditionError("CredentialsPath path must be informed")
	}
	if bqc.SendingTime <= 0 {
		return nerrors.NewFailedPreconditionError("SendingTime must be filled")
	}
	return nil
}

func (bqc *BigQueryConfig) Print() {
	// Use logger to print the configuration
	log.Info().
		Str("ProjectID", bqc.ProjectID).
		Str("Schema", bqc.Schema).
		Str("Table", bqc.Table).
		Str("CredentialsPath", bqc.CredentialsPath).
		Dur("LoopTime", bqc.SendingTime).
		Msg("BigQuery config")
}
