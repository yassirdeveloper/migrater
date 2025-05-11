package cmd

import (
	"fmt"

	"github.com/yassirdeveloper/cli/command"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/cli/operator"
	"github.com/yassirdeveloper/migrater/internal/db"
)

var jsonFilePathArgument = command.CommandArgument{
	Label:       "filepath",
	Description: "Json file of the database schema structure",
	Position:    0,
	ValueType:   command.TypeString,
}

func validateHandler(input command.CommandInput, operator operator.Operator) errors.Error {
	filePathArg, err := input.ParseArgument(jsonFilePathArgument)
	if err != nil {
		err = operator.Write(err.Display())
		if err != nil {
			return errors.NewUnexpectedError(err)
		}
	}
	filePath := filePathArg.(string)
	schema, err := db.LoadJSONSchema(filePath)
	if err != nil {
		err = operator.Write(err.Display())
		if err != nil {
			return errors.NewUnexpectedError(err)
		}
		return nil
	}
	errs := schema.Validate()
	if len(errs) > 0 {
		err = operator.Write("Invalid database schema:\n")
		if err != nil {
			return errors.NewUnexpectedError(err)
		}
		for _, err := range errs {
			err = operator.Write(fmt.Sprintf("- %s\n", err.Display()))
			if err != nil {
				return errors.NewUnexpectedError(err)
			}
		}
		return nil
	}
	err = operator.Write("Valid schema!")
	if err != nil {
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func ValidateCommand() command.Command {
	cmd := command.NewCommand(
		"validate",
		"Validate the structure of database schema in the json file.",
		validateHandler,
	)
	cmd.AddArgument(jsonFilePathArgument)
	return cmd
}
