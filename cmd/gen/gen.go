package gen

import (
	"fmt"
	"path/filepath"
	"strconv"

	injection "github.com/Neccolini/RecSimu/cmd/gen/injection"
	topology "github.com/Neccolini/RecSimu/cmd/gen/topology"
)

func GenerateNetwork(jsonFilePath string, num int, cycles int, rate float64) error {
	fileExtension := filepath.Ext(jsonFilePath)
	pngFilePath := jsonFilePath[:len(jsonFilePath)-len(fileExtension)] + ".png"

	adjacencyList := topology.RandomNetwork(num, pngFilePath)

	nodes := map[string]string{
		"0": "Coordinator",
	}

	for i := 1; i < num; i++ {
		nodes[strconv.Itoa(i)] = "Router"
	}

	injections := injection.GenerateInjectionPackets(num, cycles, rate)

	jsonOutput := JsonOutput{num, nodes, adjacencyList, cycles, injections}
	jsonOutput.WriteToFile(jsonFilePath)

	fmt.Printf("Successfully created %s and %s\n", jsonFilePath, pngFilePath)

	return nil
}
