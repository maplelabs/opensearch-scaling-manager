package main

import (
	"context"
	"scaling_manager/config"
	fetch "scaling_manager/fetchmetrics"
	"scaling_manager/logger"
	osutils "scaling_manager/opensearchUtils"
	"scaling_manager/provision"
	"scaling_manager/recommendation"
	utils "scaling_manager/utilities"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var state = new(provision.State)

var log logger.LOG

var firstExecution bool

func init() {
	log.Init("logger")
	log.Info.Println("Main module initialized")

	firstExecution = true

	configStruct, err := config.GetConfig("config.yaml", firstExecution, false)
	if err != nil {
		log.Panic.Println("The recommendation can not be made as there is an error in the validation of config file.", err)
		panic(err)
	}

	cfg := configStruct.ClusterDetails
	osutils.InitializeOsClient(cfg.OsCredentials.OsAdminUsername, cfg.OsCredentials.OsAdminPassword)

	provision.InitializeDocId()

	userCfg := configStruct.UserConfig

	if !userCfg.MonitorWithSimulator {
		go fetch.FetchMetrics(int(userCfg.PollingInterval))
	}
}

func main() {
	updatedCredsArray := make([]string, 4)

	configStruct, updatedCredsArray, err := config.UpdateEncryptedCred("config.yaml", firstExecution, updatedCredsArray)
	if err != nil {
		panic(err)
	}
	go fileWatch("config.yaml", updatedCredsArray)
	// A periodic check if there is a change in master node to pick up incomplete provisioning
	go periodicProvisionCheck(configStruct.UserConfig.PollingInterval)
	ticker := time.Tick(time.Duration(configStruct.UserConfig.PollingInterval) * time.Second)
	for range ticker {
		state.GetCurrentState()
		// The recommendation and provisioning should only happen on master node
		if utils.CheckIfMaster(context.Background(), "") && state.CurrentState == "normal" {
			//              if firstExecution || state.CurrentState == "normal" {
			firstExecution = false
			// This function will be responsible for parsing the config file and fill in task_details struct.
			var task = new(recommendation.TaskDetails)
			configStruct, err := config.GetConfig("config.yaml", false, false)
			if err != nil {
				log.Error.Println("The recommendation can not be made as there is an error in the validation of config file.")
				log.Error.Println(err.Error())
				continue
			}
			task.Tasks = configStruct.TaskDetails
			userCfg := configStruct.UserConfig
			clusterCfg := configStruct.ClusterDetails
			// This function is responsible for evaluating the task and recommend.
			recommendationList := task.EvaluateTask(userCfg.MonitorWithSimulator, userCfg.PollingInterval)
			// This function is responsible for getting the recommendation and provision.
			provision.GetRecommendation(state, recommendationList, clusterCfg, userCfg)
		}
	}
}

// Input:
// Description: It periodically checks if the master node is changed and picks up if there was any ongoing provision operation
// Output:

func periodicProvisionCheck(pollingInterval int) {
	tick := time.Tick(time.Duration(pollingInterval) * time.Second)
	previousMaster := utils.CheckIfMaster(context.Background(), "")
	for range tick {
		state.GetCurrentState()
		// Call a function which returns the current master node
		currentMaster := utils.CheckIfMaster(context.Background(), "")
		if state.CurrentState != "normal" {
			if (!previousMaster && currentMaster) || (currentMaster && firstExecution) {
				//                      if firstExecution {
				firstExecution = false
				configStruct, err := config.GetConfig("config.yaml", false, false)
				if err != nil {
					log.Warn.Println("Unable to get Config from GetConfig()", err)
					return
				}
				if strings.Contains(state.CurrentState, "scaleup") {
					log.Debug.Println("Calling scaleOut")
					isScaledUp := provision.ScaleOut(configStruct.ClusterDetails, configStruct.UserConfig, state)
					if isScaledUp {
						log.Info.Println("Scaleup completed successfully")
					} else {
						// Add a retry mechanism
						log.Warn.Println("Scaleup failed")
					}
				} else if strings.Contains(state.CurrentState, "scaledown") {
					log.Debug.Println("Calling scaleIn")
					isScaledDown := provision.ScaleIn(configStruct.ClusterDetails, configStruct.UserConfig, state)
					if isScaledDown {
						log.Info.Println("Scaledown completed successfully")
					} else {
						// Add a retry mechanism
						log.Warn.Println("Scaledown failed")
					}
				}
			}
		}
		// Update the previousMaster for next loop
		previousMaster = currentMaster
	}
}

func fileWatch(filePath string, updatedCredsArray []string) {
	//Adding file watcher to detect the change in configuration
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error.Println("Error while creating new fileWatcher : ", err)
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
				log.Info.Println("EVENT! ", event)
				time.Sleep(2 * time.Second)
				//fmt.Println("The recommendation queue will be cleared.")
				_, updatedArray, updateErr := config.UpdateEncryptedCred("config.yaml", false, updatedCredsArray)
				if updateErr != nil {
					log.Panic.Println("ConfigStruct encryption failed : ", updateErr)
				} else {
					updatedCredsArray = updatedArray
				}
			case err := <-watcher.Errors:
				log.Info.Println("Error in fileWatcher: ", err)
			}
		}
	}()

	// Adding fsnotify watcher to keep track of the changes in config file
	if err := watcher.Add(filePath); err != nil {
		log.Error.Println("Error while adding the config file changes to the fileWatcher :", err)
	}

	<-done
}
