package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/utils"
)

type Database interface {
	Connect() error
	Disconnect() error
	Execute(query string) error
	Query(query string) (results []map[string]interface{}, err error)
}

type Schema struct {
	DatabaseName string  `json:"database"`
	SchemaName   string  `json:"schema" default:"public"`
	Tables       []Table `json:"tables"`
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Name        string   `json:"name"`
	DataType    string   `json:"type"`
	Constraints []string `json:"constraints"`
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

func (s *Schema) Validate() errors.Error {
	if s.DatabaseName == "" {
		return errors.New("database name cannot be empty")
	}

	if s.SchemaName == "" {
		return errors.New("database schema name cannot be empty")
	}

	if len(s.Tables) == 0 {
		return errors.New("schema must have at least one table")
	}

	for _, table := range s.Tables {
		err := utils.ValidateSQLName(table.Name)
		if err != nil {
			return errors.New(fmt.Sprintf("invalid table name: %s (%s)", table.Name, err.Display()))
		}
		for _, column := range table.Columns {
			err = utils.ValidateSQLName(table.Name)
			if err != nil {
				return errors.New(fmt.Sprintf("invalid column name: %s (%s)", column.Name, err.Display()))
			}
			if column.DataType == "" {
				return errors.New(fmt.Sprintf("type of column %s cannot be empty", column.Name))
			}
		}
	}
	return nil
}
