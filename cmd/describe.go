package cmd

import (
	"fmt"

	"github.com/yassirdeveloper/cli/command"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/cli/operator"
	"github.com/yassirdeveloper/migrater/internal/config"
	"github.com/yassirdeveloper/migrater/internal/db"
)

var databaseOption = command.CommandOption{
	Name:        "database",
	Label:       "Database",
	Description: "Registered database name",
	Letter:      'd',
	ValueType:   command.TypeString,
}

func describeHandler(input command.CommandInput, operator operator.Operator) errors.Error {
	databaseOpt, err := input.ParseOption(databaseOption)
	if err != nil {
		return operator.Write(err.Display())
	}
	globalConfig, err := config.GetGlobalConfig()
	if err != nil {
		return operator.Write(err.Error())
	}
	var databaseConfig config.DatabaseConfig
	if databaseOpt != nil {
		databaseName := databaseOpt.(string)
		databaseConfig = globalConfig.GetDatabaseConfig(databaseName)
		if databaseConfig == nil {
			return operator.Write(fmt.Sprintf("Missing configuration for database: %s", databaseName))
		}
	} else {
		databaseConfig = globalConfig.GetDefaultDatabaseConfig()
		if databaseConfig == nil {
			return operator.Write("No default database is configured")
		}
	}
	database, err := db.GetDatabase(databaseConfig)
	if err != nil {
		return operator.Write(err.Display())
	}
	return operator.Write(database.Describe())
}

func DescribeCommand() command.Command {
	cmd := command.NewCommand(
		"describe",
		"Prints the strcuture of the database.",
		describeHandler,
	)
	cmd.AddOption(databaseOption)
	return cmd
}
