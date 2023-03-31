package config

import (
	"io/ioutil"
	"os"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/maplelabs/opensearch-scaling-manager/cluster"
	"github.com/maplelabs/opensearch-scaling-manager/logger"
	"github.com/maplelabs/opensearch-scaling-manager/recommendation"
	"gopkg.in/yaml.v3"
)

var log logger.LOG
var ConfigFileName = "config.yaml"

// Input:
//
// Description:
//
//	Initialize the Config module.
//
// Return:
func init() {
	log.Init("logger")
	log.Info.Println("Config module initialized")
}

// This struct contains the OS Admin Username and OS Admin Password via which we can connect to OS cluster.
type OsCredentials struct {
	// OsAdminUsername indicates the OS Admin Username via which OS client can connect to OS Cluster.
	OsAdminUsername string `yaml:"os_admin_username" validate:"required" json:"os_admin_username"`
	// OsAdminPassword indicates the OS Admin Password via which OS client can connect to OS Cluster.
	OsAdminPassword string `yaml:"os_admin_password" validate:"required" json:"os_admin_password"`
}

// This struct contains the Cloud Secret Key and Access Key via which we can connect to the cloud.
type CloudCredentials struct {
	PemFilePath string `yaml:"pem_file_path" validate:"required" json:"pem_file_path"`
	// SecretKey indicates the Secret key for connecting to the cloud.
	SecretKey string `yaml:"secret_key" validate:"required" json:"secret_key"`
	// AccessKey indicates the Access key for connecting to the cloud.
	AccessKey string `yaml:"access_key" validate:"required" json:"access_key"`
	Region    string `yaml:"region" validate:"required" json:"region"`
}

// This struct contains the data structure to parse the cluster details present in the configuration file.
type ClusterDetails struct {
	// ClusterStatic indicates the static configuration for the cluster.
	cluster.ClusterStatic `yaml:",inline"`
	LaunchTemplateId      string           `yaml:"launch_template_id" validate:"required" json:"launch_template_id"`
	LaunchTemplateVersion string           `yaml:"launch_template_version" validate:"required" json:"launch_template_version"`
	SshUser               string           `yaml:"os_user" validate:"required" json:"os_user"`
	OsGroup               string           `yaml:"os_group" validate:"required" json:"os_group"`
	OpensearchVersion     string           `yaml:"os_version" validate:"required" json:"os_version"`
	OpensearchHome        string           `yaml:"os_home" validate:"required" json:"os_home"`
	DomainName            string           `yaml:"domain_name" validate:"required" json:"domain_name"`
	OsCredentials         OsCredentials    `yaml:"os_credentials" json:"os_credentials"`
	CloudCredentials      CloudCredentials `yaml:"cloud_credentials" json:"cloud_credentials"`
}

// Config for application behaviour from user
type UserConfig struct {
	MonitorWithLogs      bool `yaml:"monitor_with_logs"`
	MonitorWithSimulator bool `yaml:"monitor_with_simulator"`
	PurgeAfter           int  `yaml:"purge_old_docs_after_hours" validate:"required"`
	PollingInterval      int  `yaml:"polling_interval_in_secs" validate:"required"`
	IsAccelerated        bool `yaml:"is_accelerated"`
}

// This struct contains the data structure to parse the configuration file.
type ConfigStruct struct {
	UserConfig     UserConfig            `yaml:"user_config"`
	ClusterDetails ClusterDetails        `yaml:"cluster_details"`
	TaskDetails    []recommendation.Task `yaml:"task_details" validate:"gt=0,dive"`
}

// Inputs:
//
//	path (string): The path of the configuration file.
//
// Description:
//
//	This function will be parsing the provided configuration file and populate the ConfigStruct.
//
// Return:
//
//	(ConfigStruct, error): Return the ConfigStruct and error if any
func GetConfig() (ConfigStruct, error) {
	yamlConfig, err := os.Open(ConfigFileName)
	if err != nil {
		log.Panic.Println("Unable to read the config file: ", err)
		panic(err)
	}
	defer yamlConfig.Close()
	configByte, _ := ioutil.ReadAll(yamlConfig)
	var config = new(ConfigStruct)
	err = yaml.Unmarshal(configByte, &config)
	if err != nil {
		log.Panic.Println("Unmarshal Error : ", err)
		panic(err)
	}
	err = validation(*config)
	return *config, err
}

// Inputs:
//      config (ConfigStruct): config structure populated with unmarshalled data.
//
// Description:
//      This function will be validating the configuration structure.
//
// Return:
//      (error): Return the error if there is a validation error.

func validation(config ConfigStruct) error {
	validate := validator.New()
	validate.RegisterValidation("isValidName", isValidName)
	validate.RegisterValidation("isValidTaskName", isValidTaskName)
	err := validate.Struct(config)
	return err
}

// Inputs:
//
//	fl (validator.FieldLevel): The field which needs to be validated.
//
// Description:
//
//	This function will be validating the cluster name.
//
// Return:
//
//	(bool): Return true if there is a valid cluster name else false.
func isValidName(fl validator.FieldLevel) bool {
	nameRegexString := `^[a-zA-Z][a-zA-Z0-9\-\._]+[a-zA-Z0-9]$`
	nameRegex := regexp.MustCompile(nameRegexString)

	return nameRegex.MatchString(fl.Field().String())
}

// Inputs:
//
//	fl (validator.FieldLevel): The field which needs to be validated.
//
// Description:
//
//	This function will be validating the Task name.
//
// Return:
//
//	(bool): Return true if there is a valid Task name else false.
func isValidTaskName(fl validator.FieldLevel) bool {
	TaskNameRegexString := `scale_(up|down)_by_[0-9]+`
	TaskNameRegex := regexp.MustCompile(TaskNameRegexString)

	return TaskNameRegex.MatchString(fl.Field().String())
}

// Inputs:
//
//	conf (ConfigStruct) : Credentials encrypted structure of the config.yaml file
//
// Description:
//
//	This function updates the config.yaml file with encrypted credentials ConfigStruct.
//
// Return:
//
//	(error) : Error (if any), else nil
func UpdateConfigFile(conf ConfigStruct) error {
	conf_byte, err := yaml.Marshal(&conf)
	if err != nil {
		log.Error.Println("Error marshalling the ConfigStruct : ", err)
		return err
	}

	yaml_content := "---\n" + string(conf_byte)
	err = ioutil.WriteFile(ConfigFileName, []byte(yaml_content), 0)
	if err != nil {
		log.Error.Println("Error writing the config yaml file : ", err)
		return err
	}

	return nil
}
