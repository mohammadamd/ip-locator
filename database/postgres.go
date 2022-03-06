package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	configs "simple-fh/config"
)

var connection *sql.DB

// InitializePostgres will initialize the database connection
func InitializePostgres(cfg configs.Database) {
	d, err := sql.Open("postgres",
		fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s TimeZone=%s sslmode=disable",
			cfg.Port,
			cfg.Host,
			cfg.Username,
			cfg.Password,
			cfg.DBName,
			cfg.Timezone,
		),
	)
	if err != nil {
		logrus.Panicln(err)
	}

	d.SetMaxOpenConns(cfg.MaxOpenConnections)
	d.SetMaxIdleConns(cfg.MaxIdleConnections)

	if err := d.Ping(); err != nil {
		panic(err)
	}

	connection = d
}

// GetDatabaseConnection will return the database connection
// If it wasn't initialized, it will panic
func GetDatabaseConnection() *sql.DB {
	if connection == nil {
		panic("database connection is not initialized yet")
	}

	return connection
}
