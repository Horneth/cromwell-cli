package main

import (
	"github.com/abiosoft/ishell"
	"github.com/horneth/gromwell"
	"github.com/urfave/cli"
	"os"
	"net/url"
	"flag"
)

// Cli App - global because accessed by the shell too
var app *cli.App

var cromwellClient cromwell_api.CromwellClient
// Last workflow submitted - for convenience
var last string

func startShell(commands []*ishell.Cmd ) {
	shell := ishell.New()
	shell.SetPrompt("> ")
	shell.Println("Welcome to the Cromwell CLI !")

	for _, v := range commands {
		shell.AddCmd(v)
	}

	// start shell
	shell.Start()
}

// List of all CliCommands available
func makeShellCommands(cliContext *cli.Context) []*ishell.Cmd {
	result := make([]*ishell.Cmd, len(cliCommands))
	
	// The shell inherits from all the cli commands, so add them here
	for i, v := range cliCommands {
		result[i] = &ishell.Cmd {
			Name:   v.Name,
			Help:   v.Usage,
			Func: 	func(shellContext *ishell.Context) {
				// find the cli command
				command := app.Command(shellContext.Cmd.Name)
				// create a FlagSet and pars the args
				flagSet := flag.NewFlagSet(shellContext.Cmd.Name, flag.ContinueOnError)
				// FIXME: need to use the raw args here otherwise won't work for some cases
				// See https://github.com/abiosoft/ishell/pull/26
				flagSet.Parse(shellContext.Args)
				// Create a new cli context from the args
				newContext := cli.NewContext(app, flagSet, cliContext)
				// Delegate execution of the command to the Cli command
				err := cli.HandleAction(command.Action, newContext)
				if (err != nil) { shellContext.Err(err) }
			},
		}
	}
	
	result = append(result, shellCommands...)
	return result
}

func main() {
	app = cli.NewApp()
	app.Name = "Cromwell CLI"
	app.Usage = "Run your worfklows on the command line !"
	app.Version = "0.1a"

	app.Flags = []cli.Flag {
		cli.StringFlag {
			Name: "host",
			Value: "http://localhost:8000",
			Usage: "Cromwell server url",
		},
	}

	app.Before = func (c *cli.Context) error {
		cromwellHost := c.GlobalString("host")
		cromwellServer, err := url.Parse(cromwellHost)
		if (err != nil) { return err }
		cromwellClient = cromwell_api.CromwellClient { CromwellUrl: cromwellServer }
		return nil
	}

	consoleCommand := cli.Command {
		Name:    "console",
		Usage:   "start a console",
		Action:  func(c *cli.Context) error {
			startShell(makeShellCommands(c))
			return nil
		},
	}

	app.Commands = append(cliCommands, consoleCommand)

	app.Run(os.Args)
}
