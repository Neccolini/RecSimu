/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/Neccolini/RecSimu/cmd/gen"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		nodeNum, err := cmd.Flags().GetInt("nodenum")
		if err != nil {
			log.Fatal(err)
		}

		cycles, err := cmd.Flags().GetInt("cycles")
		if err != nil {
			log.Fatal(err)
		}

		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal(err)
		}

		if err := gen.GenerateNetwork(filePath, nodeNum, cycles); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().IntP("nodenum", "n", 0, "the number of nodes")
	genCmd.Flags().IntP("cycles", "c", 0, "cycles")
	genCmd.Flags().StringP("file", "f", "", "filepath")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
