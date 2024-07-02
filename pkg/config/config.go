package pkg_config

import (
	"os"
	"strconv"

	"github.com/drink-events-backend/models"
)

func GetProjectConfig () *models.ProjectConfig {
	smtpPort, _:= strconv.Atoi(os.Getenv("SMTP_PORT"))
	accessTokenExpiry, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	refreshTokenExpiry, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	fetchNearbyPeriod, _ := strconv.Atoi(os.Getenv("FETCH_NEARBY_ACTIVE_PERIOD"))

	return &models.ProjectConfig{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		PGADMIN_DEFAULT_EMAIL: os.Getenv("PGADMIN_DEFAULT_EMAIL"),
		PGADMIN_DEFAULT_PASSWORD: os.Getenv("PGADMIN_DEFAULT_PASSWORD"),
		REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
		POSTGRES_USER: os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB: os.Getenv("POSTGRES_DB"),
		SMTP_HOST: os.Getenv("SMTP_HOST"),
		SMTP_PASSWORD: os.Getenv("SMTP_PASSWORD"),
		SMTP_PORT: smtpPort,
		SMTP_USER: os.Getenv("SMTP_USER"),
		ACCESS_TOKEN_EXPIRY:  accessTokenExpiry,
		REFRESH_TOKEN_EXPIRY: refreshTokenExpiry,
		JWT_SECRET_KEY: os.Getenv("JWT_SECRET_KEY"),
		FE_SOURCE: os.Getenv("FE_SOURCE"),
		FETCH_NEARBY_ACTIVE_PERIOD: fetchNearbyPeriod,
	}
}