package provision

import (
	"fmt"
	"scaling_manager/cluster"
)

// Input:
// Description:
//
//	getRecommendation will fetch the recommendation from recommendation queue and clear the queue.
//	It will populate the command queue which contains all the details to scale out the cluster.
//
// Return:
//
//	getRecommendation will return the command struct which will contains all the details to scale out the cluster.
func getRecommendation() Command {
	var command Command
	// Fetch the recommendation from the queue
	// Queue can be stored as document in OS, localy stored data structure or in cache memory.
	// Clear the recommendation queue
	return command
}

// Input:
// Description:
//
//	triggerRecommendation will fetch the recommendation from the queue, get the status of the provisioner
//	and cluster and trigger the provisioning.
//
// Return:
func triggerRecommendation() {
	command := getRecommendation()
	clusterCurrent := cluster.GetClusterCurrent()
	state := getState()
	if clusterCurrent.ClusterDynamic.ClusterStatus == "green" ||
		clusterCurrent.ClusterDynamic.NumRelocatingShards == 0 {
		if state.CurrentState != "provisioning" && state.CurrentState != "provision" {
			setState("provision", state.CurrentState)
			go command.Provision()
		} else {
			fmt.Println("Recommendation can not be provisioned as open search cluster is already in provisioning phase.")
		}
	}
}

// Input:
// Description:
//
//	getState will get the state of provisioning system of the scaling manager.
//
// Return:
//
//	Return the State struct populated by the function which contains the current state and previous state.
func getState() State {
	var state State
	return state
}

// Input:
//
//	currentState(string): The current state for the provisioner.
//	previousState(string): The previous state for the provisioner.
//
// Description:
//
//	setState will set the state of provisioning system of the scaling manager.
//
// Return:
func setState(currentState string, previousState string) {
	// set the state for the opensearch scaling manager
	// This state can be either pushed to OS or else kept locally.
	var state State
	state.CurrentState = currentState
	state.PreviousState = previousState
}
