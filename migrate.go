package go_web

import (
	"embed"
	"errors"
	"fmt"

	"go-web/pkg/config"
	"go-web/pkg/log"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var fs embed.FS

// version defines the current migration version, this ensures the app
// is always compatible with the version of the database.
const version = 202303301713

// Migrate migrates the database schema to the current version.
func Migrate(cfg *config.Config) error {
	logger := log.NewLogger()
	logger.Info("starting database migration",
		zap.Uint("target_version", version),
		zap.String("database", cfg.Database.Database),
	)

	sourceInstance, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}
	defer sourceInstance.Close()

	db, err := sql.Open("mysql", cfg.Database.GetDSN())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	targetInstance, err := mysql.WithInstance(db.DB(), &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration target: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceInstance, "mysql", targetInstance)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	currentVersion, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	if dirty {
		return fmt.Errorf("database is in dirty state at version %d", currentVersion)
	}

	logger.Info("current database version",
		zap.Uint("version", currentVersion),
		zap.Bool("dirty", dirty),
	)

	err = m.Migrate(version) // current version
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("database is up to date")
			return nil
		}
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	logger.Info("database migration completed successfully",
		zap.Uint("from_version", currentVersion),
		zap.Uint("to_version", version),
	)

	return nil
}
