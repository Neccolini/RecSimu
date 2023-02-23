package run

import (
	"github.com/Neccolini/RecSimu/cmd/instruction"
	"github.com/Neccolini/RecSimu/cmd/node"
	"github.com/Neccolini/RecSimu/cmd/read"
)

const (
	Add    = "add"
	Remove = "remove"
	Rejoin = "rejoin"
)

func (config *SimulationConfig) changeNetwork(ri read.RecInfo) {
	switch ri.Operation {
	case Add:
		{
			config.nodes[ri.Id], _ = node.NewNode(ri.Id, "Router", []instruction.Instruction{})
			config.adjacencyList[ri.Id] = ri.AdjacencyList
			for _, aNodeId := range ri.AdjacencyList {
				config.adjacencyList[aNodeId] = append(config.adjacencyList[aNodeId], ri.Id)
			}
		}
	case Rejoin:
		{
			config.nodes[ri.Id], _ = node.NewNode(ri.Id, "Router", []instruction.Instruction{})
			config.adjacencyList[ri.Id] = ri.AdjacencyList
			for _, aNodeId := range ri.AdjacencyList {
				config.adjacencyList[aNodeId] = append(config.adjacencyList[aNodeId], ri.Id)
			}
		}
	case Remove:
		{
			config.removeFromAdjacencyList(ri.Id)
			config.nodes[ri.Id].Remove()
			cnId := config.nodes[ri.Id].CommunicatingNodeId
			if cnId != "" {
				config.nodes[cnId].CommunicatingNodeId = ""
			}
		}
	}
}

func (config *SimulationConfig) removeFromAdjacencyList(id string) error {
	config.adjacencyList[id] = nil
	for keyId := range config.nodes {
		config.adjacencyList[keyId] = remove(config.adjacencyList[keyId], id)
	}
	return nil
}

func remove(strings []string, target string) []string {
	result := []string{}
	for _, v := range strings {
		if v != target {
			result = append(result, v)
		}
	}
	return result
}
