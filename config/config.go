package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"scaling_manager/cluster"
	"scaling_manager/task"

	"gopkg.in/yaml.v3"
)

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

// This struct contains the data structure to parse the cluster details present in the configuration file.
type ClusterDetails struct {
	// ClusterStatic indicates the static configuration for the cluster.
	cluster.ClusterStatic `yaml:",inline"`
	OsCredentials         OsCredentials    `yaml:"os_credentials"`
	CloudCredentials      CloudCredentials `yaml:"cloud_credentials"`
}

// This struct contains the data structure to parse the configuration file.
type ConfigStruct struct {
	ClusterDetails ClusterDetails `yaml:"cluster_details"`
	TaskDetails    []task.Task    `yaml:"task_details"`
}

// Inputs:
//
//	path(string): The path of the configuration file.
//
// Description:
//
//	This function will be parsing the provided configuration file and populate the ConfigStruct.
//
// Return:
//
//	Return the ConfigStruct.
func GetConfig(path string) ConfigStruct {
	yamlConfig, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to read the config file: ", err)
	}
	fmt.Println("File opened!")
	defer yamlConfig.Close()
	configByte, _ := ioutil.ReadAll(yamlConfig)
	var config = new(ConfigStruct)
	err = yaml.Unmarshal(configByte, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return *config
}
