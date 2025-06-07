package drivers

import (
	"slices"

	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/schema"
	"github.com/yassirdeveloper/migrater/internal/utils"
)

type DriverType string

const (
	MysqlDriverType    DriverType = "mysql"
	PostgresDriverType DriverType = "postgres"
	SqliteDriverType   DriverType = "sqlite"
)

var SupportedDrivers = []DriverType{
	MysqlDriverType,
	PostgresDriverType,
	SqliteDriverType,
}

func GetDriver(driverType DriverType) Driver {
	switch driverType {
	case MysqlDriverType:
		return mysqlDriverInstance
	case PostgresDriverType:
		return postgresDriverInstance
	case SqliteDriverType:
		return sqliteDriverInstance
	default:
		return nil
	}
}

type Result interface {
	Next() bool
	Scan(...any) error
}

type Driver interface {
	GetDataTypes() []schema.DataType
	Connect(utils.DSN) errors.Error
	Execute(string) errors.Error
	Query(string) (Result, errors.Error)
	Close() errors.Error
	Version() float32
	GetTableNames() ([]string, errors.Error)
	GetTable(string) (schema.Table, errors.Error)
}

func HasType(d Driver, t schema.DataType) bool {
	driverTypes := d.GetDataTypes()
	tIndex := slices.IndexFunc(
		driverTypes,
		func(s schema.DataType) bool {
			return s.Equals(t)
		},
	)
	return tIndex != -1
}
