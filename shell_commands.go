package main

import (
	"github.com/abiosoft/ishell"
	"time"
	"fmt"
	"log"
	"github.com/horneth/gromwell"
)

var shellCommands = []*ishell.Cmd{
	&ishell.Cmd{
		Name: "watch",
		Help: "get status updates for a workflow",
		Func: func(c *ishell.Context) {
			go watchWorkflow(c.Args[0], c)
		},
	},
}

func watchWorkflow(workflowId string, c *ishell.Context) {
	enhancedId, err := enhancedWorkflowId(workflowId)

	if (err != nil) {
		c.Println("Cannot find last submitted workflow")
		return
	}
	
	ticker := time.NewTicker(time.Second * 10)
	
	var currentStatus = gromwell.WorkflowStatus {
		Id: workflowId,
		Status: "Submitted",
	}
	
	c.Println(fmt.Sprintf("Watching %s", enhancedId))
	
	for range ticker.C {
		status, err := cromwellClient.GetWorkflowStatus(enhancedId)
		if (err != nil) {
			log.Print(fmt.Sprintf("Error getting status for workflow %s: %s", workflowId, err))
			return
		} else if (status.Status != currentStatus.Status) {
			currentStatus = status
			fmt.Println(currentStatus)
		}
		if (currentStatus.Status == "Succeeded" || currentStatus.Status == "Failed") {
			ticker.Stop()
			return
		}
	}
}