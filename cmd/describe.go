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
		err = operator.Write(err.Display())
		if err != nil {
			return errors.NewUnexpectedError(err)
		}
		return nil
	}
	globalConfig, err := config.GetGlobalConfig()
	if err != nil {
		err = operator.Write(err.Error())
		if err != nil {
			return errors.NewUnexpectedError(err)
		}
		return nil
	}
	var databaseConfig *config.DatabaseConfig
	if databaseOpt != nil {
		databaseName := databaseOpt.(string)
		databaseConfig := globalConfig.GetDatabaseConfig(databaseName)
		if databaseConfig == nil {
			err = operator.Write(fmt.Sprintf("Missing configuration for database: %s", databaseName))
			if err != nil {
				return errors.NewUnexpectedError(err)
			}
			return nil
		}
	} else {
		databaseConfig = globalConfig.GetDefaultDatabaseConfig()
		if databaseConfig == nil {
			err = operator.Write("No default database is configured")
			if err != nil {
				return errors.NewUnexpectedError(err)
			}
			return nil
		}
	}
	database, err := db.GetDatabase(databaseConfig)
	if err != nil {
		err = operator.Write(err.Display())
		if err != nil {
			return errors.NewUnexpectedError(err)
		}
		return nil
	}
	err = operator.Write(database.Describe())
	if err != nil {
		return errors.NewUnexpectedError(err)
	}
	return nil
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
