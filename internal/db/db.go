package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/config"
	"github.com/yassirdeveloper/migrater/internal/db/drivers"
	"github.com/yassirdeveloper/migrater/internal/schema"
	"github.com/yassirdeveloper/migrater/internal/utils"
)

type Database interface {
	Init() errors.Error
	DSN() utils.DSN
	Execute(string) errors.Error
	Query(string) (drivers.Result, errors.Error)
	Validate() []errors.Error
	Describe() string
}

func GetDatabase(config config.DatabaseConfig) (Database, errors.Error) {
	dsn, err := config.GetDSN()
	if err != nil {
		return nil, err
	}
	d := &SqlDatabase{
		DriverType: config.GetDriver(),
		Name:       config.GetName(),
		dsn:        *dsn,
	}
	err = d.Init()
	if err != nil {
		return nil, err
	}
	return d, nil
}

func LoadFromJSON(filePath string) (Database, errors.Error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot open file! %s", filePath))
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot read file: %s", filePath))
	}
	var db SqlDatabase
	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid databse structure! %s", err))
	}
	return &db, nil
}

type SqlDatabase struct {
	DriverType drivers.DriverType `json:"driver"`
	Name       string             `json:"name"`
	Tables     []schema.Table     `json:"tables"`
	driver     drivers.Driver
	dsn        utils.DSN
}

func (d *SqlDatabase) Init() errors.Error {
	driver := drivers.GetDriver(d.DriverType)
	if driver == nil {
		return errors.New(fmt.Sprintf("Invalid driver: %s", d.DriverType))
	}
	d.driver = driver
	err := d.driver.Connect(d.DSN())
	if err != nil {
		return err
	}
	tableNames, err := d.driver.GetTableNames()
	if err != nil {
		return err
	}
	d.Tables = make([]schema.Table, 0, len(tableNames))
	for _, tableName := range tableNames {
		table, err := d.driver.GetTable(tableName)
		if err != nil {
			return err
		}
		d.Tables = append(d.Tables, table)
	}
	return nil
}

func (d *SqlDatabase) DSN() utils.DSN {
	return d.dsn
}

func (d *SqlDatabase) Execute(query string) errors.Error {
	return d.driver.Execute(query)
}

func (d *SqlDatabase) Query(query string) (drivers.Result, errors.Error) {
	return d.driver.Query(query)
}

func (d *SqlDatabase) Describe() string {
	tablesSummary := ""
	for _, table := range d.Tables {
		tablesSummary += fmt.Sprintf("Table: %s\n", table.Name)
		for _, column := range table.Columns {
			tablesSummary += fmt.Sprintf("  Column: %s, Type: %s\n", column.Name, column.Type)
			if len(column.Constraints) > 0 {
				tablesSummary += fmt.Sprintf("    Constraints: %s\n", column.Constraints)
			}
		}
	}
	if tablesSummary == "" {
		tablesSummary = "No tables found."
	}
	return fmt.Sprintf("Database: %s\nDriver: %s\nTables:\n%s", d.Name, d.DriverType, tablesSummary)
}

func (s *SqlDatabase) Validate() []errors.Error {
	errs := make([]errors.Error, 0)
	driver := drivers.GetDriver(s.DriverType)
	if driver == nil {
		if s.DriverType == "" {
			errs = append(errs, errors.New("missing driver"))
		} else {
			errs = append(errs, errors.New(fmt.Sprintf("invalid driver: %s", s.DriverType)))
		}
	}
	if s.Name == "" {
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
			if !drivers.HasType(driver, column.Type) {
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
