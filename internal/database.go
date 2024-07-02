package internal_database

import (
	"database/sql"
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseOperator struct {
	*sync.RWMutex
	DB *gorm.DB
}

var (
  dbOperator *DatabaseOperator
)

func InitDB(connString string) (*gorm.DB, error) {
  if dbOperator != nil {
      return nil, fmt.Errorf("database already initialized") // Database already initialized
  }

  sqlDB, err := sql.Open("pgx", connString)
  if err != nil {
      return nil, fmt.Errorf("error opening SQL connection: %w", err)
  }

  gormDB, err := gorm.Open(postgres.New(postgres.Config{
      Conn: sqlDB,
  }), &gorm.Config{})
  if err != nil {
      return nil, fmt.Errorf("error setting up GORM: %w", err)
  }

  return gormDB, nil
}

func GetDB(connString string) (*DatabaseOperator, error) {
  if dbOperator != nil {
		return dbOperator, nil;
	}

	// Initialize Redis DB
	dbOp, initializeErr := InitDB(connString)

	if initializeErr != nil {
		return nil, initializeErr
	}

	dbOperator = &DatabaseOperator{
		new(sync.RWMutex),
		dbOp, 
	}
	return dbOperator, nil;
}