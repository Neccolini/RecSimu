package run

import (
	"fmt"

	"github.com/Neccolini/RecSimu/cmd/instruction"
	"github.com/Neccolini/RecSimu/cmd/message"
	"github.com/Neccolini/RecSimu/cmd/node"
	"github.com/Neccolini/RecSimu/cmd/random"
	"github.com/Neccolini/RecSimu/cmd/read"
	"github.com/Neccolini/RecSimu/cmd/routing"
)

type SimulationConfig struct {
	nodeNum         int
	totalCycle      int
	adjacencyList   map[string][]string
	nodes           map[string]*node.Node
	// instructionList []instruction.Instruction
	recInfo map[int][]read.RecInfo
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
		for _, info := range config.recInfo[cycle] {
			config.changeNetwork(info)
		}
		// シミュレートを実行
		if err := config.SimulateCycle(cycle); err != nil {
			return err
		}
		// todo 各サイクル後の状態を記録
		fmt.Printf("cycle %d\n", cycle)
		for _, node := range config.nodes {
			fmt.Println(node.String())
		}
	}
	return nil
}

func (config *SimulationConfig) SimulateCycle(cycle int) error {
	messageMap := map[string][]message.Message{}

	for _, node := range config.nodes {
		if !node.Alive() {
			continue
		}
		// ノードごとに送信
		node.CycleSend()
		if node.SendingMessage.IsEmpty() {
			continue
		}
		// Broadcastの場合
		if node.SendingMessage.ToId() == routing.BroadCastId {
			success := false
			for _, aNodeId := range config.adjacencyList[node.Id()] {
				if config.nodes[aNodeId].State().IsIdle() && config.nodes[aNodeId].IsJoined() {
					success = true
					messageMap[aNodeId] = append(messageMap[aNodeId], node.SendingMessage)
				}
			}
			if !success {
				fmt.Println(node.Id())
				config.nodes[node.Id()].Wait()
			}
		} else { // Broadcastでない場合
			for _, aNodeId := range config.adjacencyList[node.Id()] {
				if config.nodes[aNodeId].State().IsIdle() {
					messageMap[aNodeId] = append(messageMap[aNodeId], node.SendingMessage)
				} else if aNodeId == node.SendingMessage.ToId() { // 送信先は一つのはずなのに、そうでない方への送信が失敗すると大気モードに入るのはおかしい
					config.nodes[node.Id()].Wait() // 送信に失敗したので待機モード
				}
			}
		}
	}

	// 送信メッセージを集計
	for rNodeId, msgs := range messageMap {
		if len(msgs) > 0 {
			successMsg := random.RandomChoice(msgs)   // 複数あった場合，一つランダムで選択
			config.nodes[rNodeId].Receive(successMsg) // 受信
			for _, failedMsg := range msgs {
				if failedMsg.Id() != successMsg.Id() {
					config.nodes[failedMsg.Id()].Wait() // 送信に失敗したので待機モード
				}
			}
		}
	}

	for _, node := range config.nodes {
		if !node.Alive() {
			continue
		}
		node.CycleReceive()
		node.SimulateCycle()
	}

	return nil
}
