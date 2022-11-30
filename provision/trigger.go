package provision

import (
	"fmt"
	"scaling_manager/cluster"
)

// Input:
// Description:
//
//	GetRecommendation will fetch the recommendation from recommendation queue and clear the queue.
//	It will populate the command queue which contains all the details to scale out the cluster.
//
// Return:
func GetRecommendation() {
	var command Command
	clusterCurrent := cluster.GetClusterCurrent()
	state := GetState()
	// Fetch the recommendation from the queue
	// Queue can be stored as document in OS, localy stored data structure or in cache memory.
	if clusterCurrent.ClusterDynamic.ClusterStatus == "green" && state.CurrentState == "normal" {
		// Clear the recommendation queue
	}
	// Fill in the command struct with the recommendation queue and config file and trigger the recommendation.

	command.triggerRecommendation()
}

// Input:
// Description:
//
//	triggerRecommendation will fetch the recommendation from the queue, get the status of the provisioner
//	and cluster and trigger the provisioning.
//
// Return:
func (c *Command) triggerRecommendation() {
	clusterCurrent := cluster.GetClusterCurrent()
	state := GetState()
	if clusterCurrent.ClusterDynamic.ClusterStatus == "green" && state.CurrentState == "normal" {
		if c.Operation == "scale_up" {
			setState("provisioning_scaleup", state.CurrentState)
		} else if c.Operation == "scale_down" {
			setState("provisioning_scaledown", state.CurrentState)
		}
		go c.Provision()
	} else {
		fmt.Println("Recommendation can not be provisioned as open search cluster is already in provisioning phase or the cluster isn't healthy yet")
	}
}

// Input:
// Description:
//
//	GetState will get the state of provisioning system of the scaling manager.
//
// Return:
//
//	Return the State struct populated by the function which contains the current state and previous state.
func GetState() State {
	var state State
	state.CurrentState = "provisioning_scaledown"
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
