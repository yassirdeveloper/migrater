package main

import (
	cli "github.com/yassirdeveloper/cli"
)

const CLI_NAME = "migrater"
const CLI_VERSION = "0.0.1"

func main() {
	cli := cli.NewCli(CLI_NAME, CLI_VERSION)
	cli.Run(true)
}
