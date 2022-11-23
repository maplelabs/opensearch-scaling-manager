// This package will fetch the recommendation from the recommendation Queue and provision the scale in/out
// based on command.
package provision

import (
	"scaling_manager/cluster"
	"time"
)

var counter uint8 = 1

// This struct contains the State of the opensearch scaling manager
// States can be following:
//  1. provision
//  2. provisioning
//  3. provisioning_completed
//  4. provisioning_failed
//  5. provisioned_successfully
type State struct {
	// CurrentState indicate the current state of the scaling manager
	CurrentState string
	// PreviousState indicates the previous state of the scaling manager
	PreviousState string
	// Remark indicates the additional remarks for the state of the scaling manager
	Remark string
}

// This struct contains the OS Admin Username and OS Admin Password via which we can connect to OS cluster.
type OsCredentials struct {
	// OsAdminUsername indicates the OS Admin Username via which OS client can connect to OS Cluster.
	OsAdminUsername string `yaml:"os_admin_username"`
	// OsAdminPassword indicates the OS Admin Password via which OS client can connect to OS Cluster.
	OsAdminPassword string `yaml:"os_admin_password"`
}

// This struct contains the Cloud Secret Key and Access Key via which we can connect to the cloud.
type CloudCredentials struct {
	// SecretKey indicates the Secret key for connecting to the cloud.
	SecretKey string `yaml:"secret_key"`
	// AccessKey indicates the Access key for connecting to the cloud.
	AccessKey string `yaml:"access_key"`
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
	OsCredentials    OsCredentials    `yaml:"os_credentials"`
	CloudCredentials CloudCredentials `yaml:"cloud_credentials"`
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
	state := getState()
	setState("provisioning", state.CurrentState)
	if c.Operation == "scale_up" {
		isScaledUp := c.scaleOut(1)
		if isScaledUp {
			state = getState()
			setState("provision_completed", state.CurrentState)
			checkClusterHealth()
		} else {
			state = getState()
			// Add a retry mechanism
			setState("provision_failed", state.CurrentState)
		}
	} else if c.Operation == "scale_down" {
		isScaledDown := c.scaleIn(1)
		if isScaledDown {
			state = getState()
			setState("provision_completed", state.CurrentState)
			checkClusterHealth()
		} else {
			state = getState()
			// Add a retry mechanism
			setState("provision_failed", state.CurrentState)
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
func (c *Command) scaleOut(numNodes int) bool {
	// Spin new VMs based on number of nodes and cloud type
	// Add the newly added VM to the list of VMs
	// Configure OS on newly created VM
	// Check cluster status after the configuration
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
func (c *Command) scaleIn(numNodes int) bool {
	// Identify the node which can be removed from the cluster.
	// Configure OS to tell master node that the present node is going to be removed
	// Wait for cluster to be in stable state(Shard rebalance)
	// Shut down the node
	// Check cluster status after shutting down the node
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
	if cluster.ClusterDynamic.ClusterStatus == "green" &&
		cluster.ClusterDynamic.NumRelocatingShards == 0 {
		setState("provisioned_successfully", "provision_completed")
	} else if counter >= 3 {
		time.Sleep(180 * time.Second)
		checkClusterHealth()
	} else {
		state := getState()
		setState("provisioned_failed", state.CurrentState)
	}
}
