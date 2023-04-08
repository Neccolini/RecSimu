/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strconv"
	"strings"

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
		cycles, err := cmd.Flags().GetInt("cycles")
		if err != nil {
			log.Fatal(err)
		}
		rate, err := cmd.Flags().GetFloat64("rate")
		if err != nil {
			log.Fatal(err)
		}

		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal(err)
		}
		rawTopology, err := cmd.Flags().GetString("topology")
		topology := strings.Fields(rawTopology)

		if err != nil {
			log.Fatal(err)
		}
		switch topology[0] {
		case "random":
			{
				if len(topology) < 2 {
					log.Fatal("usage: -t random <the number of nudes>")
				}
				nodeNum, err := strconv.Atoi(topology[1])
				if err != nil {
					log.Fatal(err)
				}
				config := gen.NewConfig(topology[0], []int{nodeNum})
				if err := gen.GenerateNetwork(*config, filePath, cycles, rate); err != nil {
					log.Fatal(err)
				}
			}
		case "mesh":
			{
				if len(topology) < 3 {
					log.Fatal("usage: -t mesh <rows> <columns>")
				}
				rows, err := strconv.Atoi(topology[1])
				if err != nil {
					log.Fatal(err)
				}
				columns, err := strconv.Atoi(topology[2])
				if err != nil {
					log.Fatal(err)
				}
				config := gen.NewConfig(topology[0], []int{rows, columns})
				if err := gen.GenerateNetwork(*config, filePath, cycles, rate); err != nil {
					log.Fatal(err)
				}
			}
		default:
			{
				log.Fatalf("topology %s is not defined", topology[0])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.PersistentFlags().StringP("topology", "t", "random 100", "topology: random or mesh")
	genCmd.Flags().IntP("cycles", "c", 0, "cycles")
	genCmd.Flags().Float64P("rate", "r", 0.01, "packet injection rate")
	genCmd.Flags().StringP("file", "f", "", "filepath")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
