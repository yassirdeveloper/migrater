package drivers

import (
	"fmt"

	"github.com/jackc/pgx"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/config"
)

type postgresDriver struct {
	version   float32
	dataTypes []DataType
	conn      *pgx.Conn
}

func (d *postgresDriver) GetDataTypes() []DataType {
	return d.dataTypes
}

func (d *postgresDriver) Connect(c config.ConnectionConfig) errors.Error {
	if d.conn == nil {
		conf := c.(*config.StandardConnectionConfig)
		connConfig := pgx.ConnConfig{
			Host:     conf.Host,
			Port:     conf.Port,
			Database: conf.Database,
			User:     conf.User,
			Password: conf.Password,
		}
		conn, err := pgx.Connect(connConfig)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not establish connection: %s", err))
		}
		d.conn = conn
	}
	return nil
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

var PostgresDriver = &postgresDriver{
	version: 13.0,
	dataTypes: []DataType{
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
