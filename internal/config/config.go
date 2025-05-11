package config

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/yassirdeveloper/cli/errors"
)

var globalConfigFilePath = "config.hcl"

type GlobalConfig interface {
	GetDefaultDatabaseName() string
	GetDatabaseConfig(string) *DatabaseConfig
	Validate() errors.Error
}

func GetGlobalConfig() (GlobalConfig, errors.Error) {
	configFile, err := os.Open(globalConfigFilePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot open configuration file: %s", globalConfigFilePath))
	}
	data, err := io.ReadAll(configFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot read configuration file: %s", globalConfigFilePath))
	}
	file, diags := hclsyntax.ParseConfig(data, globalConfigFilePath, hcl.InitialPos)
	if diags.HasErrors() {
		return nil, errors.New(fmt.Sprintf("Cannot parse configuration file: %s\n%s", globalConfigFilePath, diags.Error()))
	}
	var conf globalConfig
	if err := gohcl.DecodeBody(file.Body, nil, &conf); err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot parse configuration file: %s\n%s", globalConfigFilePath, err))
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

func (c *globalConfig) GetDefaultDatabaseName() string {
	for _, d := range c.Databases {
		if d.Default {
			return d.Name
		}
	}
	var nilString string
	return nilString
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
	}
	return nil
}

type DatabaseConfig struct {
	Name    string `hcl:",label"`
	Driver  string `hcl:"driver,attr"`
	Default bool   `hcl:"default,attr"`
	DSN     string `hcl:"dsn,attr"`
}

type ConnectionConfig interface {
	DSN() string
}

type StandardConnectionConfig struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
}

func (c *StandardConnectionConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.User, c.Password, c.Host, c.Port, c.Database)
}

type FileConfig struct {
	FilePath string
}

func (c *FileConfig) DSN() string {
	return c.FilePath
}
