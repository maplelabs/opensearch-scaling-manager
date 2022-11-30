// This package will fetch the recommendation from the recommendation Queue and provision the scale in/out
// based on command.
package provision

import (
	"scaling_manager/cluster"
	"scaling_manager/config"
	"time"
	"fmt"
)

var counter uint8 = 1

// This struct contains the State of the opensearch scaling manager
// States can be of following types:
//  1. normal : This is the state when the recommnedation will be provisioned.
//  2. provision_scaleup/provision_scaledown : Once the trigger module will call provision it will set this state.
//  3. provisioning_scaleup/provisioning_scaledown : Once the provision module will start provisioning it will set this state.
//  4. provisioning_scaleup_completed/provisioning_scaledown_completed : Once the provision is completed then this state will be state.
//  5. provisioning_scaleup_failed/provisioning_scaledown_failed: If the provision is failed then this state will be set.
//  6. provisioned_successfully: If the provision is completed and cluster state is green then
//     this state will be set.
//  7. provisioned_failed: If the provision is completed and the cluster state is not green after
//     certain retries then this state will be set.
type State struct {
	// CurrentState indicate the current state of the scaling manager
	CurrentState string
	// PreviousState indicates the previous state of the scaling manager
	PreviousState string
	// Remark indicates the additional remarks for the state of the scaling manager
	Remark string
}

// This struct contains the operation and details to scale the cluster
type Command struct {
	// Operation indicates the operation will be performed by the provisioner.
	// As of now two operations can be performed by the provisioner:
	//  1) Scale up
	//  2) Scale down
	Operation string
	// NumNodes indicates the number of nodes need to be scaled in or out.
	NumNodes int
	cluster.ClusterStatic
	OsCredentials    config.OsCredentials    `yaml:"os_credentials"`
	CloudCredentials config.CloudCredentials `yaml:"cloud_credentials"`
}

// Input:
// Caller: Object of Command
// Description:
//
//	Provision will scale in/out the cluster based on the operation.
//	ToDo:
//		Think about the scenario where event based scaling needs to be performed.
//		Morning need to scale up and evening need to scale down.
//		If in morning the scale up was not successful then we should not perform the scale down.
//		May be we can keep a concept of minimum number of nodes as a configuration input.
//
// Return:
func (c *Command) Provision() {
	state := GetState()
	if c.Operation == "scale_up" {
		setState("provisioning_scaleup", state.CurrentState)
		isScaledUp := c.ScaleOut(1)
		if isScaledUp {
			state = GetState()
			setState("provisioning_scaleup_completed", state.CurrentState)
			checkClusterHealth()
		} else {
			state = GetState()
			// Add a retry mechanism
			setState("provisioning_scaleup_failed", state.CurrentState)
		}
	} else if c.Operation == "scale_down" {
		setState("provisioning_scaledown", state.CurrentState)
		isScaledDown := c.ScaleIn(1)
		if isScaledDown {
			state = GetState()
			setState("provisioning_scaledown_completed", state.CurrentState)
			checkClusterHealth()
		} else {
			state = GetState()
			// Add a retry mechanism
			setState("provisioning_scaledown_failed", state.CurrentState)
		}
	}
}

// Input:
//
//	numNodes(int): Number of nodes to scale out.
//
// Caller: Object of Command
// Description:
//
//		ScaleOut will scale out the cluster with the number of nodes.
//		This function will invoke commands to create a VM based on cloud type.
//	 	Then it will configure the opensearch on newly created nodes.
//
// Return:
//
//	Return the status of scale out of the nodes.
func (c *Command) ScaleOut(numNodes int) bool {
	// Read the current state of scaleup process and proceed with next step
	scaleup_stage := readStageFromEs()
	// If no stage was already set. The function returns an empty string. Then, start the scaleup process
        if scaleup_stage == "" {
		scaleup_stage = "start_scaleup_process"
	}
	// Spin new VMs based on number of nodes and cloud type
	if scaleup_stage == "start_scaleup_process" {
		fmt.Println("Spin new vms based on the cloud type")
		scaleup_stage = "scaleup_triggered_spin_vm"
	}
	// Add the newly added VM to the list of VMs
	// Configure OS on newly created VM
	if scaleup_stage == "scaleup_triggered_spin_vm" {
		fmt.Println("Check if the vm creation is complete and wait till done")
		fmt.Println("Add the spinned nodes into the list of vms")
		fmt.Println("Configure ES")
		scaleup_stage = "scaleup_configured"
	}
	// Check cluster status after the configuration
	if scaleup_stage == "scaleup_configured" {
		fmt.Println("Wait for the cluster health and return status")
		cluster := cluster.GetClusterCurrent()
		for i := 0 ; i <= 6; i++ {
			if cluster.ClusterDynamic.ClusterStatus == "green" {
				state := GetState()
				setState("provisioned_successfully", state.CurrentState)
				scaleup_stage = "scaleup_complete"
				break
			}
			time.Sleep(300 * time.Second)
		}
                state := GetState()
		if state.CurrentState != "provisioned_successfully" {
	                setState("provisioned_failed", state.CurrentState)
		}
	}
	return true
}

// Input:
//
//	numNodes(int): Number of nodes to scale in.
//
// Caller: Object of Command
// Description:
//
//	ScaleIn will scale in the cluster with the number of nodes.
//	This function will invoke commands to remove a node from opensearch cluster.
//
// Return:
//
//	Return the status of scale in of the nodes.
func (c *Command) ScaleIn(numNodes int) bool {
        // Read the current state of scaledown process and proceed with next step
        scaledown_stage := readStageFromEs()
	// If no stage was already set. The function returns an empty string. Then, start the scaledown process
	if scaledown_stage == "" {
		scaledown_stage = "start_scaledown_process"
	}

	// Identify the node which can be removed from the cluster.
	if scaledown_stage == "start_scaledown_process" {
		fmt.Println("Identify the node to remove from the cluster and store the node_ip")
		scaledown_stage = "scaledown_node_identified"
	}
	// Configure OS to tell master node that the present node is going to be removed
	if scaledown_stage == "scaledown_node_identified" {
		fmt.Println("Configure ES to remove the node ip from cluster")
                fmt.Println("Start ES service on the node")
                scaledown_stage = "scaledown_es_configured"
	}
	// Wait for cluster to be in stable state(Shard rebalance)
	// Shut down the node
	if scaledown_stage == "scaledown_es_configured" {
	        fmt.Println("Wait for the cluster to become healthy (in a loop) and then proceed")
	        fmt.Println("Shutdown the node")
		scaledown_stage = "scaledown_complete"
	}
	return true
}

// Input:
// Description:
//
//		checkClusterHealth will check the current cluster health and also check if there are any relocating
//	 	shards. If the cluster status is green and there are no relocating shard then we will update the status
//	 	to provisioned_successfully. Else, we will wait for 3 minutes and perform this check again for 3 times.
//
// Return:
func checkClusterHealth() {
	cluster := cluster.GetClusterCurrent()
	if cluster.ClusterDynamic.ClusterStatus == "green" {
		state := GetState()
		setState("provisioned_successfully", state.CurrentState)
	} else if counter >= 3 {
		time.Sleep(180 * time.Second)
		checkClusterHealth()
	} else {
		state := GetState()
		setState("provisioned_failed", state.CurrentState)
	}
	// We should wait for buffer period after provisioned_successfully state to stablize the cluster.
	// After that buffer period we should change the state to normal, which can tell trigger module to trigger
	// the recommendation.
}

// Input:
// Description:
//		Read the current stage that the provisioning process is in from Elasticsearch or any centralized DB which will be updated after each stage.
// Return: Stage returned from ES

func readStageFromEs() string {
	return ""
}
