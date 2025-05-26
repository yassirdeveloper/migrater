package migrater

import (
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/db"
)

type Migrater interface {
	Apply(db db.Database) errors.Error
	Diff(db db.Database) (string, errors.Error)
	Status(db db.Database) (string, errors.Error)
	Plan(db db.Database) (string, errors.Error)
}
