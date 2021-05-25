package config

import (
	"github.com/napptive/nerrors/pkg/nerrors"
	"time"
)

type BigQueryConfig struct {
	ProjectID       string
	CredentialsPath string
	LoopTime        time.Duration
}

func NewBigQueryConfig(projectId string, credentialsPath string, loopTime time.Duration) BigQueryConfig {
	return BigQueryConfig{
		ProjectID:       projectId,
		CredentialsPath: credentialsPath,
		LoopTime:        loopTime,
	}
}

func (bqc *BigQueryConfig) IsValid() error {
	if bqc.ProjectID == "" {
		return nerrors.NewFailedPreconditionError("projectID mus be filled")
	}
	if bqc.CredentialsPath == "" {
		return nerrors.NewFailedPreconditionError("credentials path must be informed")
	}
	if bqc.LoopTime <= 0 {
		return nerrors.NewFailedPreconditionError("loopTime mus be filled")
	}
	return nil
}
