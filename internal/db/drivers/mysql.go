package drivers

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/utils"
)

type mysqlDriver struct {
	version   float32
	dataTypes []DataType
	db        *sql.DB
	*mysql.MySQLDriver
}

func (d *mysqlDriver) GetDataTypes() []DataType {
	return d.dataTypes
}

func (d *mysqlDriver) Connect(dsn utils.DSN) errors.Error {
	if d.db != nil {
		errors.New("connection already exists")
	}
	db, err := sql.Open("mysql", dsn.String())
	if err != nil {
		return errors.New(fmt.Sprintf("Could not establish connection!\n%s", err))
	}
	d.db = db
	return nil
}

func (d *mysqlDriver) Execute(query string) errors.Error {
	_, err := d.db.Exec(query)
	if err != nil {
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func (d *mysqlDriver) Query(query string) (Result, errors.Error) {
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, errors.NewUnexpectedError(err)
	}
	return rows, nil
}

func (d *mysqlDriver) Version() float32 {
	return d.version
}

func (d *mysqlDriver) Close() errors.Error {
	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			return errors.NewUnexpectedError(err)
		}
		return nil
	}
	return nil
}

var mysqlDriverInstance = &mysqlDriver{
	version: 5.7,
	dataTypes: []DataType{
		"BIT",
		"TINYINT",
		"SMALLINT",
		"MEDIUMINT",
		"INT",
		"BIGINT",
		"FLOAT",
		"DOUBLE",
		"DECIMAL",
		"CHAR",
		"VARCHAR",
		"BINARY",
		"VARBINARY",
		"BLOB",
		"TEXT",
		"ENUM",
		"SET",
		"DATE",
		"TIME",
		"TIMESTAMP",
		"DATETIME",
		"YEAR",
		"GEOMETRY",
		"POINT",
		"LINESTRING",
		"POLYGON",
		"GEOMETRYCOLLECTION",
		"MULTIPOINT",
		"MULTILINESTRING",
		"MULTIPOLYGON",
		"GEOMETRYCOLLECTION",
	},
}
