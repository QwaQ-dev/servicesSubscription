package repository

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	"github.com/QwaQ-dev/servicesSubscription/internal/config"
	"github.com/QwaQ-dev/servicesSubscription/pkg/sl"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	migrateiofs "github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func InitDatabase(cfg config.Database, log *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBhost, cfg.Port, cfg.DBusername, cfg.DBpassword, cfg.DBname, cfg.SSLMode,
	))
	if err != nil {
		log.Error("Error with connecting to database", sl.Err(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error("Error with pinging database", sl.Err(err))
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Error("failed to create migrate driver", sl.Err(err))
		return nil, err
	}

	sourceDriver, err := migrateiofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Error("failed to create iofs source driver", sl.Err(err))
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		cfg.DBname,
		driver,
	)
	if err != nil {
		log.Error("failed to create migrate instance", sl.Err(err))
		return nil, err
	}

	runMigrations(m, "up", log)

	log.Info("Database connected successfully")
	return db, nil
}

func runMigrations(m *migrate.Migrate, direction string, log *slog.Logger) error {
	switch direction {
	case "up":
		log.Info("Running database migrations: UP")
		err := m.Up()
		if err != nil {
			if err == migrate.ErrNoChange {
				log.Info("No new migrations to apply")
				return nil
			}
			log.Error("Migration UP failed", sl.Err(err))
			return err
		}
		log.Info("Migrations UP applied successfully")

	case "down":
		log.Info("Running database migrations: DOWN")
		err := m.Down()
		if err != nil {
			if err == migrate.ErrNoChange {
				log.Info("No migrations to revert")
				return nil
			}
			log.Error("Migration DOWN failed", sl.Err(err))
			return err
		}
		log.Info("Migrations DOWN applied successfully")

	case "drop":
		log.Warn("Dropping all database objects")
		err := m.Drop()
		if err != nil {
			log.Error("Migration DROP failed", sl.Err(err))
			return err
		}
		log.Info("Database dropped successfully")

	default:
		return fmt.Errorf("unknown migration direction: %s", direction)
	}
	return nil
}
