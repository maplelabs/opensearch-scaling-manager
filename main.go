package main

import (
	"fmt"
	"scaling_manager/cluster"
	"scaling_manager/config"
	"scaling_manager/provision"
	"scaling_manager/task"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var State = new(provision.State)

func main() {
	// A periodic check if there is a change in master node to pick up incomplete provisioning
	go periodicProvisionCheck()
	// The polling interval is set to 5 minutes and can be configured.
	go fileWatch()
	ticker := time.Tick(300 * time.Second)
	for range ticker {
		// This function is responsible for fetching the metrics and pushing it to the index.
		// In starting we will call simulator to provide this details with current timestamp.
		// fetch.FetchMetrics()
		// This function will be responsible for parsing the config file and fill in task_details struct.
		var task = new(task.TaskDetails)
		configStruct, err := config.GetConfig("config.yaml")
		if err != nil {
			fmt.Println("The recommendation can not be made as there is an error in the validation of config file.")
			fmt.Println(err)
			continue
		} else {
			fmt.Println(configStruct)
		}
		task.Tasks = configStruct.TaskDetails
		// This function is responsible for evaluating the task and recommend.
		recommendation_list := task.EvaluateTask()
		// This function is responsible for getting the recommendation and provision.
		provision.GetRecommendation(State, recommendation_list)
	}
}

// Input:
// Description: It periodically checks if the master node is changed and picks up if there was any ongoing provision operation
// Output:

func periodicProvisionCheck() {
	tick := time.Tick(30 * time.Second)
	previous_master := cluster.GetCurrentMasterIp()
	for range tick {
		current_state := State.GetCurrentState()
		// Call a function which returns the current master node
		current_master := cluster.GetCurrentMasterIp()
		if current_state != "normal" {
			if cluster.CheckIfMaster() {
				if previous_master != current_master {
					// Create a new command struct and call the scaleIn or scaleOut functions
					// Call these scaleOut and scaleIn functions using goroutines so that this periodic check continues
					// command struct to be filled with the recommendation queue and config file
					var command provision.Command
					if strings.Contains(current_state, "scaleup") {
						fmt.Println("Calling scaleOut")
						isScaledUp := command.ScaleOut(1, State)
						if isScaledUp {
							fmt.Println("Scaleup completed successfully")
						} else {
							// Add a retry mechanism
							fmt.Println("Scaleup failed")
						}
					} else if strings.Contains(current_state, "scaledown") {
						fmt.Println("Calling scaleIn")
						isScaledDown := command.ScaleIn(1, State)
						if isScaledDown {
							fmt.Println("Scaledown completed successfully")
						} else {
							// Add a retry mechanism
							fmt.Println("Scaledown failed")
						}
					}
				}
			}
		}
		// Update the repvious_master for next loop
		previous_master = current_master
	}
}

func fileWatch() {
	//Adding file watcher to detect the change in configuration
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()
	done := make(chan bool)

	//A go routine that keeps checking for change in configuration
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				//If there is change in config then clear recommendation queue
				//clearRecommendationQueue()
				fmt.Printf("EVENT! %#v\n", event)
				fmt.Println("The recommendation queue will be cleared.")
			case err := <-watcher.Errors:
				fmt.Println("ERROR in file watcher", err)
			}
		}
	}()

	// Adding fsnotify watcher to keep track of the changes in config file
	if err := watcher.Add("config.yaml"); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}
