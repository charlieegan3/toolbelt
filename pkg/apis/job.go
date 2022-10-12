package apis

import (
	"context"
	"time"
)

// Job is an interface which defines a cron job to be run by the toolbelt's job runner
type Job interface {
	// Name returns the name of the job for display purposes
	Name() string

	// Run runs the job
	Run(ctx context.Context) error

	// Timeout returns the time after which the job should be preemptively killed
	Timeout() time.Duration

	// Schedule returns the crontab schedule for the job
	Schedule() string
}

// ExternalJobRunner is an interface which defines a runner for jobs outside the toolbelt. This is
// used for jobs which need other binaries, more compute/ram etc.
type ExternalJobRunner interface {
	Configure(config map[string]any) error
	RunJob(job *ExternalJob) error
}

// ExternalJob is an interface for a job to run on by an ExternalJobRunner
type ExternalJob interface {
	// Name returns the name of the job for display purposes
	Name() string

	// Config returns the configuration for the job to be handed to the external job runner
	Config() map[string]any
}
