package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

var DB *sqlx.DB

type SSLMode string

const (
	SSLModeEnable  SSLMode = "enable"
	SSLModeDisable SSLMode = "disable"
)

// InitDB function connects with a given database and returns error if there is any error
func InitDB() error {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), SSLModeDisable)
	DB, err = sqlx.Open("postgres", connStr)

	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}
	return migrateUp(DB)
}

func ShutdownDatabase() error {
	return DB.Close()
}

// migrateUp function migrate the database and handles the migration logic
func migrateUp(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations/",
		"postgres", driver)

	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Tx provides the transaction wrapper
func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %+v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				logrus.Errorf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logrus.Errorf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}
