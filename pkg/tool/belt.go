package tool

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/charlieegan3/toolbelt/pkg/apis"
)

// Belt is the main struct for the Tool Belt. It contains the base router which all tools are registered to
type Belt struct {
	Router *mux.Router

	server *http.Server

	config map[string]any

	db *sql.DB

	cron *cron.Cron
}

// NewBelt creates a new Belt struct with an initalized router
func NewBelt() *Belt {
	return &Belt{
		Router: mux.NewRouter(),
		cron:   cron.New(),
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

	if tool.FeatureSet().Database {
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
		for _, job := range tool.Jobs() {
			err := b.AddJob(job)
			if err != nil {
				return fmt.Errorf("failed to add job for tool %s: %w", tool.Name(), err)
			}
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

func (b *Belt) StartServer(host, port string) chan os.Signal {
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	return c
}

func (b *Belt) StopServer(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return b.server.Shutdown(ctx)
}

func (b *Belt) AddJob(job apis.Job) error {
	err := b.cron.AddFunc(
		job.Schedule(),
		func() {
			log.Printf("running job %s", job.Name())
			ctx, cancel := context.WithTimeout(context.Background(), job.Timeout())
			defer cancel()

			doneCh := make(chan error, 1)
			panicCh := make(chan interface{}, 1)

			go func() {
				defer func() {
					if p := recover(); p != nil {
						panicCh <- p
					}
				}()

				doneCh <- job.Run()
			}()

			select {
			case err := <-doneCh:
				if err != nil {
					log.Printf("error running job %s: %v", job.Name(), err)
				} else {
					log.Printf("ran job %s", job.Name())
				}
			case p := <-panicCh:
				log.Printf("error running job %s, panicked: %v", job.Name(), p)
			case <-ctx.Done():
				log.Printf("error running job %s, timed out", job.Name())
			}

		},
	)
	if err != nil {
		return fmt.Errorf("failed to add job %s to cron: %v", job.Name(), err)
	}

	return nil
}

func (b *Belt) RunJobs() {
	go func() {
		b.cron.Start()
	}()
}
