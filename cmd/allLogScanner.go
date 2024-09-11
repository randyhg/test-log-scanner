package cmd

import (
	"github.com/randyhg/test-log-scanner/config"
	"github.com/randyhg/test-log-scanner/service"
	"github.com/randyhg/test-log-scanner/util"
	"github.com/randyhg/test-log-scanner/util/mylog"
	"github.com/spf13/cobra"
	"time"
)

var allLogScannerCmd = &cobra.Command{
	Use:   "all-log-scanner",
	Short: "Scan All Log",
	Run:   startScanAll,
}

func init() {
	rootCmd.AddCommand(allLogScannerCmd)
}

func startScanAll(cmd *cobra.Command, args []string) {
	config.Init()
	util.InitDB()
	util.InitRedis()

	mylog.Info("start scanning")
	for {
		if err := service.ScanGzFiles("some-gz-log-url"); err != nil {
			mylog.Error(err)
		}
		time.Sleep(10 * time.Second)
		continue
	}
}
