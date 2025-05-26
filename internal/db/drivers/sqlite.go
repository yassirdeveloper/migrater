package drivers

import (
	"database/sql"
	"fmt"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/utils"
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

func (d *sqliteDriver) Connect(dsn utils.DSN) errors.Error {
	if d.db != nil {
		return errors.New("connection already exists")
	}
	db, err := sql.Open("sqlite3", dsn.String())
	if err != nil {
		return errors.New(fmt.Sprintf("Could not establish connection!\n%s", err))
	}
	d.db = db
	return nil
}

func (d *sqliteDriver) Execute(query string) errors.Error {
	_, err := d.db.Exec(query)
	if err != nil {
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func (d *sqliteDriver) Query(query string) (Result, errors.Error) {
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, errors.NewUnexpectedError(err)
	}
	return rows, nil
}

func (d *sqliteDriver) Version() float32 {
	return d.version
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

var sqliteDriverInstance = &sqliteDriver{
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
