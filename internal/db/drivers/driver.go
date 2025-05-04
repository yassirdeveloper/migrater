package drivers

import (
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/config"
)

type DataType string

type Driver interface {
	GetDataTypes() []DataType
	Connect(config.Config) errors.Error
	Close() errors.Error
}

var Drivers = map[string]Driver{
	"mysql":    MysqlDriver,
	"postgres": PostgresDriver,
	"sqlite":   SqliteDriver,
}
