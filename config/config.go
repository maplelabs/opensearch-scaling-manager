package config

import (
	"io/ioutil"
	"os"
	"regexp"
	"scaling_manager/cluster"
	"scaling_manager/crypto"
	"scaling_manager/logger"
	"scaling_manager/recommendation"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var log logger.LOG

// Input:
//
// Description:
//
//	Initialize the logger module.
//
// Return:
func init() {
	log.Init("logger")
	log.Info.Println("Config module initialized")
}

// This struct contains the OS Admin Username and OS Admin Password via which we can connect to OS cluster.
type OsCredentials struct {
	// OsAdminUsername indicates the OS Admin Username via which OS client can connect to OS Cluster.
	OsAdminUsername string `yaml:"os_admin_username" validate:"required"`
	// OsAdminPassword indicates the OS Admin Password via which OS client can connect to OS Cluster.
	OsAdminPassword string `yaml:"os_admin_password" validate:"required"`
}

// This struct contains the Cloud Secret Key and Access Key via which we can connect to the cloud.
type CloudCredentials struct {
	// SecretKey indicates the Secret key for connecting to the cloud.
	SecretKey string `yaml:"secret_key" validate:"required"`
	// AccessKey indicates the Access key for connecting to the cloud.
	AccessKey string `yaml:"access_key" validate:"required"`
}

// This struct contains the data structure to parse the cluster details present in the configuration file.
type ClusterDetails struct {
	// ClusterStatic indicates the static configuration for the cluster.
	cluster.ClusterStatic `yaml:",inline"`
	OsCredentials         OsCredentials    `yaml:"os_credentials"`
	CloudCredentials      CloudCredentials `yaml:"cloud_credentials"`
}

// Config for application behaviour from user
type UserConfig struct {
	MonitorWithLogs      bool `yaml:"monitor_with_logs"`
	MonitorWithSimulator bool `yaml:"monitor_with_simulator"`
	PollingInterval      int  `yaml:"polling_interval_in_secs"`
}

// This struct contains the data structure to parse the configuration file.
type ConfigStruct struct {
	UserConfig     UserConfig            `yaml:"user_config"`
	ClusterDetails ClusterDetails        `yaml:"cluster_details"`
	TaskDetails    []recommendation.Task `yaml:"task_details" validate:"gt=0,dive"`
}

// Inputs:
//
//		path (string): The path of the configuration file.
//	 firstExecution (bool) : true if the function is called for the first time, else false
//	 fromFileHandler (bool) : true if the function is called from file event handler
//	 function, else false
//
// Description:
//
//	This function will be parsing the provided configuration file and populate the ConfigStruct.
//
// Return:
//
// ConfigStruct : Structure of the config.yaml file.
// Error : Error (if any), else nil
func GetConfig(path string, firstExecution bool, fromFileHandler bool) (ConfigStruct, error) {
	var config = new(ConfigStruct)
	yamlConfig, err := os.Open(path)
	if err != nil {
		log.Panic.Println("Unable to read the config file: ", err)
		return *config, err
	}
	defer yamlConfig.Close()
	configByte, _ := ioutil.ReadAll(yamlConfig)

	err = yaml.Unmarshal(configByte, &config)
	if err != nil {
		log.Panic.Println("ConfigStruct unmarshal error : ", err)
		return *config, err
	}
	err = validation(*config)
	if err != nil {
		log.Panic.Println("ConfigStruct validation error : ", err)
		return *config, err
	}

	if !firstExecution && !fromFileHandler {
		decryptedOsUsername, err := crypto.GetDecryptedData(config.ClusterDetails.OsCredentials.OsAdminUsername)
		if decryptedOsUsername != "" {
			config.ClusterDetails.OsCredentials.OsAdminUsername = decryptedOsUsername
		} else {
			return *config, err
		}

		decryptedOsPassword, err := crypto.GetDecryptedData(config.ClusterDetails.OsCredentials.OsAdminPassword)
		if decryptedOsPassword != "" {
			config.ClusterDetails.OsCredentials.OsAdminPassword = decryptedOsPassword
		} else {
			return *config, err
		}

		decryptedCloudSecretkey, err := crypto.GetDecryptedData(config.ClusterDetails.CloudCredentials.SecretKey)
		if decryptedCloudSecretkey != "" {
			config.ClusterDetails.CloudCredentials.SecretKey = decryptedCloudSecretkey
		} else {
			return *config, err
		}

		decryptedCloudAccesskey, err := crypto.GetDecryptedData(config.ClusterDetails.CloudCredentials.AccessKey)
		if decryptedCloudAccesskey != "" {
			config.ClusterDetails.CloudCredentials.AccessKey = decryptedCloudAccesskey
		} else {
			return *config, err
		}
	}

	return *config, err
}

// Inputs:
//
//	config(ConfigStruct): config structure populated with unmarshalled data.
//
// Description:
//
//	This function will be validating the configuration structure.
//
// Return:
//
//	Return the error if there is a validation error.
func validation(config ConfigStruct) error {
	validate := validator.New()
	validate.RegisterValidation("isValidName", isValidName)
	validate.RegisterValidation("isValidTaskName", isValidTaskName)
	err := validate.Struct(config)
	return err
}

// Inputs:
//
//	field(validator.FieldLevel): The field which needs to be validated.
//
// Description:
//
//	This function will be validating the cluster name.
//
// Return:
//
//	Return true if there is a valida cluster name else false.
func isValidName(fl validator.FieldLevel) bool {
	nameRegexString := `^[a-zA-Z][a-zA-Z0-9\-\._]+[a-zA-Z0-9]$`
	nameRegex := regexp.MustCompile(nameRegexString)

	return nameRegex.MatchString(fl.Field().String())
}

// Inputs:
//
//	field(validator.FieldLevel): The field which needs to be validated.
//
// Description:
//
//	This function will be validating the Task name.
//
// Return:
//
//	Return true if there is a valid Task name else false.
func isValidTaskName(fl validator.FieldLevel) bool {
	TaskNameRegexString := `scale_(up|down)_by_[0-9]+`
	TaskNameRegex := regexp.MustCompile(TaskNameRegexString)

	return TaskNameRegex.MatchString(fl.Field().String())
}

// Inputs:
//
// filePath (string): Path if the config file.
// initialRun (bool): true if the function is called for the first time, else false
// updatedCredsArray ([]string): array which stores the previous creds
//
// Description:
//
//	This function updates the config file with the encrypted creds on the first run.
//	 During the second run, the function checks if there was any updates made to the
//	 credentials in the config.yaml file, and updates only if the changes are made to
//	 the credentials
//
// Return:
//
// ConfigStruct : Structure of the config.yaml file.
// []string : Updated credentials, if there was any changes in the config.yaml file
// error : Error (if any), else nil
func UpdateEncryptedCred(filePath string, initialRun bool, updatedCredsArray []string) (ConfigStruct, []string, error) {
	update_flag := false
	configStruct, err := GetConfig(filePath, initialRun, true)
	if err != nil {
		log.Panic.Println("The recommendation can not be made as there is an error in the validation of config file.", err)
		return configStruct, updatedCredsArray, err
	}
	unencryptedConfigStruct := configStruct
	//log.Info.Println(unencryptedConfigStruct)

	if initialRun {
		update_flag = true

		configStruct.ClusterDetails.OsCredentials.OsAdminUsername, err = crypto.GetEncryptedData(configStruct.ClusterDetails.OsCredentials.OsAdminUsername)
		updatedCredsArray[0] = configStruct.ClusterDetails.OsCredentials.OsAdminUsername
		if err != nil {
			return unencryptedConfigStruct, updatedCredsArray, err
		}

		configStruct.ClusterDetails.OsCredentials.OsAdminPassword, err = crypto.GetEncryptedData(configStruct.ClusterDetails.OsCredentials.OsAdminPassword)
		updatedCredsArray[1] = configStruct.ClusterDetails.OsCredentials.OsAdminPassword
		if err != nil {
			return unencryptedConfigStruct, updatedCredsArray, err
		}

		configStruct.ClusterDetails.CloudCredentials.SecretKey, err = crypto.GetEncryptedData(configStruct.ClusterDetails.CloudCredentials.SecretKey)
		updatedCredsArray[2] = configStruct.ClusterDetails.CloudCredentials.SecretKey
		if err != nil {
			return unencryptedConfigStruct, updatedCredsArray, err
		}

		configStruct.ClusterDetails.CloudCredentials.AccessKey, err = crypto.GetEncryptedData(configStruct.ClusterDetails.CloudCredentials.AccessKey)
		updatedCredsArray[3] = configStruct.ClusterDetails.CloudCredentials.AccessKey
		if err != nil {
			return unencryptedConfigStruct, updatedCredsArray, err
		}
	} else {
		if configStruct.ClusterDetails.OsCredentials.OsAdminUsername != updatedCredsArray[0] {
			update_flag = true
			configStruct.ClusterDetails.OsCredentials.OsAdminUsername, err = crypto.GetEncryptedData(configStruct.ClusterDetails.OsCredentials.OsAdminUsername)
			updatedCredsArray[0] = configStruct.ClusterDetails.OsCredentials.OsAdminUsername
			if err != nil {
				return unencryptedConfigStruct, updatedCredsArray, err
			}
		}

		if configStruct.ClusterDetails.OsCredentials.OsAdminPassword != updatedCredsArray[1] {
			update_flag = true
			configStruct.ClusterDetails.OsCredentials.OsAdminPassword, err = crypto.GetEncryptedData(configStruct.ClusterDetails.OsCredentials.OsAdminPassword)
			updatedCredsArray[1] = configStruct.ClusterDetails.OsCredentials.OsAdminPassword
			if err != nil {
				return unencryptedConfigStruct, updatedCredsArray, err
			}
		}

		if configStruct.ClusterDetails.CloudCredentials.SecretKey != updatedCredsArray[2] {
			update_flag = true
			configStruct.ClusterDetails.CloudCredentials.SecretKey, err = crypto.GetEncryptedData(configStruct.ClusterDetails.CloudCredentials.SecretKey)
			updatedCredsArray[2] = configStruct.ClusterDetails.CloudCredentials.SecretKey
			if err != nil {
				return unencryptedConfigStruct, updatedCredsArray, err
			}
		}

		if configStruct.ClusterDetails.CloudCredentials.AccessKey != updatedCredsArray[3] {
			update_flag = true
			configStruct.ClusterDetails.CloudCredentials.AccessKey, err = crypto.GetEncryptedData(configStruct.ClusterDetails.CloudCredentials.AccessKey)
			updatedCredsArray[3] = configStruct.ClusterDetails.CloudCredentials.AccessKey
			if err != nil {
				return unencryptedConfigStruct, updatedCredsArray, err
			}
		}
	}
	time.Sleep(3 * time.Second)
	if update_flag {
		err = UpdateConfigFile(configStruct)
		if err != nil {
			return unencryptedConfigStruct, updatedCredsArray, err
		}
	} else {
		log.Info.Println("Credentials not updated, hence config file update not required")
	}

	return unencryptedConfigStruct, updatedCredsArray, nil
}

// Inputs:
//
// ConfigStruct : Credentials encrypted structure of the config.yaml file
//
// Description:
//
//	This function updates the config.yaml file with encrypted credentials ConfigStruct.
//
// Return:
//
// Error : Error (if any), else nil
func UpdateConfigFile(conf ConfigStruct) error {
	conf_byte, err := yaml.Marshal(&conf)
	if err != nil {
		log.Error.Println("Error marshalling the ConfigStruct : ", err)
		return err
	}

	yaml_content := "---\n" + string(conf_byte)
	err = ioutil.WriteFile("config.yaml", []byte(yaml_content), 0)
	if err != nil {
		log.Error.Println("Error writing the config yaml file : ", err)
		return err
	}

	return nil
}
