//Package tests contains supporting code for running tests.
package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"ekraal.org/avatarlysis/business/data/schema"
	"ekraal.org/avatarlysis/foundation/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	//Success success color marker.
	Success = "\u2713"
	//Failed failure color marker.
	Failed = "\u2717"
)

const (
	dbImage = "postgres:12.2-alpine"
	dbPort  = "5432"
	user1ID = "5cf37266-3473-4006-984f-9325122678b7"
	user2ID = "45b5fbd3-755f-4379-8f07-a58d4a30fa2f"
)

//NewUnit creates a test database inside a docker container.It creates the
//required table structure but the database is otherwise empty.It returns
//the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T) (*sqlx.DB, func()) {
	c := startContainer(t, dbImage, dbPort)

	cfg := database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	}

	db, err := database.Open(cfg)
	if err != nil {
		t.Fatalf("opening database connection: %v", err)
	}

	t.Log("waiting for database to be ready...")

	//Wait for the database to be ready. Wait 100ms longer between each attempt.
	//Do not try more than 20 times..
	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		dumpContainerLogs(t, c.ID)
		stopContainer(t, c.ID)
		t.Fatalf("database never ready: %v", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		stopContainer(t, c.ID)
		t.Fatalf("migrating error: %s", err)
	}

	//teardown should be invoked when the caller is done with
	//the database.
	tearDown := func() {
		t.Helper()
		db.Close()
		stopContainer(t, c.ID)
	}

	return db, tearDown
}

//Test owns state for running and shutting down tests.
type Test struct {
	DB      *sqlx.DB
	Log     *log.Logger
	t       *testing.T
	cleanup func()
}

//NewIntegration creates a database and seeds it.
func NewIntegration(t *testing.T) *Test {
	db, cleanup := NewUnit(t)

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	log := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	test := Test{
		DB:      db,
		Log:     log,
		t:       t,
		cleanup: cleanup,
	}

	return &test
}

//Teardown releases any resources used for the test.
func (test *Test) Teardown() {
	test.cleanup()
}

//Context returns an app level context for testing.
func Context() context.Context {
	values := Values{
		TraceID: uuid.New().String(),
		Now:     time.Now(),
	}

	return context.WithValue(context.Background(), KeyValues, &values)
}

//represents the type of value for the context key.
type ctxKey int

//KeyValues is how request values are stored/retrieved.
const KeyValues ctxKey = 1

//Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}
