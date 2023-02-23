package read

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Input struct {
	NodeNum         int                 `json:"num"`
	Cycle           int                 `json:"cycles"`
	AdjacencyList   map[string][]string `json:"adjacencies"`
	NodesType       map[string]string   `json:"nodes"`
	ReconfigureInfo []RecInfo           `json:"reconfigure"`
}

type InputShaped struct {
	NodeNum         int                 `json:"num"`
	Cycle           int                 `json:"cycles"`
	AdjacencyList   map[string][]string `json:"adjacencies"`
	NodesType       map[string]string   `json:"nodes"`
	ReconfigureInfo map[int][]RecInfo
}

func ReadJsonFile(path string) InputShaped {
	byteArray, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error while reading %s: %v\n", path, err)
	}
	var i Input
	if err := json.Unmarshal(byteArray, &i); err != nil {
		log.Fatalf("error while reading %s: %v\n", path, err)
	}

	mapedRecInfo := MapCycle2RecInfo(i.ReconfigureInfo)

	return InputShaped{
		NodeNum:         i.NodeNum,
		Cycle:           i.Cycle,
		AdjacencyList:   i.AdjacencyList,
		NodesType:       i.NodesType,
		ReconfigureInfo: mapedRecInfo,
	}
}
