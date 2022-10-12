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
	// Name returns the name of the job runner, used by jobs to select this runner
	Name() string
	// Configure sets the config for the job runner
	Configure(config map[string]any) error

	// RunJob takes an external job and runs it
	RunJob(job ExternalJob) error
}

// ExternalJob is an interface for a job to run on by an ExternalJobRunner
type ExternalJob interface {
	// Name returns the name of the job
	Name() string

	// RunnerName returns the name of the runner to use
	RunnerName() string

	// Config returns the configuration for the job to be handed to the external job runner
	Config() map[string]any
}
