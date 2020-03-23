package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDatabase migrates given swl version files to given host
func MigrateDatabase(schemaLoacation, URL string) error {
	m, err := migrate.New(
		schemaLoacation,
		URL,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		return err
	}
	return nil
}
