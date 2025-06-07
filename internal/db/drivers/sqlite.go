package drivers

import (
	"database/sql"
	"fmt"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/schema"
	"github.com/yassirdeveloper/migrater/internal/utils"
)

type sqliteDriver struct {
	version   float32
	dataTypes []schema.DataType
	db        *sql.DB
	*sqlite3.SQLiteDriver
}

func (d *sqliteDriver) GetDataTypes() []schema.DataType {
	return d.dataTypes
}

func (d *sqliteDriver) Connect(dsn utils.DSN) errors.Error {
	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			return errors.NewUnexpectedError(err)
		}
	}
	db, err := sql.Open("sqlite3", dsn.String())
	if err != nil {
		return errors.New(fmt.Sprintf("Could not establish connection!\n%s", err))
	}
	if err := db.Ping(); err != nil {
		return errors.New(fmt.Sprintf("Could not ping database!\n%s", err))
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

func (d *sqliteDriver) GetTableNames() ([]string, errors.Error) {
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, errors.NewUnexpectedError(err)
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, errors.NewUnexpectedError(err)
		}
		tableNames = append(tableNames, tableName)
	}
	return tableNames, nil
}

func (d *sqliteDriver) GetTable(tableName string) (schema.Table, errors.Error) {
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	rows, err := d.db.Query(query)
	if err != nil {
		return schema.Table{}, errors.NewUnexpectedError(err)
	}
	defer rows.Close()

	var columns []schema.Column
	for rows.Next() {
		var cid int
		var name string
		var ctype string
		var notnull int
		var dfltValue sql.NullString
		var pk int

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			return schema.Table{}, errors.NewUnexpectedError(err)
		}

		column := schema.Column{
			Name:       name,
			Type:       schema.DataType(ctype),
		}
		columns = append(columns, column)
	}

	return schema.Table{
		Name:    tableName,
		Columns: columns,
	}, nil
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
	dataTypes: []schema.DataType{
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
