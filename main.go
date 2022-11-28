package main

import (
	"fmt"
	"scaling_manager/config"
	provision "scaling_manager/provision"
	"scaling_manager/task"
	cluster "scaling_manager/cluster"
	"time"
	"strings"
)

func main() {
	go periodicProvisionCheck()
	// The polling interval is set to 5 minutes and can be configured.
	ticker := time.Tick(300 * time.Second)
	for range ticker {
		// This function is responsible for fetching the metrics and pushing it to the index.
		// In starting we will call simulator to provide this details with current timestamp.
		// fetch.FetchMetrics()
		// This function will be responsible for parsing the config file and fill in task_details struct.
		var task = new(task.TaskDetails)
		configStruct := config.GetConfig("config.yaml")
		task.Tasks = configStruct.TaskDetails
		// This function is responsible for evaluating the task and recommend.
		task.EvaluateTask()
		// This function is responsible for getting the recommendation and provision.
		provision.GetRecommendation()
	}
}


// Input:
// Description: It periodically checks if the master node is changed and picks up if there was any ongoing provision operation
// Output:

func periodicProvisionCheck() {
        tick := time.Tick(10 * time.Second)
	previous_master := cluster.GetCurrentMasterIp()
        for range tick {
                state := provision.GetState()
                // Call a function which returns the current master node
                current_master := cluster.GetCurrentMasterIp()
                if strings.HasPrefix(state.CurrentState, "provision") {
	                if cluster.CheckIfMaster() {
                                if previous_master != current_master {
                                        // Create a new command struct and call the scaleIn or scaleOut functions
                                        // Call these scaleOut and scaleIn functions using goroutines so that this periodic check continues
                                        // command struct to be filled with the recommendation queue and config file
                                        var command provision.Command
                                        if strings.Contains(state.CurrentState, "scaleup") {
                                                fmt.Println("Calling scaleOut with goroutine")
                                                go command.ScaleOut(1)
                                        } else if strings.Contains(state.CurrentState, "scaledown") {
                                                fmt.Println("Calling scaleIn with goroutine")
                                                go command.ScaleIn(1)
                                        }
                                }
                        }
                }
                // Update the repvious_master for next loop
                previous_master = current_master
        }
}
