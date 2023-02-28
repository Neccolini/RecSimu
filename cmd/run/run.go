package run

import (
	"log"

	"github.com/Neccolini/RecSimu/cmd/debug"
	"github.com/Neccolini/RecSimu/cmd/instruction"
	"github.com/Neccolini/RecSimu/cmd/message"
	"github.com/Neccolini/RecSimu/cmd/node"
	"github.com/Neccolini/RecSimu/cmd/read"
	"github.com/Neccolini/RecSimu/cmd/routing"
)

type SimulationConfig struct {
	nodeNum       int
	totalCycle    int
	adjacencyList map[string][]string
	nodes         map[string]*node.Node
	// instructionList []instruction.Instruction
	recInfo    map[int][]read.RecInfo
	messageMap map[string][]message.Message
}

func NewSimulationConfig(nodeNum int, cycle int, adjacencyList map[string][]string, nodesType map[string]string, recInfo map[int][]read.RecInfo) *SimulationConfig {
	config := &SimulationConfig{}
	config.nodeNum = nodeNum
	config.totalCycle = cycle
	config.adjacencyList = adjacencyList
	config.nodes = make(map[string]*node.Node, nodeNum)
	config.recInfo = recInfo
	for id, nType := range nodesType {
		config.nodes[id], _ = node.NewNode(id, nType, []instruction.Instruction{}) // todo エラー処理？
	}

	return config
}

func (config *SimulationConfig) Simulate(outputFile string) error {
	// サイクルごとのシミュレートを実行
	for cycle := 1; cycle <= config.totalCycle; cycle++ {
		// todo トポロジーの変更

		// シミュレートを実行
		if err := config.SimulateCycle(cycle); err != nil {
			return err
		}
		// todo 各サイクル後の状態を記録
		debug.Debug.Printf("cycle %d\n", cycle)
		for _, node := range config.nodes {
			debug.Debug.Println(node.String())
		}
	}
	return nil
}

func (config *SimulationConfig) SimulateCycle(cycle int) error {
	config.messageMap = map[string][]message.Message{}

	// サイクルの更新
	for _, node := range config.nodes {
		if !node.Alive() {
			continue
		}
		node.SimulateCycle()
	}

	// メッセージをブロードキャスト
	for _, node := range config.nodes {
		if node.IsSendTrying() {
			msg := node.SendingMessage
			if msg.Id() != node.Id() {
				log.Fatalf("msgId and nodeId is different: %s %s", msg.Id(), node.Id())
			}
			config.broadCastMessage(msg)
		}
	}

	// 受信側の処理：複数から送られてきた場合失敗
	for distId, msgs := range config.messageMap {
		distNode := config.nodes[distId]
		if !distNode.IsIdle() && !distNode.IsWaiting() {
			continue
		}
		if len(msgs) == 1 {
			msg := msgs[0]
			if !config.nodes[msg.Id()].IsSending() {
				config.nodes[msg.Id()].SetSending()
			}
			if !config.nodes[distId].IsReceiving() {
				config.nodes[distId].SetReceiving(msg)
			}
		}
	}
	return nil
}

func (config *SimulationConfig) broadCastMessage(msg message.Message) {
	node := config.nodes[msg.Id()]
	success := false
	for _, adjacentId := range config.adjacencyList[node.Id()] {
		adjacentNode := config.nodes[adjacentId]
		if adjacentId == msg.ToId() && !adjacentNode.IsIdle() && !adjacentNode.IsWaiting() {
			// 宛先がある場合，その宛先がbusyなら失敗
			return
		}
		if msg.ToId() == routing.BroadCastId && adjacentNode.IsJoined() &&
			(adjacentNode.IsIdle() || adjacentNode.IsWaiting()) {
			success = true
		}
	}

	if msg.ToId() == routing.BroadCastId && !success {
		return
	}

	for _, adjacentId := range config.adjacencyList[node.Id()] {
		config.messageMap[adjacentId] = append(config.messageMap[adjacentId], msg)
	}
}
