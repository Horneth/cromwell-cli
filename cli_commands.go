package main

import (
	"github.com/urfave/cli"
	"fmt"
)


type WrongArgumentCount int
type CromwellAction func (c *cli.Context) (string, error)

func (c WrongArgumentCount) Error() string {
	return fmt.Sprintf("Wrong number of argument: %d", int(c))
}

// CLI Command definition

var cliCommands = []cli.Command {
	cli.Command {
		Name: "run",
		Usage: "submit a workflow",
		ArgsUsage: "wdlSource [workflowInputs] [workflowOptions]",
		Action: submitWorkflow,
	},
	cli.Command {
		Name: "status",
		Usage: "get the status of a workflow",
		ArgsUsage: "workflowId",
		Action: getStatus,
	},
	cli.Command {
		Name: "outputs",
		Usage: "get the outputs of a workflow",
		ArgsUsage: "workflowId",
		Action: getOutputs,
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
		Action: getMetadata,
	},
	cli.Command {
		Name: "abort",
		Usage: "abort a workflow",
		ArgsUsage: "workflowId",
		Action: abortWorkflow,
	},
	cli.Command {
		Name: "cromwell-version",
		Usage: "return the version of the underlying cromwell server",
		ArgsUsage: "",
		Action: cromwellVersion,
	},
}
