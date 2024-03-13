package subcommands

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/config-source/cdb/internal/settings"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"

	postgres "github.com/golang-migrate/migrate/v4/database/pgx/v5"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
)

var rollback bool
var steps int

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := sql.Open("pgx", settings.DBUrl())
		if err != nil {
			return err
		}

		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return err
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres", driver)
		if err != nil {
			return fmt.Errorf("failed to load migrations: %w", err)
		}

		var migrationErr error
		if steps > 0 {
			if rollback {
				steps = -1 * steps
			}

			migrationErr = m.Steps(steps)
		} else if rollback {
			migrationErr = m.Down()
		} else {
			migrationErr = m.Up()
		}

		if errors.Is(migrationErr, migrate.ErrNoChange) {
			fmt.Println("Database is up to date.")
			return nil
		}

		return migrationErr
	},
}

func init() {
	migrateCmd.Flags().BoolVarP(&rollback, "rollback", "r", false, "If provided rollback database migrations")
	migrateCmd.Flags().IntVarP(&steps, "steps", "s", -1, "Number of steps to migrate or rollback")

	rootCmd.AddCommand(migrateCmd)
}
