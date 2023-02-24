package cmd

import (
	crypto "github.com/maplelabs/opensearch-scaling-manager/crypto"
	app "github.com/maplelabs/opensearch-scaling-manager/scaleManager"
	"github.com/spf13/cobra"
)

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
	crypto.Initialize("fetchMetrics")
	app.Initialize()
	app.StartFetchMetrics()
}
