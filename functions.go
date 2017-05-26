package main

import (
	"github.com/urfave/cli"
	"errors"
	"fmt"
	"strings"
	"strconv"
	"io/ioutil"
	"github.com/horneth/gromwell"
)

func submitWorkflow(c *cli.Context) (interface{}, error) {
	var wdlPath string
	var workflowInputs string
	var workflowOptions string
	var status gromwell.WorkflowStatus
	
	if err := validateNbArgs(c.NArg(), []int {1, 2, 3}); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return status, err
	}

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

	return status, err
}

func getStatus(c *cli.Context) (interface{}, error) {
	var status gromwell.WorkflowStatus
	if err := validateNbArgs(c.NArg(), []int {1}); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return status, err
	}
	status, err := cromwellClient.GetWorkflowStatus(enhancedWorkflowId(c.Args().First()))
	return status, err
}

func getOutputs(c *cli.Context) (interface{}, error) {
	var outputs gromwell.WorkflowOutputs
	if err := validateNbArgs(c.NArg(), []int {1}); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return outputs, err
	}
	outputs, err := cromwellClient.GetWorkflowOutputs(enhancedWorkflowId(c.Args().First()))
	return outputs, err
}

func getMetadata(c *cli.Context) (interface{}, error) {
	fmt.Println(c.Args())
	var metadata gromwell.WorkflowMetadata
	if err := validateNbArgs(c.NArg(), []int {1}); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return metadata, err
	}
	metadata, err := cromwellClient.GetWorkflowMetadata(enhancedWorkflowId(c.Args().First()))
	if (err != nil) { return metadata, err }
	if output := c.String("output"); output != "" {
		err = ioutil.WriteFile(output, metadata.Metadata, 0644)
		if (err == nil) { return fmt.Sprintf("Metadata written to %s ", output), err }
	}
	
	return metadata, err
}

func abortWorkflow(c *cli.Context) (interface{}, error) {
	var status gromwell.WorkflowStatus
	if err := validateNbArgs(c.NArg(), []int {1}); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return status, err
	}
	status, err := cromwellClient.AbortWorkflow(enhancedWorkflowId(c.Args().First()))
	return status, err
}

func cromwellVersion(c *cli.Context) (interface{}, error) {
	var version string
	if err := validateNbArgs(c.NArg(), []int {0}); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return version, err
	}
	version, err := cromwellClient.Version()
	return version, err
}

func enhancedWorkflowId(id string) string {
	if (id == "last" && last != "") { return last }
	return id
}

func validateNbArgs(nbArgs int, allowed []int) error {
	for _, v := range allowed {
		if (nbArgs == v) { return nil }
	}
	
	return errors.New(fmt.Sprintf("Invalid number of arguemnts: %d. Expecting one of %s", nbArgs, conjoin(allowed)))
}

func conjoin(items []int) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return strconv.Itoa(items[0])
	}
	if len(items) == 2 { // "a and b" not "a, and b"
		return strconv.Itoa(items[0]) + " " + "," + " " + strconv.Itoa(items[1])
	}

	sep := ", "
	pieces := []string{strconv.Itoa(items[0])}
	for _, item := range items[1 : len(items)-1] {
		pieces = append(pieces, sep, strconv.Itoa(item))
	}
	pieces = append(pieces, sep, strconv.Itoa(items[len(items)-1]))

	return strings.Join(pieces, "")
}
