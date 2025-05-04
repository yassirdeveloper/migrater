package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/db/drivers"
	"github.com/yassirdeveloper/migrater/internal/utils"
)

type Database interface {
	Connect() error
	Disconnect() error
	Execute(query string) error
	Query(query string) (results []map[string]interface{}, err error)
}

type Schema struct {
	Driver       string  `json:"driver"`
	DatabaseName string  `json:"database"`
	SchemaName   string  `json:"schema"`
	Tables       []Table `json:"tables"`
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Name        string           `json:"name"`
	Type        drivers.DataType `json:"type"`
	Constraints []string         `json:"constraints"`
}

func LoadJSONSchema(filePath string) (Schema, errors.Error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Schema{}, errors.New(fmt.Sprintf("Cannot open file: %s", filePath))
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return Schema{}, errors.New(fmt.Sprintf("Cannot read file: %s", filePath))
	}
	var schema Schema
	err = json.Unmarshal(data, &schema)
	if err != nil {
		return Schema{}, errors.New(fmt.Sprintf("Invalid databse structure: %s", err))
	}
	return schema, nil
}

func (s *Schema) Validate() []errors.Error {
	errs := make([]errors.Error, 0)
	driver := drivers.Drivers[s.Driver]
	if driver == nil {
		if s.Driver == "" {
			errs = append(errs, errors.New("missing driver"))
		} else {
			errs = append(errs, errors.New(fmt.Sprintf("invalid driver: %s", s.Driver)))
		}
	}
	if s.DatabaseName == "" {
		errs = append(errs, errors.New("database name cannot be empty"))
	}

	if len(s.Tables) == 0 {
		errs = append(errs, errors.New("schema must have at least one table"))
	}

	tableNames := make(map[string]bool)
	for _, table := range s.Tables {
		if tableNames[table.Name] {
			errs = append(errs, errors.New(fmt.Sprintf("duplicate table name: %s", table.Name)))
		}
		tableNames[table.Name] = true

		err := utils.ValidateSQLName(table.Name)
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("invalid table name: %s (%s)", table.Name, err.Display())))
		}

		columnNames := make(map[string]bool)
		for _, column := range table.Columns {
			if columnNames[column.Name] {
				errs = append(errs, errors.New(fmt.Sprintf("duplicate column name: %s in table %s", column.Name, table.Name)))
			}
			columnNames[column.Name] = true

			err = utils.ValidateSQLName(column.Name)
			if err != nil {
				errs = append(errs, errors.New(fmt.Sprintf("invalid column name: %s (%s)", column.Name, err.Display())))
			}

			if column.Type == "" {
				errs = append(errs, errors.New(fmt.Sprintf("type of column %s cannot be empty", column.Name)))
			}

			if !slices.Contains(driver.GetDataTypes(), column.Type) {
				errs = append(errs, errors.New(fmt.Sprintf("invalid data type: %s for column %s", column.Type, column.Name)))
			}

			// Add constraints validation here
			// for _, constraint := range column.Constraints {
			// 	if !isValidConstraint(constraint) {
			// 		errs = append(errs, errors.New(fmt.Sprintf("invalid constraint: %s for column %s", constraint, column.Name)))
			// 	}
			// }
		}
	}
	return errs
}
