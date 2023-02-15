package run

import (
	"fmt"
	"sync"
	"time"
)

type NodeInitInfo struct {
	Id          int
	NodeType    string
	MessageList []Message
}

func Node(nodeInitInfo NodeInitInfo, control chan Control, wg *sync.WaitGroup) {
	running := false
	var c Control
	for {
		select {
		case c = <-control:
			if c.Run == true {
				running = true
			}
			if c.Exit == true {
				return
			}
		default:
		}
		if running == true {
			fmt.Printf("node %d has started\n", nodeInitInfo.Id)
			time.Sleep(1000 * time.Millisecond) // 処理
			fmt.Printf("node %d has stopped\n", nodeInitInfo.Id)

			control <- Control{Run: false, Exit: false, MessageBuffer: MessageBuffer{}}
			running = false
			wg.Done()
		}
	}
}
