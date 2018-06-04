package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	// PostgreSQL adapter
	_ "github.com/lib/pq"
)

const (
	dbDialect = "postgres"
)

var (
	sqlDB *sql.DB
)

// ConnectionParams ...
type ConnectionParams struct {
	Host     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

// SetDB - you should NOT use this function directly,
// this is mainly exposed for testing purposes.
// If you want to open a database connection you should do it
// through InitAndOpenDatabase
func SetDB(sdb *sql.DB) {
	sqlDB = sdb
}

// Close ...
func Close() {
	closeDB(sqlDB)
	SetDB(nil)
}

func closeDB(dbToClose *sql.DB) {
	if dbToClose != nil {
		if err := dbToClose.Close(); err != nil {
			// impossible (to test - don't try to get to 100% coverage ;), the only case when there would be an error here
			// is if the database would fail, which is out of scope for unit/integration tests
			// see: https://golang.org/src/database/sql/sql.go?s=17459:17486#L612
			log.Printf(" [!] Exception: Failed to close DB: %+v", err)
		}
	}
}

// NewConnectionParamsFromEnvs ...
func NewConnectionParamsFromEnvs(defaults ConnectionParams) ConnectionParams {
	connParams := defaults

	if connParams.Host == "" {
		connParams.Host = os.Getenv("DB_HOST")
	}
	if connParams.DBName == "" {
		connParams.DBName = os.Getenv("DB_NAME")
	}
	if connParams.User == "" {
		connParams.User = os.Getenv("DB_USER")
	}
	if connParams.Password == "" {
		connParams.Password = os.Getenv("DB_PSW")
	}
	if connParams.SSLMode == "" {
		connParams.SSLMode = os.Getenv("DB_SSL_MODE")
	}

	return connParams
}

// Validate ...
func (cp ConnectionParams) Validate() error {
	if cp.Host == "" {
		return errors.New("No Database Host specified")
	}
	if cp.DBName == "" {
		return errors.New("No Database Name specified")
	}
	if cp.User == "" {
		return errors.New("No Database User specified")
	}
	if cp.Password == "" {
		return errors.New("No Database Password specified")
	}
	return nil
}

// ToConnectionStringParam ...
func (cp ConnectionParams) ToConnectionStringParam() (string, error) {
	if err := cp.Validate(); err != nil {
		return "", err
	}

	// every required input found!
	dbConnectionParams := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		cp.Host, cp.DBName, cp.User, cp.Password)

	// optionals
	if cp.SSLMode != "" {
		dbConnectionParams += " sslmode=" + cp.SSLMode
	}

	return dbConnectionParams, nil
}

// InitAndOpenDatabaseWithConnectionParams ...
func InitAndOpenDatabaseWithConnectionParams(dbConnectionParams string) error {
	db, err := sql.Open(dbDialect, dbConnectionParams)
	if err != nil {
		// impossible (to test - don't try to get to 100% coverage ;), the only case when there would be an error here
		// is if the driver is incorrect, connection params are not validated at all at this point
		// see: https://golang.org/src/database/sql/sql.go?s=16039:16096#L558
		return errors.Wrap(err, "Failed to open database")
	}

	if err = db.Ping(); err != nil {
		closeDB(db)
		return errors.Wrap(err, "Failed to ping database")
	}

	SetDB(db)
	return nil
}

// InitAndOpenDatabase ...
func InitAndOpenDatabase() error {
	dbConnectionParams := os.Getenv("DATABASE_URL")
	if dbConnectionParams == "" {
		log.Println(" (!) No DATABASE_URL specified! Checking DB connection params environment variables ...")

		connParams := NewConnectionParamsFromEnvs(ConnectionParams{})
		connParamsStr, err := connParams.ToConnectionStringParam()
		if err != nil {
			return errors.Wrap(err, "No DATABASE_URL specified and parsing params from envs resulted in an error")
		}
		dbConnectionParams = connParamsStr
	}

	return InitAndOpenDatabaseWithConnectionParams(dbConnectionParams)
}
