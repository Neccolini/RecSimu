package run

import (
	"fmt"
	"sync"
	"time"
)

type SimulationConfig struct {
	NodeNum       int
	NodeInfo      []NodeInitInfo
	Cycle         int
	AdjacencyList [][]int
	Messages      []chan MessageBuffer
	ControlNodes  []chan Control
}

type Control struct {
	Run           bool
	Exit          bool
	MessageBuffer MessageBuffer
}

func NewSimulationConfig(nodeNum int, nodeInfo []NodeInitInfo, cycle int, adjacencyList [][]int) *SimulationConfig {
	config := new(SimulationConfig)
	config.NodeNum = nodeNum
	config.NodeInfo = nodeInfo
	config.Cycle = cycle
	config.AdjacencyList = adjacencyList
	config.Messages = make([]chan MessageBuffer, nodeNum)
	config.ControlNodes = make([]chan Control, nodeNum)

	for i := 0; i < nodeNum; i++ {
		config.Messages[i] = make(chan MessageBuffer, 2) // todo channel buffer size?
		config.ControlNodes[i] = make(chan Control, 2)
	}
	return config
}

func (config *SimulationConfig) Simulate() {
	var wg = new(sync.WaitGroup)

	for i := 0; i < config.NodeNum; i++ {
		go Node(config.NodeInfo[i], config.ControlNodes[i], wg)
	}

	for cycle := 0; cycle < config.Cycle; cycle++ {
		// 各ノードの動作を開始
		wg.Add(config.NodeNum)
		for i := 0; i < config.NodeNum; i++ {
			config.ControlNodes[i] <- Control{Run: true, Exit: false, MessageBuffer: MessageBuffer{}}
		}
		// 各ノードの動作の終了
		wg.Wait()
		// メインルーチンの動作開始
		fmt.Println("Main begin")
		for i := 0; i < config.NodeNum; i++ {
			fmt.Println(<-config.ControlNodes[i])
		}
		// config.DeliverMessage()
		time.Sleep(1000 * time.Millisecond) // 処理
		fmt.Println("Main end")
		// メインルーチンの動作終了
	}
	for i := 0; i < config.NodeNum; i++ {
		config.ControlNodes[i] <- Control{false, true, MessageBuffer{}}
	}
	wg.Wait()
	return
}

func (config *SimulationConfig) DeliverMessage() {
	// 各ノードについて
	for i := 0; i < config.NodeNum; i++ {
		select {
		case control := <-config.ControlNodes[i]:
			// メッセージが届いていたら
			if len(control.MessageBuffer.Messages) != 0 {
				// メッセージバッファのメッセージを一つずつ読む
				for _, message := range control.MessageBuffer.Messages {
					// メッセージが入っていたら
					if !message.IsEmpty() {
						// 隣接リストを参照しメッセージを配信
						for _, node := range config.AdjacencyList[message.Id] {
							messageBuffer := <-config.Messages[node]
							messageBuffer.Messages = append(messageBuffer.Messages, message)
							config.Messages[node] <- messageBuffer
						}
					}
				}
			}
		default:
			fmt.Println("No data received")
		}
	}
}
