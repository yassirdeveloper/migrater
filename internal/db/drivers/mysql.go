package drivers

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/schema"
	"github.com/yassirdeveloper/migrater/internal/utils"
)

type mysqlDriver struct {
	version   float32
	dataTypes []schema.DataType
	db        *sql.DB
	*mysql.MySQLDriver
}

func (d *mysqlDriver) GetDataTypes() []schema.DataType {
	return d.dataTypes
}

func (d *mysqlDriver) Connect(dsn utils.DSN) errors.Error {
	if d.db != nil && d.db.Ping() == nil {
		return errors.New("connection already exists")
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

func (d *mysqlDriver) GetTableNames() ([]string, errors.Error) {
	query := "SHOW TABLES"
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
	if err := rows.Err(); err != nil {
		return nil, errors.NewUnexpectedError(err)
	}

	return tableNames, nil
}

func (d *mysqlDriver) GetTable(tableName string) (schema.Table, errors.Error) {
	query := fmt.Sprintf("SHOW COLUMNS FROM `%s`", tableName)
	rows, err := d.db.Query(query)
	if err != nil {
		return schema.Table{}, errors.NewUnexpectedError(err)
	}
	defer rows.Close()

	var columns []schema.Column
	for rows.Next() {
		var column schema.Column
		var isPrimaryKey string
		var columnType, isNullable, defaultValue, extra string
		if err := rows.Scan(&column.Name, &columnType, &isNullable, &isPrimaryKey, &defaultValue, &extra); err != nil {
			return schema.Table{}, errors.NewUnexpectedError(err)
		}
		column.Type = schema.DataType(column.Type)
		if isNullable == "NO" {
			column.Constraints = append(column.Constraints, schema.NotNullConstraint{})
		}
		if isPrimaryKey == "PRI" {
			column.Constraints = append(column.Constraints, schema.PrimaryKeyConstraint{})
		}
		if defaultValue != "" {
			column.Constraints = append(column.Constraints, schema.DefaultConstraint{Value: defaultValue})
		}
		columns = append(columns, column)
	}
	if err := rows.Err(); err != nil {
		return schema.Table{}, errors.NewUnexpectedError(err)
	}

	return schema.Table{
		Name:    tableName,
		Columns: columns,
	}, nil
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
	dataTypes: []schema.DataType{
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
