package gen

import (
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/Neccolini/RecSimu/cmd/injection"
)

const FlitLen = 6

func GenerateInjectionPackets(nodeNum int, maxCycles int, rate float64) []injection.Injection {
	baseCycles := nodeNum * 100
	if baseCycles >= maxCycles+1000 {
		log.Fatalf("The number of cycles is too small")
	}
	cycles := maxCycles - baseCycles
	packetsNum := totalPacketNum(cycles, rate)
	interval := cycles / packetsNum
	res := []injection.Injection{}

	for i := 0; i < nodeNum; i++ {
		nodeId := strconv.Itoa(i)
		for j := baseCycles; j <= maxCycles; j += interval {
			randId := rand.Intn(nodeNum)
			if randId == i {
				offset := rand.Intn(nodeNum-1) + 1
				randId = (randId + offset) % nodeNum
			}
			toId := strconv.Itoa(randId)

			randCycle := rand.Intn(interval-2) + 1
			injection := generatePacket(j, nodeId, toId, j+randCycle)
			res = append(res, injection)
		}
	}
	return res
}

func totalPacketNum(cycles int, rate float64) int {
	return int(math.Ceil(float64(cycles) * rate / float64(FlitLen)))
}

func generatePacket(injectionNum int, fromId string, toId string, cycle int) injection.Injection {
	InjectionId := "Inj_" + fromId + "_" + strconv.Itoa(injectionNum)
	return injection.Injection{InjectionId: InjectionId, Cycle: cycle, FromId: fromId, DistId: toId, Data: ""}
}
