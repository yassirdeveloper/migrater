package config

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/migrater/internal/db/drivers"
	"github.com/yassirdeveloper/migrater/internal/utils"
)

var globalConfigFilePath = "config.hcl"

type GlobalConfig interface {
	GetDefaultDatabaseConfig() *DatabaseConfig
	GetDatabaseConfig(string) *DatabaseConfig
	Validate() errors.Error
}

func GetGlobalConfig() (GlobalConfig, errors.Error) {
	configFile, err := os.Open(globalConfigFilePath)
	if err != nil {
		return nil, errors.New("Cannot open configuration file")
	}
	data, err := io.ReadAll(configFile)
	if err != nil {
		return nil, errors.New("Cannot read configuration file")
	}
	file, diags := hclsyntax.ParseConfig(data, globalConfigFilePath, hcl.InitialPos)
	if diags.HasErrors() {
		return nil, errors.New(fmt.Sprintf("Cannot parse configuration file!\n%s", diags.Error()))
	}
	var conf globalConfig
	if err := gohcl.DecodeBody(file.Body, nil, &conf); err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot decode configuration file!\n%s", err))
	}
	err = conf.Validate()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid configuration file: %s", err))
	}
	return &conf, nil
}

type globalConfig struct {
	Databases []*DatabaseConfig `hcl:"database,block"`
}

func (c *globalConfig) GetDefaultDatabaseConfig() *DatabaseConfig {
	for _, d := range c.Databases {
		if d.Default {
			return d
		}
	}
	return nil
}

func (c *globalConfig) GetDatabaseConfig(name string) *DatabaseConfig {
	for _, d := range c.Databases {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (c *globalConfig) Validate() errors.Error {
	// todo add more validation
	defaultCount := 0
	uniqueNames := make(map[string]bool)
	for _, d := range c.Databases {
		if d.Name == "" {
			return errors.New("Database name field is required")
		}
		if d.Driver == "" {
			return errors.New("Database driver field is required")
		}
		if d.DSN == "" {
			return errors.New("Database dsn field is required")
		}
		if uniqueNames[d.Name] {
			return errors.New(fmt.Sprintf("Duplicate database name found: %s", d.Name))
		}
		uniqueNames[d.Name] = true
		if d.Default {
			defaultCount++
		}
		driverName := d.Driver
		if driver := drivers.GetDriver(driverName); driver == nil {
			return errors.New(fmt.Sprintf("Invalid driver name: %s\nSupported drivers: %s", d.Driver, drivers.SupportedDrivers[:]))
		}
	}
	if defaultCount > 1 {
		return errors.New("Only one database can be set as default")
	}
	return nil
}

type DatabaseConfig struct {
	Name    string             `hcl:",label"`
	Driver  drivers.DriverType `hcl:"driver,attr"`
	Default bool               `hcl:"default,attr"`
	DSN     string             `hcl:"dsn,attr"`
}

func (d *DatabaseConfig) GetDSN() (utils.DSN, errors.Error) {
	var format utils.DSNFormat
	switch d.Driver {
	case drivers.MysqlDriverType:
		format = utils.DSNFormatMySQL
	case drivers.PostgresDriverType:
		format = utils.DSNFormatPostgres
	case drivers.SqliteDriverType:
		format = utils.DSNFormatSQLite
	default:
		return utils.DSN{}, errors.New(fmt.Sprintf("Unsupported driver: %s", d.Driver))
	}
	parsedDSN, err := utils.ToDSN(d.DSN, format)
	if err != nil {
		return utils.DSN{}, errors.New(fmt.Sprintf("Invalid DSN: %s", err))
	}
	return *parsedDSN, nil
}
