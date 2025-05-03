package migrater

import "github.com/yassirdeveloper/migrater/internal/db"

type Migrater interface {
	Apply(schema db.Schema) error
	Diff(schema db.Schema) (diff string, err error)
	Status(schema db.Schema) (status string, err error)
	Plan(schema db.Schema) (plan string, err error)
}
