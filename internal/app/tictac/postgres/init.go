package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required for SQL access
	migrate "github.com/rubenv/sql-migrate"
)

// Config defines the options that are used when connecting to a PostgreSQL instance
type Config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	Name        string
	SSLMode     string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
}

// Connect creates a connection to the PostgreSQL instance and applies any
// unapplied database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := MigrateDB(db, migrate.Up); err != nil {
		return nil, err
	}
	return db, nil
}

//MigrateDB means to migrate prepared data to the DB.
func MigrateDB(db *sqlx.DB, direction migrate.MigrationDirection) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "tictac_1",
				Up: []string{`
					CREATE TABLE IF NOT EXISTS tictac (
						value         int
					);

					INSERT INTO tictac VALUES (1);
					`,
				},
				Down: []string{
					"DROP TABLE tictac;",
				},
			},
		},
	}
	_, err := migrate.Exec(db.DB, "postgres", migrations, direction)
	return err
}
