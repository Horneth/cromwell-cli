package main

import (
	"github.com/urfave/cli"
	"errors"
	"fmt"
	"github.com/horneth/gromwell"
)

func submitWorkflow(c *cli.Context) error {
	var wdlPath string
	var workflowInputs string
	var workflowOptions string
	
	if err := validateNbArgs(c, []int {1, 2, 3}); err != nil { return cli.NewExitError(err, 1) }

	wdlPath = c.Args().First()

	if (c.NArg() == 2) {
		workflowInputs = c.Args()[1]
	}

	if (c.NArg() == 3) {
		workflowOptions = c.Args()[2]
	}

	submitCommand := gromwell.SubmitCommand {
		WdlSource: wdlPath,
		WorkflowInputs: workflowInputs,
		WorkflowOptions: workflowOptions,
	}

	status, err := cromwellClient.SubmitWorkflow(submitCommand)

	if (err == nil) {
		last = status.Id
	}

	fmt.Println(status)

	return nil
}

func getStatus(c *cli.Context) error {
	if err := validateOneArg(c); err != nil { return cli.NewExitError(err, 1) }
	
	id, err := enhancedWorkflowId(c.Args().First())
	if (err != nil) { return cli.NewExitError(err, 1) }
	
	status, err := cromwellClient.GetWorkflowStatus(id)
	if (err != nil) { return cli.NewExitError(err, 1) }
	
	fmt.Println(status.String())

	return nil
}

func getOutputs(c *cli.Context) error {
	if err := validateOneArg(c); err != nil { return cli.NewExitError(err, 1) }

	id, err := enhancedWorkflowId(c.Args().First())
	if (err != nil) { return cli.NewExitError(err, 1) }
	
	outputs, err := cromwellClient.GetWorkflowOutputs(id)
	if (err != nil) { return cli.NewExitError(err, 1) }
	
	if output := c.String("output"); output != "" {
		if err = outputs.ToFile(output); err != nil {
			return cli.NewExitError(err, 1)
		}
		
		fmt.Printf("Outputs written to %s\n", output)
		return nil
	}

	fmt.Println(outputs.String())

	return nil
}

func getMetadata(c *cli.Context) error {
	// Validate nb args
	if err := validateOneArg(c); err != nil { return cli.NewExitError(err, 1) }

	id, err := enhancedWorkflowId(c.Args().First())
	if (err != nil) { return cli.NewExitError(err, 1) }
	
	metadata, err := cromwellClient.GetWorkflowMetadata(id)
	if (err != nil) { return cli.NewExitError(err, 1) }
	
	if output := c.String("output"); output != "" {
		if err = metadata.ToFile(output); err != nil {
			return cli.NewExitError(err, 1)
		}
		
		fmt.Printf("Metadata written to %s\n", output)
		return nil
	}

	fmt.Println(metadata.String())

	return nil
}

func abortWorkflow(c *cli.Context) error {
	if err := validateOneArg(c); err != nil { return cli.NewExitError(err, 1) }

	id, err := enhancedWorkflowId(c.Args().First())
	if (err != nil) { return cli.NewExitError(err, 1) }
	
	status, err := cromwellClient.AbortWorkflow(id)
	if (err != nil) { return cli.NewExitError(err, 1) }

	fmt.Println(status.String())

	return nil
}

func cromwellVersion(c *cli.Context) error {
	// Validate nb args
	if err := validateNoArg(c); err != nil { return err }
	
	version, err := cromwellClient.Version()

	if (err != nil) { return cli.NewExitError(err, 1) }
	
	fmt.Println(version)
	
	return nil
}

func enhancedWorkflowId(id string) (string, error) {
	if (id == "last") {
		if (last != "") {
			return last, nil
		} else {
			return "", errors.New("Cannot find last submitted workflow")
		}
	}
	return id, nil
}

func validateNbArgs(c *cli.Context, allowed []int) error {
	for _, v := range allowed {
		if (c.NArg() == v) { return nil }
	}
	cli.ShowCommandHelp(c, c.Command.Name)
	return errors.New(fmt.Sprintf("Invalid number of arguemnts: %d. Expecting one of %s", c.NArg(), conjoin(allowed)))
}

func validateOneArg(c *cli.Context) error {
	return validateNbArgs(c, []int {1})
}

func validateNoArg(c *cli.Context) error {
	return validateNbArgs(c, []int {0})
}
