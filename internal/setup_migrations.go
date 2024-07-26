package internal_database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	pkg_config "github.com/drink-events-backend/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func SetupMigrations() error {
	db, err := sql.Open("postgres", pkg_config.GetProjectConfig().DATABASE_URL)
	if err != nil {
		return errors.New("error setting up DB")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return errors.New("error creating PostgreSQL driver instance") 
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting current working directory:", err)
	}
	migrationsPath := filepath.Join(workingDir, "internal", "migrations")
	absoluteMigrationsPath := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(
			absoluteMigrationsPath,
			"postgres", 
			driver,
	)

	if err != nil {
		return errors.New("error creating migrate instance")
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("error setting up migrations : ", err)
	}

	return nil
}
