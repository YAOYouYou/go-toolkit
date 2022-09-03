package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DbConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string

	DriverName string
}

func NewDB(c *DbConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port, c.Database)
	db, err := sqlx.Open(c.DriverName, dsn)
	if err != nil {
		return nil, err
	}
	return db, err
}
