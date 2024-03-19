package internal_database

import (
    "database/sql"
    "fmt"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var (
  db *gorm.DB
)

func InitDB(connString string) error {
  if db != nil {
      return nil // Database already initialized
  }

  sqlDB, err := sql.Open("pgx", connString)
  if err != nil {
      return fmt.Errorf("error opening SQL connection: %w", err)
  }

  gormDB, err := gorm.Open(postgres.New(postgres.Config{
      Conn: sqlDB,
  }), &gorm.Config{})
  if err != nil {
      return fmt.Errorf("error setting up GORM: %w", err)
  }

  db = gormDB
  return nil
}

func GetDB() (*gorm.DB, error) {
  if db == nil {
      return nil, fmt.Errorf("database not initialized")
  }
  return db, nil
}