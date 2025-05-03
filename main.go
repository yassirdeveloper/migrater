package main

import (
	"log"

	cli "github.com/yassirdeveloper/cli"
	"github.com/yassirdeveloper/migrater/cmd"
)

const CLI_NAME = "migrater"
const CLI_VERSION = "0.0.1"

func main() {
	cli, err := cli.NewCli(CLI_NAME, CLI_VERSION)
	if err != nil {
		log.Fatal(err)
	}
	cli.AddCommand(cmd.ValidateCommand())
	cli.Run(true)
}
