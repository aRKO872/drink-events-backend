package dao

import (
	internal_database "github.com/drink-events-backend/internal"
	pkg_config "github.com/drink-events-backend/pkg/config"
)

type DatabaseAccessOperator struct {
	internal_database.DatabaseOperator
}

func GetDBAccessOperator() (
	*DatabaseAccessOperator,
	error,
) {
	databaseInit, getDBErr := internal_database.GetDB(pkg_config.GetProjectConfig().DATABASE_URL)
	if getDBErr != nil {
		return nil, getDBErr
	}

	return &DatabaseAccessOperator{
		*databaseInit,
	}, nil
}