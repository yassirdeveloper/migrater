package schema

import (
	"fmt"
	"strings"
)

type DataType string

func (t DataType) Equals(other DataType) bool {
	return strings.EqualFold(string(t), string(other))
}

func (t DataType) String() string {
	return string(t)
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Name        string       `json:"name"`
	Type        DataType     `json:"type"`
	Default     string       `json:"default,omitempty"`
	Constraints []Constraint `json:"constraints"`
}

type Constraint interface {
	Name() string
}

type NotNullConstraint struct{}

func (c NotNullConstraint) Name() string {
	return "NotNull"
}

type PrimaryKeyConstraint struct{}

func (c PrimaryKeyConstraint) Name() string {
	return "PrimaryKey"
}

type UniqueConstraint struct{}

func (c UniqueConstraint) Name() string {
	return "Unique"
}

type CheckConstraint struct{}

func (c CheckConstraint) Name() string {
	return "Check"
}

type DefaultConstraint struct {
	Value string `json:"value"`
}

func (c DefaultConstraint) Name() string {
	return fmt.Sprintf("Default(%s)", c.Value)
}

type ForeignKeyConstraint struct {
	ReferencedTable  string `json:"referenced_table"`
	ReferencedColumn string `json:"referenced_column"`
	OnDelete         string `json:"on_delete,omitempty"`
	OnUpdate         string `json:"on_update,omitempty"`
}

func (c ForeignKeyConstraint) Name() string {
	return fmt.Sprintf("ForeignKey(%s.%s)", c.ReferencedTable, c.ReferencedColumn)
}
