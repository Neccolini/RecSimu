/*
Copyright © 2023 Neccolini <shun11202991@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/Neccolini/RecSimu/cmd/debug"
	"github.com/Neccolini/RecSimu/cmd/read"
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
		filepath, err := cmd.Flags().GetString("input")
		if err != nil {
			log.Fatal(err)
		}
		debugBoolean, err := cmd.Flags().GetBool("debug")
		if err != nil {
			log.Fatal(err)
		}

		input := read.ReadJsonFile(filepath)
		config := run.NewSimulationConfig(input.NodeNum, input.Cycle, input.AdjacencyList, input.NodesType, input.ReconfigureInfo)
		debug.Debug.On = debugBoolean

		config.Simulate("test")

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("input", "i", "", "input configuration file")
	runCmd.Flags().StringP("output", "o", "", "output file")
	runCmd.Flags().BoolP("debug", "d", false, "debug mode")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
