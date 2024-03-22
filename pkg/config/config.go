package pkg_config

import (
	"os"
	"strconv"
)

type ProjectConfig struct {
	DATABASE_URL string
	PGADMIN_DEFAULT_EMAIL string
	PGADMIN_DEFAULT_PASSWORD string
	REDIS_PASSWORD string
	POSTGRES_USER string
	POSTGRES_PASSWORD string
	POSTGRES_DB string
	SMTP_HOST string
	SMTP_PASSWORD string
	SMTP_PORT int
	SMTP_USER string
}

func GetProjectConfig () *ProjectConfig {
	smtpPort, _:= strconv.Atoi(os.Getenv("SMTP_PORT"))

	return &ProjectConfig{
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
	}
}