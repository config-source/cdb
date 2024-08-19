package postgresutils

// It would be nice to put this behind a build tag but it breaks so much auto
// completion and helpfulness that I've decided to leave this dead code in the
// binary.

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	postgres "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresPooler interface {
	Raw() *pgxpool.Pool
}

type TestRepository struct {
	conn      *pgx.Conn
	Repo      PostgresPooler
	testName  string
	TestDBURL string
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func findMigrationPath() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	curDir, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	for {
		migrationPath := filepath.Join(curDir, "migrations")
		stat, err := os.Stat(migrationPath)
		if os.IsNotExist(err) || !stat.IsDir() {
			curDir = filepath.Dir(curDir)
			if curDir == "/" {
				panic(fmt.Errorf("Unable to find migrations directory from %s!", path))
			}

			continue
		} else if err != nil {
			panic(err)
		}

		return fmt.Sprintf("file://%s", migrationPath)
	}
}

var migrationPath = findMigrationPath()

// Start creates a test database, migrates it. Make sure to defer
// TestRepository.Cleanup() after calling Start.
//
// Use TestRepository.TestDBURL to connect to this new database.
func (tr *TestRepository) Start(testName string) error {
	port := os.Getenv("PGPORT")
	if port == "" {
		port = "5432"
	}

	host := os.Getenv("PGHOST")
	if host == "" {
		host = "localhost"
	}

	connUrlPrefix := fmt.Sprintf("postgres://%s:%s", host, port)

	connUrl := fmt.Sprintf("%s/%s", connUrlPrefix, "postgres")
	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		return err
	}
	tr.conn = conn
	tr.testName = toSnakeCase(testName)

	ctx := context.Background()

	tr.Cleanup()

	_, err = tr.conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", tr.testName))
	if err != nil {
		return err
	}

	tr.TestDBURL = fmt.Sprintf("%s/%s", connUrlPrefix, tr.testName)
	db, err := sql.Open("pgx", tr.TestDBURL)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	return m.Up()
}

// Cleanup deletes the test database.
func (tr *TestRepository) Cleanup() {
	if tr.Repo != nil {
		tr.Repo.Raw().Close()
	}

	_, err := tr.conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", tr.testName))
	if err != nil {
		panic(err)
	}
}

func InitTestDB(t *testing.T) *TestRepository {
	t.Parallel()
	t.Helper()

	tr := TestRepository{}
	err := tr.Start(t.Name())
	if err != nil {
		t.Fatal(err)
	}

	return &tr
}
