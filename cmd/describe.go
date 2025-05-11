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
	var databaseName string
	if databaseOpt != nil {
		databaseName = databaseOpt.(string)
	} else {
		globalConfig, err := config.GetGlobalConfig()
		if err != nil {
			err = operator.Write(err.Error())
			if err != nil {
				return errors.NewUnexpectedError(err)
			}
			return nil
		}
		databaseName = globalConfig.GetDefaultDatabaseName()
		if databaseName == "" {
			err = operator.Write("No default database is configured")
			if err != nil {
				return errors.NewUnexpectedError(err)
			}
			return nil
		}
	}
	database := db.GetDatabase(databaseName)
	if database == nil {
		err = operator.Write(fmt.Sprintf("No database is registered with the name: %s", databaseName))
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
		"Prints the strcuture of the database schema.",
		describeHandler,
	)
	cmd.AddOption(databaseOption)
	return cmd
}
