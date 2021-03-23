package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cage1016/ms-sample/internal/pkg/gormtracing"
	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config defines the options that are used when connecting to a PostgreSQL instance
type Config struct {
	Host        string `envconfig:"QS_DB_HOST" default:"localhost"`
	Port        string `envconfig:"QS_DB_PORT" default:"5432"`
	User        string `envconfig:"QS_DB_USER" default:"postgres"`
	Pass        string `envconfig:"QS_DB_PASS" default:"password"`
	Name        string `envconfig:"QS_DB" default:"tictac"`
	SSLMode     string `envconfig:"QS_DB_SSL_MODE" default:"disable"`
	SSLCert     string `envconfig:"QS_DB_SSL_CERT"`
	SSLKey      string `envconfig:"QS_DB_SSL_KEY"`
	SSLRootCert string `envconfig:"QS_DB_SSL_ROOT_CERT"`
}

func (cfg Config) ToURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)
}

// Connect creates a connection to the PostgreSQL instance and applies any
// unapplied database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.ToURL()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Use(gormtracing.NewGormTracePlugin())

	err = db.AutoMigrate(
	// &model.TicTac{},
	)
	if err != nil {
		return nil, err
	}

	db2, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := MigrateDB(db2, migrate.Up); err != nil {
		return nil, err
	}
	return db, nil
}

//MigrateDB means to migrate prepared data to the DB.
func MigrateDB(db *sql.DB, direction migrate.MigrationDirection) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "tictacs_1",
				Up: []string{`
					CREATE TABLE IF NOT EXISTS tictacs (
						value         int
					);

					INSERT INTO tictacs VALUES (0);
					`,
				},
				Down: []string{
					"DROP TABLE tictacs;",
				},
			},
		},
	}
	_, err := migrate.Exec(db, "postgres", migrations, direction)
	return err
}
