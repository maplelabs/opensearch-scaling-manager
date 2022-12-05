package main

import (
	"fmt"
	"scaling_manager/config"
	provision "scaling_manager/provision"
	"scaling_manager/task"
	"time"
	"github.com/fsnotify/fsnotify"
)

func main() {
	// The polling interval is set to 5 minutes and can be configured.
	go fileWatch()
	ticker := time.Tick(5 * time.Second)
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
		//configJson,_ := json.Marshal(configStruct.ClusterDetails)
		arr := task.EvaluateTask()
		// This function is responsible for getting the recommendation and provision.
		provision.GetRecommendation(arr)
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