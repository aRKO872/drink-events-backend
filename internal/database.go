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
      return fmt.Errorf("database already initialized") // Database already initialized
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

func GetDB(connString string) (*gorm.DB, error) {
  if db == nil {
    if err := InitDB(connString); err != nil {
      return nil, fmt.Errorf("error initializing db")
    }
  }
  return db, nil
}