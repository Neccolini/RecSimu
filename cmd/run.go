/*
Copyright © 2023 Neccolini <shun11202991@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/Neccolini/RecSimu/cmd/run"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run simulation",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatal(err)
		}
		m := map[int][]int{0: {1}, 1: {0}}
		l := map[int]string{0: "Coordinator", 1: "Router"}
		config := run.NewSimulationConfig(2, 2, m, l)

		config.Simulate("test")

		fmt.Println("OK")

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("config", "c", "", "configuration file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
