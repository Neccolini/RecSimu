package read

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Input struct {
	NodeNum       int                 `json:"num"`
	Cycle         int                 `json:"cycles"`
	AdjacencyList map[string][]string `json:"adjacencies"`
	NodesType     map[string]string   `json:"nodes"`
}

func ReadJsonFile(path string) Input {
	byteArray, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error while reading %s: %v\n", path, err)
	}
	var jsonObj Input
	if err := json.Unmarshal(byteArray, &jsonObj); err != nil {
		log.Fatalf("error while reading %s: %v\n", path, err)
	}
	return jsonObj
}
