// Package tests contains supporting code for running tests.
package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"jnk-ardan-service/business/data/schema"
	"jnk-ardan-service/business/sys/database"
	"jnk-ardan-service/foundation/docker"
	"jnk-ardan-service/foundation/logger"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// DBContainer provides configuration for a container to run.
type DBContainer struct {
	Image string
	Port  string
	Args  []string
}

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test (cleanup function).
func NewUnit(t *testing.T, dbc DBContainer) (*zap.SugaredLogger, *sqlx.DB, func()) {
	// I dont want the log output to intermix with testing output
	// capture the logs during the unit test
	// any time logging to stdOut it will go into the os.Pipe()
	// read all logs out at the end
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	// if this fails, the test fails. Recall logic in foundation/docker package
	c := docker.StartContainer(t, dbc.Image, dbc.Port, dbc.Args...)

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Waiting for database to be ready ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// if Migrate or Seed fail, dump logs and stop container
	// helps to identify bugs
	if err := schema.Migrate(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(t, c.ID)
		t.Fatalf("Migrating error: %s", err)
	}

	if err := schema.Seed(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(r, c.ID)
		t.Fatalf("Seeding error: %s", err)
	}

	log, err := logger.New("TEST")
	if err != nil {
		t.Fatalf("logger error: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
		docker.StopContainer(t, c.ID)

		log.Sync()

		// read the logs out of the buffer and put them in their own section after the tests
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old
		fmt.Println("**************** LOGS ****************")
		fmt.Print(buf.String())
		fmt.Println("**************** LOGS ****************")
	}

	return log, db, teardown
}
