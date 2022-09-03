package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	conf := &DbConfig{
		Host: "127.0.0.1",
		Port: 3306,
		Username: "root",
		Password: "123456",
		Database: "test",
		DriverName: "mysql",
	}

	db, _ := NewDB(conf)
	err := db.Ping()
	assert.ErrorIs(t, err, nil, "except get nil but got ", err)
}
