package run

import (
	"fmt"
	"log"

	"github.com/Neccolini/RecSimu/cmd/instruction"
	"github.com/Neccolini/RecSimu/cmd/node"
)

type SimulationConfig struct {
	nodeNum       int
	totalCycle    int
	adjacencyList map[int][]int
	nodes         map[int]*node.Node
}

func NewSimulationConfig(nodeNum int, cycle int, adjacencyList map[int][]int) *SimulationConfig {
	config := &SimulationConfig{}
	config.nodeNum = nodeNum
	config.totalCycle = cycle
	config.adjacencyList = adjacencyList
	config.nodes = make(map[int]*node.Node, nodeNum)
	for i := 0; i < nodeNum; i++ {
		nodeI, err := node.NewNode(i, "router", []instruction.Instruction{})
		if err != nil {
			log.Fatal(err)
		}
		config.nodes[i] = nodeI
	}
	return config
}

func (config *SimulationConfig) Simulate(outputFile string) error {
	// サイクルごとのシミュレートを実行
	for cycle := 1; cycle <= config.totalCycle; cycle++ {
		// todo トポロジーの変更

		// シミュレートを実行
		if err := config.SimulateCycle(cycle);err != nil {
			return err
		}
	}
	return nil
}

func (config *SimulationConfig) SimulateCycle(cycle int) error {
	// ノードごとにシミュレート
	for _, node := range config.nodes {
		id := node.Id()

		if err := node.SimulateCycle(); err != nil {
			return fmt.Errorf("error during node %d: %w", id, err)
		}
	}

	// メッセージを配信
	for _, node := range config.nodes {
		if !node.SendMessage.IsEmpty() {
			// メッセージをブロードキャストする
			for _, adjacentNodeId := range config.adjacencyList[node.Id()] {
				config.nodes[adjacentNodeId].Receive(node.SendMessage)
			}
			config.nodes[node.Id()].SendMessage.Clear()
		}
	}
	return nil
}
