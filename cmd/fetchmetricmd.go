package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/maplelabs/opensearch-scaling-manager/config"
	fetch "github.com/maplelabs/opensearch-scaling-manager/fetchmetrics"
	"github.com/maplelabs/opensearch-scaling-manager/logger"
	osutils "github.com/maplelabs/opensearch-scaling-manager/opensearchUtils"
	"github.com/maplelabs/opensearch-scaling-manager/provision"
)

// Logger variable used across the package for logging.
var log logger.LOG

// Start Command to start the Scaling Manager service
var fetchMetricStartCmd = &cobra.Command{
	Use:   "fetchmetric start",
	Short: "Start fetchmetric",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fetchMetricstart()
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
func fetchMetricstart() {
	configStruct, err := config.GetConfig("config.yaml")
	if err != nil {
		log.Panic.Println("The recommendation can not be made as there is an error in the validation of config file.", err)
		panic(err)
	}
	cfg := configStruct.ClusterDetails

	osutils.InitializeOsClient(cfg.OsCredentials.OsAdminUsername, cfg.OsCredentials.OsAdminPassword)

	provision.InitializeDocId()

	userCfg := configStruct.UserConfig

	if !userCfg.MonitorWithSimulator {
		fetch.FetchMetrics(userCfg.PollingInterval, userCfg.PurgeAfter)
	} else {
		log.Warn.Println("MonitorWithSimulator is enabled. Please disable and re-run the fetch-metrics module.")
		os.Exit(1)
	}
}
