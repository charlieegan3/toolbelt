package tool

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"

	"github.com/charlieegan3/toolbelt/pkg/apis"
	utilsHTTP "github.com/charlieegan3/toolbelt/pkg/utils/http"
)

// Belt is the main struct for the Tool Belt. It contains the base router which all tools are registered to
type Belt struct {
	Router *mux.Router

	server *http.Server

	config map[string]any

	db *sql.DB

	jobs map[string][]apis.Job

	externalJobRunners map[string]apis.ExternalJobRunner
}

// NewBelt creates a new Belt struct with an initalized router
func NewBelt() *Belt {
	r := mux.NewRouter()
	r.Use(utilsHTTP.InitMiddlewareLogging())

	return &Belt{
		Router: r,
		jobs:   make(map[string][]apis.Job),
	}
}

// AddExternalJobRunner adds a new external job runner to the belt. Jobs can be run using this runner by referencing the runner's name
func (b *Belt) AddExternalJobRunner(runner apis.ExternalJobRunner) {
	if b.externalJobRunners == nil {
		b.externalJobRunners = make(map[string]apis.ExternalJobRunner)
	}

	b.externalJobRunners[runner.Name()] = runner

}

// ExternalJobsFunc returns a function which can be used to run jobs from external sources
func (b *Belt) ExternalJobsFunc() func(job apis.ExternalJob) error {
	return func(job apis.ExternalJob) error {
		runner, ok := b.externalJobRunners[job.RunnerName()]
		if !ok {
			return fmt.Errorf("failed to find runner %s", job.RunnerName())
		}

		return runner.RunJob(job)
	}
}

// AddTool adds a new tool to the belt. Each tool is given a subrouter with the base path set to the tool's HTTPPath
func (b *Belt) AddTool(tool apis.Tool) error {
	if tool.FeatureSet().Config {
		toolConfig, ok := b.config[tool.Name()]
		if !ok {
			return fmt.Errorf("tool %s requires config but none was provided", tool.Name())
		}

		err := tool.SetConfig(toolConfig.(map[string]any))
		if err != nil {
			return fmt.Errorf("failed to set config for tool %s: %w", tool.Name(), err)
		}
	}

	if tool.FeatureSet().ExternalJobs {
		tool.ExternalJobsFuncSet(b.ExternalJobsFunc())
	}

	if tool.FeatureSet().Database {
		if b.db == nil {
			return fmt.Errorf("tool %s requires a database but none was provided", tool.Name())
		}

		migrations, path, err := tool.DatabaseMigrations()
		if err != nil {
			return fmt.Errorf("failed to get database migrations for tool %s: %w", tool.Name(), err)
		}

		driver, err := postgres.WithInstance(b.db, &postgres.Config{
			MigrationsTable: fmt.Sprintf("schema_migrations_%s", strings.ReplaceAll(tool.Name(), "-", "_")),
		})
		if err != nil {
			return fmt.Errorf("failed to create database driver for tool %s: %w", tool.Name(), err)
		}

		source, err := iofs.New(migrations, path)
		m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
		if err != nil {
			return fmt.Errorf("failed to create database migrate instance for tool %s: %w", tool.Name(), err)
		}

		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run database migrations for tool %s: %w", tool.Name(), err)
		}

		tool.DatabaseSet(b.db)
	}

	if tool.FeatureSet().HTTP {
		path := tool.HTTPPath()
		if path == "" {
			return fmt.Errorf("tool %s cannot use the HTTP feature with a blank HTTPPath", tool.Name())
		}
		toolRouter := b.Router.PathPrefix(fmt.Sprintf("/%s", path)).Subrouter()
		err := tool.HTTPAttach(toolRouter)
		if err != nil {
			return fmt.Errorf("failed to attach tool: %v", err)
		}
	}

	if tool.FeatureSet().Jobs {
		loadedJobs, err := tool.Jobs()
		if err != nil {
			return fmt.Errorf("failed to get jobs for tool %s: %w", tool.Name(), err)
		}
		for _, job := range loadedJobs {
			b.AddJob(tool.Name(), job)
		}
	}

	return nil
}

func (b *Belt) SetConfig(config map[string]any) {
	b.config = config
}

func (b *Belt) SetDatabase(db *sql.DB) {
	b.db = db
}

func (b *Belt) DatabaseDownMigrate(tool apis.Tool) error {
	if b.db == nil {
		return fmt.Errorf("tool %s requires a database but none was provided", tool.Name())
	}

	migrations, path, err := tool.DatabaseMigrations()
	if err != nil {
		return fmt.Errorf("failed to get database migrations for tool %s: %w", tool.Name(), err)
	}

	driver, err := postgres.WithInstance(b.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver for tool %s: %w", tool.Name(), err)
	}

	source, err := iofs.New(migrations, path)
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create database migrate instance for tool %s: %w", tool.Name(), err)
	}

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run database down migrations for tool %s: %w", tool.Name(), err)
	}
	return nil
}

func (b *Belt) RunServer(ctx context.Context, host, port string) {
	b.server = &http.Server{
		Handler:      b.Router,
		Addr:         fmt.Sprintf("%s:%s", host, port),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	go func() {
		err := b.server.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down server")

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := b.server.Shutdown(ctx); err != nil {
			log.Fatalf("Graceful shutdown failed: %s", err)
		}
		log.Println("Server gracefully stopped")
	}
}

func (b *Belt) AddJob(toolName string, job apis.Job) {
	if _, ok := b.jobs[toolName]; !ok {
		b.jobs[toolName] = []apis.Job{}
	}

	b.jobs[toolName] = append(b.jobs[toolName], job)
}

func (b *Belt) RunJobs(ctx context.Context) {
	crn := cron.New()

	for toolName, jobs := range b.jobs {
		for i := range jobs {
			job := b.jobs[toolName][i]

			jobRef := fmt.Sprintf("%s/%s", toolName, job.Name())

			log.Printf("loaded job %q with schedule %q", jobRef, job.Schedule())

			err := crn.AddFunc(
				job.Schedule(),
				func() {
					log.Printf("running job %q", jobRef)
					jobCtx, cancel := context.WithTimeout(ctx, job.Timeout())
					defer cancel()

					doneCh := make(chan error, 1)
					panicCh := make(chan interface{}, 1)

					go func() {
						defer func() {
							if p := recover(); p != nil {
								panicCh <- p
							}
						}()

						doneCh <- job.Run(jobCtx)
					}()

					select {
					case err := <-doneCh:
						if err != nil {
							log.Printf("error running job %q: %v", jobRef, err)
						} else {
							log.Printf("ran job %q", jobRef)
						}
					case p := <-panicCh:
						log.Printf("error running job %q, panicked: %v", jobRef, p)
					case <-ctx.Done():
						if ctx.Err() == context.DeadlineExceeded {
							log.Printf("parent context timed out during job %q", jobRef)
						} else if ctx.Err() == context.Canceled {
							log.Printf("parent context cancelled during job %q", jobRef)
						}
					}
				},
			)
			if err != nil {
				log.Printf("failed to add job %q to cron: %v", jobRef, err)
			}
		}
	}

	log.Printf("job worker started")
	go func() {
		crn.Start()
	}()

	select {
	case <-ctx.Done():
		log.Println("stopping job worker")
		crn.Stop()
	}
}
