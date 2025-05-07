package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"vehicle-sharing-go/pkg/cmd"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "inventory-message-relay-binlog",
	Short: "Read events about an outbox table from MySQL binlog, write them on Kafka preserving order. Persist state on Redis.",
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
		cmd.InitConfig(cmd.NewConfigYaml("./app/inventory/cmd/vehicle/message-relay-binlog/cmd", "message-relay"))
	})

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
