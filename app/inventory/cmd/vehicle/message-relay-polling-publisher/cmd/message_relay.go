package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers []string
	GroupId string
	Topics  []string
}

type DbConfig struct {
	Conn DbConn
	Name string
}

type DbConn struct {
	Host     string
	Port     int
	User     string
	Password string
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	viper.AutomaticEnv()

	rootCmd.AddCommand(runCmd)
}
