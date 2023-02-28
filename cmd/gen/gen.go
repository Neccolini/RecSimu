package gen

import (
	"path/filepath"
	"strconv"

	gen "github.com/Neccolini/RecSimu/cmd/gen/topology"
)

func GenerateNetwork(jsonFilePath string, num int, cycles int) error {
	fileExtension := filepath.Ext(jsonFilePath)
	pngFilePath := jsonFilePath[:len(jsonFilePath)-len(fileExtension)] + ".png"

	adjacencyList := gen.RandomNetwork(num, pngFilePath)

	nodes := map[string]string{
		"0": "Coordinator",
	}

	for i := 1; i < num; i++ {
		nodes[strconv.Itoa(i)] = "Router"
	}

	jsonOutput := JsonOutput{num, nodes, adjacencyList, cycles}
	jsonOutput.WriteToFile(jsonFilePath)

	return nil
}
