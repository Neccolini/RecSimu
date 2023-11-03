package gen

import (
	"fmt"
	"path/filepath"
	"strconv"

	injection "github.com/Neccolini/RecSimu/cmd/gen/injection"
	topology "github.com/Neccolini/RecSimu/cmd/gen/topology"
)

type Config struct {
	topology string
	nums     []int
}

func NewConfig(topology string, nums []int) *Config {
	return &Config{topology, nums}
}

func GenerateNetwork(config Config, jsonFilePath string, cycles int, rate float64) error {
	fileExtension := filepath.Ext(jsonFilePath)
	pngFilePath := jsonFilePath[:len(jsonFilePath)-len(fileExtension)] + ".png"

	var nodeNum int
	var adjacencyList map[string][]string

	switch config.topology {
	case "random":
		{
			nodeNum = config.nums[0]
			adjacencyList = topology.RandomNetwork(config.nums[0], pngFilePath)
		}
	case "mesh":
		{
			nodeNum = config.nums[0] * config.nums[1]
			adjacencyList = topology.MeshNetwork(config.nums[0], config.nums[1], pngFilePath)
		}
	}

	nodes := []map[string]string{
		{
			"node_id":   "0",
			"node_type": "coordinator",
		},
	}

	for i := 1; i < nodeNum; i++ {
		nodes = append(nodes, map[string]string{
			"node_id":   strconv.Itoa(i),
			"node_type": "router",
		})
		// nodes[strconv.Itoa(i)] = "router"
	}

	injections := injection.GenerateInjectionPackets(nodeNum, cycles, rate)

	jsonOutput := JsonOutput{nodeNum, nodes, adjacencyList, cycles, injections}
	jsonOutput.WriteToFile(jsonFilePath)

	fmt.Printf("Successfully created %s and %s\n", jsonFilePath, pngFilePath)

	return nil
}
