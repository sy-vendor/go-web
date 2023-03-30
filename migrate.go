package go_web

import (
	"embed"
	"errors"
	"os"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

// version defines the current migration version, this ensures the app
// is always compatible with the version of the database.
const version = 202303301713

// Migrate migrates the Postgres schema to the current version.
func Migrate() error {
	sourceInstance, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if !strings.Contains(databaseURL, "multiStatements") {
		if strings.Contains(databaseURL, "?") {
			databaseURL += "&multiStatements=true"
		} else {
			databaseURL += "?multiStatements=true"
		}
	}
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	targetInstance, err := mysql.WithInstance(db.DB(), &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", sourceInstance, "nft", targetInstance)
	if err != nil {
		return err
	}
	err = m.Migrate(version) // current version
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return sourceInstance.Close()
}
