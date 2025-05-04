package drivers

import (
	"database/sql"
	"fmt"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/config"
)

type sqliteDriver struct {
	version   float32
	dataTypes []DataType
	db        *sql.DB
	*sqlite3.SQLiteDriver
}

func (d *sqliteDriver) GetDataTypes() []DataType {
	return d.dataTypes
}

func (d *sqliteDriver) Connect(c config.Config) errors.Error {
	if d.db != nil {
		return errors.New("connection already exists")
	}
	db, err := sql.Open("sqlite3", c.DSN())
	if err != nil {
		return errors.New(fmt.Sprintf("Could not establish connection: %s", err))
	}
	d.db = db
	return nil
}

func (d *sqliteDriver) Close() errors.Error {
	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			return errors.NewUnexpectedError(err)
		}
		return nil
	}
	return nil
}

var SqliteDriver = &sqliteDriver{
	version: 3.35,
	dataTypes: []DataType{
		"INTEGER",
		"REAL",
		"TEXT",
		"BLOB",
		"NULL",
		"DATE",
		"TIME",
		"DATETIME",
		"BOOLEAN",
	},
}
