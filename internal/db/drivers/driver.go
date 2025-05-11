package drivers

import (
	"slices"
	"strings"

	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/config"
)

type DataType string

func (t DataType) Equals(other DataType) bool {
	return strings.EqualFold(string(t), string(other))
}

type Driver interface {
	GetDataTypes() []DataType
	Connect(config.ConnectionConfig) errors.Error
	Close() errors.Error
}

func HasType(d Driver, t DataType) bool {
	driverTypes := d.GetDataTypes()
	tIndex := slices.IndexFunc(
		driverTypes,
		func(s DataType) bool {
			return s.Equals(t)
		},
	)
	return tIndex != -1
}

var Drivers = map[string]Driver{
	"mysql":    MysqlDriver,
	"postgres": PostgresDriver,
	"sqlite":   SqliteDriver,
}
