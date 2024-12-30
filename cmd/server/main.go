package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "server",
	Long: "Api Rest School",
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("failed to execute command. err: %v", err)
		os.Exit(1)
	}
}
