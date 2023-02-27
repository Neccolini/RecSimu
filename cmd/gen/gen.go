package gen

import (
	"strconv"

	gen "github.com/Neccolini/RecSimu/cmd/gen/topology"
)

func GenerateNetwork(filepath string, num int, cycles int) error {
	adjacencyList := gen.RandomNetwork(num)

	nodes := map[string]string{
		"0": "Coordinator",
	}

	for i := 1; i < num; i++ {
		nodes[strconv.Itoa(i)] = "Router"
	}

	jsonOutput := JsonOutput{num, nodes, adjacencyList, cycles}
	jsonOutput.WriteToFile(filepath)

	return nil
}
