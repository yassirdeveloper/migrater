package drivers

import (
	"fmt"

	"github.com/jackc/pgx"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/schema"
	"github.com/yassirdeveloper/migrater/internal/utils"
)

type postgresDriver struct {
	version   float32
	dataTypes []schema.DataType
	conn      *pgx.Conn
}

func (d *postgresDriver) GetDataTypes() []schema.DataType {
	return d.dataTypes
}

func (d *postgresDriver) Connect(dsn utils.DSN) errors.Error {
	if d.conn == nil {
		connConfig := pgx.ConnConfig{
			Host:     dsn.Host,
			Port:     dsn.Port,
			Database: dsn.Database,
			User:     dsn.User,
			Password: dsn.Password,
		}
		conn, err_ := pgx.Connect(connConfig)
		if err_ != nil {
			return errors.New(fmt.Sprintf("Could not establish connection!\n%s", err_))
		}
		d.conn = conn
	}
	return nil
}

func (d *postgresDriver) Execute(query string) errors.Error {
	_, err := d.conn.Exec(query)
	if err != nil {
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func (d *postgresDriver) Query(query string) (Result, errors.Error) {
	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, errors.NewUnexpectedError(err)
	}
	return rows, nil
}

func (d *postgresDriver) GetTableNames() ([]string, errors.Error) {
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	rows, err := d.conn.Query(query)
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

func (d *postgresDriver) GetTable(tableName string) (schema.Table, errors.Error) {
	query := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = '%s'", tableName)
	rows, err := d.conn.Query(query)
	if err != nil {
		return schema.Table{}, errors.NewUnexpectedError(err)
	}
	defer rows.Close()

	var columns []schema.Column
	for rows.Next() {
		var column schema.Column
		if err := rows.Scan(&column.Name, &column.Type); err != nil {
			return schema.Table{}, errors.NewUnexpectedError(err)
		}
		columns = append(columns, column)
	}

	return schema.Table{
		Name:    tableName,
		Columns: columns,
	}, nil
}

func (d *postgresDriver) Version() float32 {
	return d.version
}

func (d *postgresDriver) Close() errors.Error {
	if d.conn != nil {
		err := d.conn.Close()
		if err != nil {
			return errors.NewUnexpectedError(err)
		}
		return nil
	}
	return nil
}

var postgresDriverInstance = &postgresDriver{
	version: 13.0,
	dataTypes: []schema.DataType{
		"SMALLINT",
		"INTEGER",
		"BIGINT",
		"DECIMAL",
		"NUMERIC",
		"REAL",
		"DOUBLE PRECISION",
		"SMALLSERIAL",
		"SERIAL",
		"BIGSERIAL",
		"MONEY",
		"BOOLEAN",
		"BYTEA",
		"CHAR",
		"VARCHAR",
		"CHARACTER",
		"TEXT",
		"DATE",
		"TIME",
		"TIMESTAMP",
		"TIMESTAMPTZ",
		"INTERVAL",
		"TIME WITH TIME ZONE",
		"TIMESTAMP WITH TIME ZONE",
		"NAME",
		"OID",
		"OIDVECTOR",
		"REGCLASS",
		"REGCONFIG",
		"REGDICTIONARY",
		"REGNAMESPACE",
		"REGOPER",
		"REGOPERATOR",
		"REGPROC",
		"REGPROCEDURE",
		"REGTYPE",
		"REGCLASS",
		"REGTYPE",
		"REGROLE",
		"REGTEMPLATE",
		"REGTRIGGER",
		"REGWINDOW",
		"UUID",
		"XML",
		"BOX",
		"CIRCLE",
		"LINE",
		"LSEG",
		"PATH",
		"POINT",
		"POLYGON",
		"INET",
		"CIDR",
		"MACADDR",
		"MACADDR8",
		"BIT",
		"BIT VARYING",
		"VARBIT",
		"TSVECTOR",
		"TSQUERY",
		"REGCONFIG",
		"REGDICTIONARY",
	},
}
