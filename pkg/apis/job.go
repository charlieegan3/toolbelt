package apis

import "time"

type Job interface {
	// Name returns the name of the job for display purposes
	Name() string

	// Run runs the job
	Run() error

	// Timeout returns the time after which the job should be preemptively killed
	Timeout() time.Duration

	// Schedule returns the crontab schedule for the job
	Schedule() string
}
