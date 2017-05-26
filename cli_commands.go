package main

import (
	"github.com/urfave/cli"
	"fmt"
)


type WrongArgumentCount int
type CromwellAction func (c *cli.Context) (interface{}, error)

func (c WrongArgumentCount) Error() string {
	return fmt.Sprintf("Wrong number of argument: %d", int(c))
}

// CLI Command definition

var cliCommands = []cli.Command {
	cli.Command {
		Name: "run",
		Usage: "submit a workflow",
		ArgsUsage: "wdlSource [workflowInputs] [workflowOptions]",
		Action: cliAction(submitWorkflow),
	},
	cli.Command {
		Name: "status",
		Usage: "get the status of a workflow",
		ArgsUsage: "workflowId",
		Action: cliAction(getStatus),
	},
	cli.Command {
		Name: "outputs",
		Usage: "get the outputs of a workflow",
		ArgsUsage: "workflowId",
		Action: cliAction(getOutputs),
	},
	cli.Command {
		Name: "metadata",
		Usage: "get the metadata of a workflow",
		ArgsUsage: "workflowId",
		Flags: []cli.Flag {
			cli.StringFlag {
				Name: "output, o",
				Usage: "write the metadata to `FILE`",
			},
		},
		Action: cliAction(getMetadata),
	},
	cli.Command {
		Name: "abort",
		Usage: "abort a workflow",
		ArgsUsage: "workflowId",
		Action: cliAction(abortWorkflow),
	},
	cli.Command {
		Name: "cromwell-version",
		Usage: "return the version of the underlying cromwell server",
		ArgsUsage: "",
		Action: cliAction(cromwellVersion),
	},
}

// Method to get a cliAction from a cromwellCommand 
func cliAction(action CromwellAction) cli.ActionFunc {
	return func (c *cli.Context) error {
		result, err := action(c)
		if (err != nil) { return cli.NewExitError(err, 1) }
		if (result != nil) { fmt.Println(result) }
		return nil
	}
}
