package run

import (
	"fmt"
	"log"

	"github.com/Neccolini/RecSimu/cmd/debug"
	"github.com/Neccolini/RecSimu/cmd/injection"
	"github.com/Neccolini/RecSimu/cmd/message"
	routing "github.com/Neccolini/RecSimu/cmd/network/routing"
	"github.com/Neccolini/RecSimu/cmd/node"
	"github.com/Neccolini/RecSimu/cmd/rec"
	"github.com/Neccolini/RecSimu/cmd/utils"
)

type SimulationConfig struct {
	nodeNum        int
	totalCycle     int
	adjacencyList  map[string][]string
	nodes          map[string]*node.Node
	injectionTable injection.InjectionTable
	recInfo        map[int][]rec.RecInfo
	fromId2ToId    map[string][]string // fromId -> toId
	toId2FromId    map[string][]string // toId -> fromId
}

func NewSimulationConfig(nodeNum int, cycle int, adjacencyList map[string][]string, nodesType map[string]string, recInfo map[int][]rec.RecInfo, iTable injection.InjectionTable) *SimulationConfig {
	config := &SimulationConfig{}
	config.nodeNum = nodeNum
	config.totalCycle = cycle
	config.adjacencyList = adjacencyList
	config.nodes = make(map[string]*node.Node, nodeNum)
	config.recInfo = recInfo
	config.injectionTable = iTable

	for id, nType := range nodesType {
		config.nodes[id], _ = node.NewNode(id, nType) // todo エラー処理？
	}

	return config
}

func (config *SimulationConfig) Simulate() error {
	// サイクルごとのシミュレートを実行
	for cycle := 1; cycle <= config.totalCycle; cycle++ {
		// todo トポロジーの変更
		config.changeNode(cycle)

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

	// シミュレーション結果の集計
	averageLatency := 0.0
	totalPackets := 0
	failedPackets := 0
	for _, node := range config.nodes {
		averageLatency += float64(node.Performance.TotalLatency) / float64(node.Performance.TotalPacketNum)
		totalPackets += node.Performance.TotalPacketNum
		failedPackets += node.Performance.FailedPacketNum()
		if node.Performance.RecResult() != nil {
			fmt.Printf("reconfiguration: %s %v\n", node.Id(), node.Performance.RecResult())
		}
	}
	fmt.Printf("total packets: %d / %d\n", totalPackets-failedPackets, totalPackets)
	fmt.Printf("average latency: %.5f [cycle]\n", averageLatency/float64(config.nodeNum))

	return nil
}

func (config *SimulationConfig) SimulateCycle(cycle int) error {
	config.toId2FromId = map[string][]string{}
	config.fromId2ToId = map[string][]string{}

	config.inject(cycle)

	// サイクルの更新
	for _, node := range config.nodes {
		if !node.Alive() {
			continue
		}
		node.SimulateCycle(cycle)
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

func (config *SimulationConfig) inject(cycle int) error {
	injections, err := config.injectionTable.QueryByCycle(cycle)
	if err != nil {
		return err
	}
	for _, i := range injections {
		node := config.nodes[i.FromId]
		if node.IsJoined() {
			node.InjectMessage(i.InjectionId, i.DistId, i.Data)
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
		msg := config.nodes[fromId].SendingMessage
		if len(toIds) == 0 {
			continue
		}
		config.nodes[fromId].SetSending()

		injection, _ := config.injectionTable.QueryById(msg.MessageId)
		for _, toId := range toIds {
			if injection.DistId == toId {
				// perfomance measure end
				config.nodes[injection.FromId].MessageReached(injection.InjectionId)
			}
			config.nodes[toId].SetReceiving(msg)
		}
	}
}

func (config *SimulationConfig) changeNode(cycle int) {
	for _, recInfo := range config.recInfo[cycle] {
		switch recInfo.Operation {
		case rec.Add, rec.Rejoin:
			{
				config.AddNode(recInfo)
			}
		case rec.Remove:
			{
				config.RemoveNode(recInfo.Id, cycle)
			}
		}
	}
}

func (config *SimulationConfig) AddNode(recInfo rec.RecInfo) {
	config.adjacencyList[recInfo.Id] = make([]string, 0, len(recInfo.AdjacencyList))
	for _, aId := range recInfo.AdjacencyList {
		config.adjacencyList[recInfo.Id] = append(config.adjacencyList[recInfo.Id], aId)
		config.adjacencyList[aId] = append(config.adjacencyList[aId], recInfo.Id)
	}
	config.nodes[recInfo.Id], _ = node.NewNode(recInfo.Id, recInfo.NodeType)
}

func (config *SimulationConfig) RemoveNode(id string, cycle int) {
	// 接続を削除
	for _, aId := range config.adjacencyList[id] {
		config.adjacencyList[aId] = utils.Remove(config.adjacencyList[aId], id)
	}
	for _, aId := range config.adjacencyList[id] {
		if config.nodes[aId].RoutingFunction.ParentId() == id {
			config.nodes[aId].Performance.RecStart(cycle)
			config.nodes[aId].InitReconfiguration(id)
		}
	}
	config.adjacencyList[id] = nil
	config.nodes[id].Remove()
}
