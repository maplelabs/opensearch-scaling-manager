package main

import (
	"scaling_manager/config"
	"scaling_manager/provision"
	"scaling_manager/task"
	"time"
)

func main() {
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
