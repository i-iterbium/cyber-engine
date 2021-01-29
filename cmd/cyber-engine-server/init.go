package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/i-iterbium/cyber-engine/internal/pkg/config"
	"github.com/i-iterbium/cyber-engine/internal/pkg/database"
	"github.com/i-iterbium/cyber-engine/internal/pkg/maputil"
	_ "github.com/lib/pq"
	"github.com/x-foby/go-short/log"
)

func init() {
	log.SetLevel(log.INFO)

	config.Init()
	registerDrivers()
}

func registerDrivers() {
	database.RegisterDriver("postgres", database.DriverSetting{
		GetConnectionString: getConnectionString,
		AfterConnection:     afterConnection,
	})
}

func getConnectionString(s map[string]interface{}) (string, error) {
	defaultHost := "localhost"
	defaultPort := 5432

	host, err := maputil.GetStringFromMap(s, "host", &defaultHost)
	if err != nil {
		return "", err
	}

	port, err := maputil.GetIntFromMap(s, "port", &defaultPort)
	if err != nil {
		return "", err
	}

	user, err := maputil.GetStringFromMap(s, "user", nil)
	if err != nil {
		return "", err
	}

	password, err := maputil.GetStringFromMap(s, "password", nil)
	if err != nil {
		return "", err
	}

	dbName, err := maputil.GetStringFromMap(s, "name", nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable binary_parameters=yes", host, port, user, password, dbName), nil
}

func afterConnection(db *sql.DB, cs database.ConnectionSetting) error {
	connMaxLifetime, err := maputil.GetIntFromMap(cs.ConnectionStringParams, "connMaxLifetime", nil)
	if err != nil {
		return err
	}

	maxIdleConns, err := maputil.GetIntFromMap(cs.ConnectionStringParams, "maxIdleConns", nil)
	if err != nil {
		return err
	}

	maxOpenConns, err := maputil.GetIntFromMap(cs.ConnectionStringParams, "maxOpenConns", nil)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Millisecond)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	return nil
}
