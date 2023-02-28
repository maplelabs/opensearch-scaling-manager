package cmd

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/maplelabs/opensearch-scaling-manager/config"
	crypto "github.com/maplelabs/opensearch-scaling-manager/crypto"
	app "github.com/maplelabs/opensearch-scaling-manager/scaleManager"
	"github.com/spf13/cobra"
)

// Start Command to start the Scaling Manager service
var cryptoStartCmd = &cobra.Command{
	Use:   "crypto",
	Short: "Handles crypto module",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		start, _ := cmd.Flags().GetString("start")
		stop_pid, _ := cmd.Flags().GetString("stop")
		decrypt_credentials, _ := cmd.Flags().GetString("decrypt")
		if start != "" {
			CryptoStart()
		} else if stop_pid != "" {
			CryptoStop(stop_pid)
		} else if decrypt_credentials != "" {
			decryptedCredentials := CryptoDecrypt(decrypt_credentials)
			fmt.Println(decryptedCredentials)
		} else if start != "" && stop_pid != "" {
			log.Panic.Println("Please provide either start or stop command.")
		}
	},
}

// Input:
//
// Description:
//
//	Initializes the start command, adds the required flags
//
// Return:
func init() {
	log.Init("logger")
	cryptoStartCmd.PersistentFlags().String("start", "", "To start crypto")
	cryptoStartCmd.PersistentFlags().String("stop", "", "To stop crypto")
	cryptoStartCmd.PersistentFlags().String("decrypt", "", "To decrypt the credentials. Please Pass os_credentials to decrypt os credentials or cloud_credentials to decrypt cloud credentials")
}

// Input:
//
// Description:
//
//	The Function initilazes and starts the execution of Scaling Manager
//
// Return:
//
//	(error): Returns error upon unsuccessful execution
func CryptoStart() {
	crypto.Initialize()
	configStruct, err := config.GetConfig()
	if err != nil {
		log.Panic.Println("The recommendation can not be made as there is an error in the validation of config file.", err)
		panic(err)
	}
	app.FileWatch(configStruct, "crypto")
}

// Input:
//
// Description:
//
//		Function reads the Process Id file and stops the running instance
//	 of Scaling Manager.
//
// Return:
//
// (error): Returns error upon unsuccessful execution.
func CryptoStop(pid string) error {
	log.Info.Println("Stopping fetch metrics")
	var pid_int int
	var err error

	pid_int, err = strconv.Atoi(string(pid))
	proc, err := os.FindProcess(pid_int)
	if err != nil {
		log.Error.Println("Process not found ", err)
		return err
	}

	err = proc.Signal(os.Interrupt)
	if err != nil {
		log.Error.Println("Unable to terminate process ", err)
		return err
	}

	time.Sleep(2 * time.Second)

	proc, err = os.FindProcess(pid_int)
	if err != nil {
		log.Info.Println("Process Terminate Successful")
		return nil
	}

	err = proc.Signal(os.Signal(syscall.Signal(0)))
	if err == nil {
		log.Info.Printf("Process with Pid %v is still running.", pid_int)
		log.Info.Println("Scale Manager currently in the provision phase and will be shut down once it is completed")
	} else {
		log.Info.Printf("Process with pid %v is not running.", pid_int)
		log.Info.Println("Process Terminate Successful")
	}
	return nil
}

func CryptoDecrypt(credentials string) string {
	var decryptedCredentials string

	configStruct, err := config.GetConfig()
	if err != nil {
		log.Panic.Println("The recommendation can not be made as there is an error in the validation of config file.", err)
		panic(err)
	}
	if credentials == "os_username" {
		crypto.GetDecryptedOsCreds(&configStruct.ClusterDetails.OsCredentials)
		decryptedCredentials = configStruct.ClusterDetails.OsCredentials.OsAdminUsername
	} else if credentials == "os_password" {
		crypto.GetDecryptedOsCreds(&configStruct.ClusterDetails.OsCredentials)
		decryptedCredentials = configStruct.ClusterDetails.OsCredentials.OsAdminPassword
	} else if credentials == "cloud_secret_key" {
		crypto.GetDecryptedCloudCreds(&configStruct.ClusterDetails.CloudCredentials)
		decryptedCredentials = configStruct.ClusterDetails.CloudCredentials.SecretKey
	} else if credentials == "cloud_access_key" {
		crypto.GetDecryptedCloudCreds(&configStruct.ClusterDetails.CloudCredentials)
		decryptedCredentials = configStruct.ClusterDetails.CloudCredentials.AccessKey
	} else {
		log.Panic.Println("Please pass correct arguement.")
		log.Panic.Println("Please Pass os_username, os_password, cloud_secret_key or cloud_access_key to decrypt the credentials")
	}
	return decryptedCredentials
}
