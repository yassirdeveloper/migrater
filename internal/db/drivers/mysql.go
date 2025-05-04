package drivers

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/config"
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

func (d *mysqlDriver) Connect(c config.Config) errors.Error {
	if d.db != nil {
		errors.New("connection already exists")
	}
	db, err := sql.Open("mysql", c.DSN())
	if err != nil {
		return errors.New(fmt.Sprintf("Could not establish connection: %s", err))
	}
	d.db = db
	return nil
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

var MysqlDriver = &mysqlDriver{
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
