package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"vehicle-sharing-go/pkg/cmd"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "inventory-message-relay-binlog",
	Short: "Poll events from Mysql outbox table and write them on Kafka preserving order.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		cmd.InitConfig(cmd.NewConfigYaml("./app/inventory/cmd/vehicle/message-relay-polling-publisher/cmd", "message-relay"))
	})
}
