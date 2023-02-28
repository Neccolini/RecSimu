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
	recInfo     map[int][]read.RecInfo
	fromId2ToId map[string][]string // fromId -> toId
	toId2FromId map[string][]string // toId -> fromId
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
		debug.Debug.Printf("\ncycle %d\n", cycle)
		for _, node := range config.nodes {
			debug.Debug.Println(node.String())
		}
	}
	return nil
}

func (config *SimulationConfig) SimulateCycle(cycle int) error {
	config.toId2FromId = map[string][]string{}
	config.fromId2ToId = map[string][]string{}

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

	config.deliverMessages()

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
		if !config.nodes[adjacentId].IsIdle() && !config.nodes[adjacentId].IsWaiting() {
			continue
		}
		config.toId2FromId[adjacentId] = append(config.toId2FromId[adjacentId], node.Id())
		config.fromId2ToId[node.Id()] = append(config.fromId2ToId[node.Id()], adjacentId)
	}
}

func (config *SimulationConfig) deliverMessages() {
	// 受信側で複数届いていたら失敗
	for _, fromIds := range config.toId2FromId {
		if len(fromIds) >= 2 {
			for _, fromId := range fromIds {
				delete(config.fromId2ToId, fromId)
			}
		}
	}

	// 送信成功
	for fromId, toIds := range config.fromId2ToId {
		if len(toIds) == 0 {
			continue
		}
		config.nodes[fromId].SetSending()
		for _, toId := range toIds {
			config.nodes[toId].SetReceiving(config.nodes[fromId].SendingMessage)
		}
	}
}
