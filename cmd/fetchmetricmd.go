package cmd

import (
	"os"
	"strconv"

	"github.com/maplelabs/opensearch-scaling-manager/config"
	app "github.com/maplelabs/opensearch-scaling-manager/scaleManager"
	"github.com/spf13/cobra"
)

// Start Command to start the Scaling Manager service
var fetchMetricStartCmd = &cobra.Command{
	Use:   "fetchmetric",
	Short: "Handle fetchmetric module",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		start, _ := cmd.Flags().GetString("start")
		stop_pid, _ := cmd.Flags().GetString("stop")

		if start != "" {
			fetchMetricStart()
		} else if stop_pid != "" {
			fetchMetricStop(stop_pid)
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
	fetchMetricStartCmd.PersistentFlags().String("start", "", "To start fetchmetrics")
	fetchMetricStartCmd.PersistentFlags().String("stop", "", "To stop fetchmetrics")
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
func fetchMetricStart() {
	app.Initialize()
	app.StartFetchMetrics()
	configStruct, err := config.GetConfig()
	if err != nil {
		log.Panic.Println("The recommendation can not be made as there is an error in the validation of config file.", err)
		panic(err)
	}
	app.FileWatch(configStruct, "fetchMetrics")
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
func fetchMetricStop(pid string) error {
	log.Info.Println("Stopping fetch metrics")
	var pid_int int
	var err error

	pid_int, err = strconv.Atoi(string(pid))
	proc, err := os.FindProcess(pid_int)
	if err != nil {
		log.Error.Println("Process not found ", err)
		return err
	}

	err = proc.Signal(os.Kill)
	if err != nil {
		log.Error.Println("Unable to kill process ", err)
		return err
	}
	return nil
}
