package infra

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresConfig is a configuration for postgres.
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	SSLCert  string
	SSLKey   string
}

// NewPostgres creates new connection to postgres and return the connection
func NewPostgres(cfg PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s TimeZone=Asia/Bangkok",
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.User,
		cfg.Password,
		cfg.SSLMode,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Silent),
	})
}
