package gen

import (
	"encoding/json"
	"log"
	"os"

	injection "github.com/Neccolini/RecSimu/cmd/injection"
)

type JsonOutput struct {
	Num         int                   `json:"node_num"`
	Nodes       []map[string]string   `json:"nodes"`
	Adjacencies map[string][]string   `json:"neighbors"`
	Cycles      int                   `json:"total_cycles"`
	Injections  []injection.Injection `json:"packets"`
}

type Node struct {
	Id       string
	NodeType string
}

func (j *JsonOutput) WriteToFile(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	if err := encoder.Encode(*j); err != nil {
		log.Fatal(err)
	}
	return nil
}
