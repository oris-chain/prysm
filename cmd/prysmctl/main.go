package main

import (
	"os"

	"github.com/prysmaticlabs/prysm/cmd/prysmctl/get"
	"github.com/prysmaticlabs/prysm/cmd/prysmctl/ssz"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/prysmaticlabs/prysm/cmd/prysmctl/checkpoint"
)

var prysmctlCommands []*cli.Command

func main() {
	app := &cli.App{
		Commands: prysmctlCommands,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	prysmctlCommands = append(prysmctlCommands, checkpoint.Commands...)
	prysmctlCommands = append(prysmctlCommands, get.Commands...)
	prysmctlCommands = append(prysmctlCommands, ssz.Commands...)
}
