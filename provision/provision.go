// This package will fetch the recommendation from the recommendation Queue and provision the scale in/out
// based on command.
package provision

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	ansibleutils "github.com/maplelabs/opensearch-scaling-manager/ansible_scripts"
	"github.com/maplelabs/opensearch-scaling-manager/cluster"
	"github.com/maplelabs/opensearch-scaling-manager/cluster_sim"
	"github.com/maplelabs/opensearch-scaling-manager/config"
	"github.com/maplelabs/opensearch-scaling-manager/crypto"
	osutils "github.com/maplelabs/opensearch-scaling-manager/opensearchUtils"
	utils "github.com/maplelabs/opensearch-scaling-manager/utilities"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/maplelabs/opensearch-scaling-manager/logger"

	"github.com/tkuchiki/faketime"
)

var log = new(logger.LOG)

// Input:
//
// Description:
//
//	Initialize the Provision module.
//
// Return:
func init() {
	log.Init("logger")
	log.Info.Println("Provisioner module initiated")
}

// Input:
//
//	clusterCfg (config.ClusterDetails): Cluster Level config details
//	usrCfg (config.UserConfig): User defined config for applicatio behavior
//	numNodes (int): Number of nodes to be scaled up/down
//	operation (string): scaleup or scaledown operation
//	RulesResponsible (string): A string that contains the rules responsible for the decision of operation being performed
//
// Description:
//
//	TriggerProvision will call scale in/out the cluster based on the operation.
//	ToDo:
//	        Think about the scenario where event based scaling needs to be performed.
//	        Morning need to scale up and evening need to scale down.
//	        If in morning the scale up was not successful then we should not perform the scale down.
//	        May be we can keep a concept of minimum number of nodes as a configuration input.
//
// Return:
func TriggerProvision(clusterCfg config.ClusterDetails, usrCfg config.UserConfig, numNodes int, t *time.Time, operation, RulesResponsible string) {
	state.GetCurrentState()
	if operation == "scale_up" {
		state.PreviousState = state.CurrentState
		state.CurrentState = "provisioning_scaleup"
		state.NumNodes = numNodes
		state.RemainingNodes = numNodes
		state.RuleTriggered = "scale_up"
		state.RulesResponsible = RulesResponsible
		state.UpdateState()
		isScaledUp, err := ScaleOut(clusterCfg, usrCfg, t)
		if isScaledUp {
			log.Info.Println("Scaleup successful")
			PushToOs("Success", err)
		} else {
			log.Error.Println(err)
			state.GetCurrentState()
			// Add a retry mechanism
			state.PreviousState = state.CurrentState
			state.CurrentState = "provisioning_scaleup_failed"
			state.UpdateState()
			PushToOs("Failed", err)
		}
		// Set the state back to normal to continue further
		SetStateBackToNormal()
	} else if operation == "scale_down" {
		state.PreviousState = state.CurrentState
		state.CurrentState = "provisioning_scaledown"
		state.NumNodes = numNodes
		state.RemainingNodes = numNodes
		state.RuleTriggered = "scale_down"
		state.RulesResponsible = RulesResponsible
		state.UpdateState()
		isScaledDown, err := ScaleIn(clusterCfg, usrCfg, t)
		if isScaledDown {
			log.Info.Println("Scaledown successful")
			PushToOs("Success", err)
		} else {
			log.Error.Println(err)
			state.GetCurrentState()
			// Add a retry mechanism
			state.PreviousState = state.CurrentState
			state.CurrentState = "provisioning_scaledown_failed"
			state.UpdateState()
			PushToOs("Failed", err)
		}
		// Set the state back to normal to continue further
		SetStateBackToNormal()
	}
}

// Input:
//
//	clusterCfg (config.ClusterDetails): Cluster Level config details
//	usrCfg (config.UserConfig): User defined config for applicatio behavior
//
// Description:
//
//	ScaleOut will scale out the cluster with the number of nodes.
//	This function will invoke commands to create a VM based on cloud type.
//	Then it will configure the opensearch on newly created nodes.
//
// Return:
//
//	(bool): Return the status of scale out of the nodes.
func ScaleOut(clusterCfg config.ClusterDetails, usrCfg config.UserConfig, t *time.Time) (bool, error) {
	// Read the current state of scaleup process and proceed with next step
	// If no stage was already set. The function returns an empty string. Then, start the scaleup process
	state.GetCurrentState()
	crypto.GetDecryptedCloudCreds(&clusterCfg.CloudCredentials)
	crypto.GetDecryptedOsCreds(&clusterCfg.OsCredentials)
	var newNodeIp string
	simFlag := usrCfg.MonitorWithSimulator
	monitorWithLogs := usrCfg.MonitorWithLogs
	isAccelerated := usrCfg.IsAccelerated

	switch state.CurrentState {
	case "provisioning_scaleup":
		log.Info.Println("Starting scaleUp process")
		if simFlag && isAccelerated {
			fakeSleep(t)
		}
		state.PreviousState = state.CurrentState
		state.CurrentState = "start_scaleup_process"
		state.ProvisionStartTime = time.Now().UnixMilli()
		state.UpdateState()
		fallthrough
		// Spin new VMs based on number of nodes and cloud type
	case "start_scaleup_process":
		if monitorWithLogs {
			log.Info.Println("Spin new vms based on the cloud type")
			time.Sleep(time.Duration(usrCfg.PollingInterval) * time.Second)
			if simFlag && isAccelerated {
				fakeSleep(t)
			}
			log.Info.Println("Spinning AWS instance")
			time.Sleep(time.Duration(usrCfg.PollingInterval) * time.Second)
			if simFlag && isAccelerated {
				fakeSleep(t)
			}
		} else {
			var err error
			newNodeIp, err = SpinNewVm(clusterCfg.LaunchTemplateId, clusterCfg.LaunchTemplateVersion, clusterCfg.CloudCredentials)
			if err != nil {
				return false, err
			}
		}
		log.Info.Println("Spinned a new node: ", newNodeIp)
		state.NodeIp = newNodeIp
		state.PreviousState = state.CurrentState
		state.CurrentState = "scaleup_triggered_spin_vm"
		state.UpdateState()
		fallthrough
	// Add the newly added VM to the list of VMs
	// Configure OS on newly created VM
	case "scaleup_triggered_spin_vm":
		state.GetCurrentState()
		newNodeIp = state.NodeIp
		if monitorWithLogs {
			log.Info.Println("Adding the spinned nodes into the list of vms")
			time.Sleep(time.Duration(usrCfg.PollingInterval) * time.Second)
			if simFlag && isAccelerated {
				fakeSleep(t)
			}
			log.Info.Println("Configure ES")
			time.Sleep(time.Duration(usrCfg.PollingInterval) * time.Second)
			if simFlag && isAccelerated {
				fakeSleep(t)
			}
			log.Info.Println("Configuring in progress")
			if simFlag && isAccelerated {
				fakeSleep(t)
			}
		} else {
			log.Info.Println("Configuring Opensearch on new node...")
			hostsFileName := "ansible_scripts/hosts"
			username := clusterCfg.SshUser
			f, err := os.OpenFile(hostsFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				log.Fatal.Println(err)
				return false, err
			}
			defer f.Close()
			nodes := utils.GetNodes()
			dataWriter := bufio.NewWriter(f)
			dataWriter.WriteString("[current_nodes]\n")
			for _, nodeIdMap := range nodes {
				_, writeErr := dataWriter.WriteString(nodeIdMap.(map[string]string)["name"] + " ansible_user=" + username + " roles=master,data,ingest ansible_private_host=" + nodeIdMap.(map[string]string)["hostIp"] + " ansible_ssh_private_key_file=" + clusterCfg.CloudCredentials.PemFilePath + "\n")
				if writeErr != nil {
					log.Error.Println("Error writing the node data into hosts file", writeErr)
				}
			}
			dataWriter.WriteString("[new_node]\n")
			dataWriter.WriteString("node-" + strings.ReplaceAll(newNodeIp, ".", "-") + " ansible_user=" + username + " roles=master,data,ingest ansible_private_host=" + newNodeIp + " ansible_ssh_private_key_file=" + clusterCfg.CloudCredentials.PemFilePath + "\n")
			dataWriter.Flush()
			ansibleErr := ansibleutils.CallAnsible(username, hostsFileName, clusterCfg, "scale_up")
			if ansibleErr != nil {
				if newNodeIp != "" {
					log.Warn.Println("Terminating the instance as the ansible script failed.")
					terminateErr := TerminateInstance(newNodeIp, clusterCfg.CloudCredentials)
					if terminateErr != nil {
						log.Fatal.Println(terminateErr)
					}
				}
				return false, ansibleErr
			}
		}
		state.PreviousState = state.CurrentState
		state.CurrentState = "provisioning_scaleup_configured"
		state.UpdateState()
		fallthrough
	case "provisioning_scaleup_configured":
		state.GetCurrentState()
		newNodeIp = state.NodeIp
		// Check if node has joined the cluster
		log.Info.Println("Waiting for new node to join the cluster...")
		time.Sleep(40 * time.Second)
		nodesInfo := utils.GetNodes()
		var joined bool
		for _, nodeIdInfo := range nodesInfo {
			if nodeIdInfo.(map[string]string)["hostIp"] == newNodeIp {
				joined = true
				break
			}
		}
		if !joined {
			errMsg := "The new node doesn't seem to have joined the cluster. Please login into new node and check for opensearch logs for more details."
			return joined, errors.New(errMsg)
		}

		// Install and start scaling manager on new node
		hostsFileName := "ansible_scripts/install_hosts"
		f, err := os.OpenFile(hostsFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal.Println(err)
			return false, err
		}
		defer f.Close()
		dataWriter := bufio.NewWriter(f)
		dataWriter.WriteString("[new_node]\n")
		dataWriter.WriteString("node-" + strings.ReplaceAll(newNodeIp, ".", "-") + " ansible_user=" + clusterCfg.SshUser + " roles=master,data,ingest ansible_private_host=" + newNodeIp + " ansible_ssh_private_key_file=" + clusterCfg.CloudCredentials.PemFilePath + "\n")
		dataWriter.Flush()

		ansibleErr := ansibleutils.UpdateWithTags(clusterCfg.SshUser, hostsFileName, []string{"install", "update_config", "update_pem", "update_secret", "start"})
		if ansibleErr != nil {
			log.Error.Println(ansibleErr)
			log.Error.Println("Node scaled up but unable to run scaling manager on new node. Please check ansible logs for more details. (logs/playbook.log)")
		}
		state.PreviousState = state.CurrentState
		state.CurrentState = "provisioning_scaleup_completed"
		state.UpdateState()
		fallthrough
	// Check cluster status after the configuration
	case "provisioning_scaleup_completed":
		if simFlag {
			SimulateSharRebalancing("scaleOut", state.NumNodes, isAccelerated)
		}
		log.Info.Println("Waiting for the cluster to become healthy")
		if simFlag && isAccelerated {
			fakeSleep(t)
		}
		CheckClusterHealth(usrCfg, t)
	}
	return true, nil
}

// Input:
//
//	clusterCfg (config.ClusterDetails): Cluster Level config details
//	usrCfg (config.UserConfig): User defined config for application behavior
//
// Description:
//
//	ScaleIn will scale in the cluster with the number of nodes.
//	This function will invoke commands to remove a node from opensearch cluster.
//
// Return:
//
//	(bool): Return the status of scale in of the nodes.
func ScaleIn(clusterCfg config.ClusterDetails, usrCfg config.UserConfig, t *time.Time) (bool, error) {
	// Read the current state of scaledown process and proceed with next step
	// If no stage was already set. The function returns an empty string. Then, start the scaledown process
	crypto.GetDecryptedCloudCreds(&clusterCfg.CloudCredentials)
	crypto.GetDecryptedOsCreds(&clusterCfg.OsCredentials)
	state.GetCurrentState()
	var removeNodeIp, removeNodeName string
	var nodes map[string]interface{}
	monitorWithLogs := usrCfg.MonitorWithLogs
	simFlag := usrCfg.MonitorWithSimulator
	isAccelerated := usrCfg.IsAccelerated
	if state.CurrentState == "provisioning_scaledown" {
		log.Info.Println("Staring scaleDown process")
		state.PreviousState = state.CurrentState
		state.CurrentState = "start_scaledown_process"
		state.ProvisionStartTime = time.Now().UnixMilli()
		state.UpdateState()
	}
	// Identify the node which can be removed from the cluster.
	switch state.CurrentState {
	case "start_scaledown_process":
		log.Info.Println("Identify the node to remove from the cluster and store the node_ip")
		if monitorWithLogs {
			time.Sleep(time.Duration(usrCfg.PollingInterval) * time.Second)
			if simFlag && isAccelerated {
				fakeSleep(t)
			}
		} else {
			nodes = utils.GetNodes()
			for nodeId, nodeIdInfo := range nodes {
				if !(utils.CheckIfMaster(context.Background(), nodeId)) {
					removeNodeIp = nodeIdInfo.(map[string]string)["hostIp"]
					removeNodeName = nodeIdInfo.(map[string]string)["name"]
					break
				}
			}
		}
		state.NodeIp = removeNodeIp
		state.NodeName = removeNodeName
		log.Info.Println("Node identified for removal: ", removeNodeName, removeNodeIp)
		state.PreviousState = state.CurrentState
		state.CurrentState = "scaledown_node_identified"
		state.UpdateState()
		fallthrough
	// Configure OS to tell master node that the present node is going to be removed
	case "scaledown_node_identified":
		state.GetCurrentState()
		removeNodeIp = state.NodeIp
		removeNodeName = state.NodeName
		if monitorWithLogs {
			log.Info.Println("Configure ES to remove the node ip from cluster")
			time.Sleep(time.Duration(usrCfg.PollingInterval) * time.Second)
			if simFlag && isAccelerated {
				fakeSleep(t)
			}
			log.Info.Println("Shutdown the node by ssh")
			time.Sleep(time.Duration(usrCfg.PollingInterval) * time.Second)
			if simFlag && isAccelerated {
				fakeSleep(t)
			}
		} else {
			log.Info.Println("Configuring to remove the node from cluster through ansible")
			hostsFileName := "ansible_scripts/hosts"
			username := clusterCfg.SshUser
			f, err := os.OpenFile(hostsFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				log.Error.Println(err)
				return false, err
			}
			defer f.Close()
			dataWriter := bufio.NewWriter(f)
			dataWriter.WriteString("[current_nodes]\n")
			for _, nodeIdInfo := range nodes {
				if nodeIdInfo.(map[string]string)["hostIp"] != removeNodeIp {
					_, writeErr := dataWriter.WriteString(nodeIdInfo.(map[string]string)["name"] + " ansible_user=" + username + " roles=master,data,ingest ansible_private_host=" + nodeIdInfo.(map[string]string)["hostIp"] + " ansible_ssh_private_key_file=" + clusterCfg.CloudCredentials.PemFilePath + "\n")
					if writeErr != nil {
						log.Error.Println("Error writing the node data into hosts file", writeErr)
					}
				}
			}
			dataWriter.WriteString("[remove_node]\n")
			dataWriter.WriteString(removeNodeName + " ansible_user=" + username + " roles=master,data,ingest ansible_private_host=" + removeNodeIp + " ansible_ssh_private_key_file=" + clusterCfg.CloudCredentials.PemFilePath + "\n")
			dataWriter.Flush()
			log.Info.Println("Removing node ***********************************:", removeNodeName)
			ansibleErr := ansibleutils.CallAnsible(username, hostsFileName, clusterCfg, "scale_down")
			if ansibleErr != nil {
				return false, ansibleErr
			}
		}
		state.PreviousState = state.CurrentState
		state.CurrentState = "provisioned_scaledown_on_cluster"
		state.UpdateState()
		fallthrough
	case "provisioned_scaledown_on_cluster":
		state.GetCurrentState()
		removeNodeIp = state.NodeIp
		log.Info.Println("Terminating the instance")
		terminateErr := TerminateInstance(removeNodeIp, clusterCfg.CloudCredentials)
		if terminateErr != nil {
			log.Fatal.Println(terminateErr)
			return false, terminateErr
		}
		state.PreviousState = state.CurrentState
		state.CurrentState = "provisioning_scaledown_completed"
		state.UpdateState()
		fallthrough
	// Wait for cluster to be in stable state(Shard rebalance)
	// Shut down the node
	case "provisioning_scaledown_completed":
		if simFlag {
			SimulateSharRebalancing("scaleIn", state.NumNodes, isAccelerated)
		}
		log.Info.Println("Wait for the cluster to become healthy and then proceed")
		CheckClusterHealth(usrCfg, t)
		if simFlag && isAccelerated {
			fakeSleep(t)
		}
	}
	return true, nil
}

// Input:
//
//	usrCfg (config.UserConfig): User defined config for application behavior
//
// Description:
//
//	CheckClusterHealth will check the current cluster health and also check if there are any relocating
//	shards. If the cluster status is green and there are no relocating shard then we will update the status
//	to provisioned_successfully. Else, we will wait for 3 minutes and perform this check again for 3 times.
//
// Return:
func CheckClusterHealth(usrCfg config.UserConfig, t *time.Time) {
	var timedOut bool
	simFlag := usrCfg.MonitorWithSimulator
	isAccelerated := usrCfg.IsAccelerated
	state.GetCurrentState()
	clusterDynamic, _ := cluster.GetClusterCurrent(false)
	if clusterDynamic.NumUnassignedShards > 0 {
		log.Info.Println("Retrying to reroute unassigned shards once before waiting for rebalancing")
		_, err := osutils.RerouteRetryFailed(context.Background())
		if err != nil {
			log.Error.Println("Failed to retry reroute", err)
		}
	}
	for {
		if simFlag {
			_ = cluster_sim.GetClusterCurrent(isAccelerated)
		} else {
			_, timedOut = cluster.GetClusterCurrent(true)
		}
		if !timedOut {
			state.PreviousState = state.CurrentState
			if strings.Contains(state.PreviousState, "scaleup") {
				state.CurrentState = "provisioned_scaleup_successfully"
			} else {
				state.CurrentState = "provisioned_scaledown_successfully"
			}
			state.UpdateState()
			break
		} else {
			log.Info.Println("Waiting for cluster to rebalance.......")
			time.Sleep(time.Duration(usrCfg.PollingInterval) * time.Second)
			if simFlag && isAccelerated {
				fakeSleep(t)
			}
		}
	}
}

// Inputs:
//
//	operation (string): An input required by the simulator to know the operation being pperformed (scaleup/scaledown)
//	numNode (int): The number of nodes added/removed from the cluster to simulate
//
// Description:
//
//	Calls the simulator api with the details of nodes added/removed to simulate the shard rebalancing operation
//
// Return:
func SimulateSharRebalancing(operation string, numNode int, isAccelerated bool) {
	// Add logic to call the simulator's end point
	var byteStr string
	if isAccelerated {
		t_now := time.Now()
		time_now := fmt.Sprintf("%02d:%02d:%02d", t_now.Hour(), t_now.Minute(), t_now.Second())
		date_now := fmt.Sprintf("%02d-%02d-%d", t_now.Day(), t_now.Month(), t_now.Year())
		byteStr = fmt.Sprintf("{\"nodes\":%d, \"time_now\":\"%s%s%s\"}", numNode, date_now, " ", time_now)
	} else {
		byteStr = fmt.Sprintf("{\"nodes\":%d}", numNode)
	}
	var jsonStr = []byte(byteStr)
	log.Debug.Println(string(jsonStr))
	var urlLink string
	if operation == "scaleOut" {
		urlLink = fmt.Sprintf("http://localhost:5000/provision/addnode")
	} else {
		urlLink = fmt.Sprintf("http://localhost:5000/provision/remnode")
	}

	req, err := http.NewRequest("POST", urlLink, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Panic.Println(err)
		panic(err)
	}

	if resp.StatusCode != 200 {
		log.Error.Println(resp.Status)
	}

	defer resp.Body.Close()
}

// Inputs:
//
//	t *time.Time
//
// Description:
//
//	Accelerate the sleep to the time duration mentioned in the input.
//
// Return:
func fakeSleep(t *time.Time) {
	*t = t.Add(time.Minute * 5)
	f := faketime.NewFaketimeWithTime(*t)
	f.Do()
	log.Info.Println(time.Now())
}

// Inputs:
//
// Description:
//
//	Sets the CurrentState to normal, updates the other fields with default and updates the opensearch document with the same
//
// Return:
func SetStateBackToNormal() {
	state.LastProvisionedTime = time.Now().UnixMilli()
	state.ProvisionStartTime = 0
	state.PreviousState = state.CurrentState
	state.CurrentState = "normal"
	state.RuleTriggered = ""
	state.RemainingNodes = 0
	state.UpdateState()
	log.Info.Println("State set back to normal")
}

// Inputs:
//
//	status (string): Status of the Provisioning
//	err (error): Error if any during provisioning
//
// Description:
//
//	Adds a document to Opensearch representing the status of the provisioning that took place
//
// Return:
func PushToOs(status string, err error) {
	provisionState := make(map[string]interface{}, 0)
	provisionState["RuleTriggered"] = state.RuleTriggered
	provisionState["ProvisionStartTime"] = state.ProvisionStartTime
	provisionState["ProvisionEndTime"] = time.Now().UnixMilli()
	provisionState["NumNodes"] = state.NumNodes
	provisionState["Status"] = status
	if err != nil {
		provisionState["FailureReason"] = err.Error()
	}
	provisionState["RulesResponsible"] = state.RulesResponsible
	provisionState["TimeTaken"] = fmt.Sprint((time.UnixMilli(provisionState["ProvisionEndTime"].(int64))).Sub(time.UnixMilli(provisionState["ProvisionStartTime"].(int64))))
	provisionState["StatTag"] = "ProvisionStats"
	provisionState["_documentType"] = "ProvisionStats"
	provisionState["Timestamp"] = time.Now().UnixMilli()

	doc, err := json.Marshal(provisionState)
	if err != nil {
		log.Panic.Println("json.Marshal ERROR: ", err)
		panic(err)
	}

	indexResponse, err := osutils.IndexMetrics(context.Background(), doc)
	if err != nil {
		log.Panic.Println("Failed to insert provision stats document: ", err)
		panic(err)
	}
	defer indexResponse.Body.Close()
	log.Debug.Println("Update resp: ", indexResponse)
}
