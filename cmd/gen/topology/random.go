package gen

import (
	"log"

	"github.com/Neccolini/RecSimu/cmd/set"
	"gonum.org/v1/plot/plotter"
)

type randomTopology struct {
	idList           []string
	nodesPositionMap map[position]string
	spaces           set.Set[position]
}

func RandomNetwork(nodeNum int, plotFilePath string) map[string][]string {
	idList := nodeIdGen(nodeNum)
	rt := NewRandomTopology(idList)
	rt.create()
	adjacencyList := rt.buildNetwork()

	plotData := plotter.XYs{}
	for xy := range rt.nodesPositionMap {
		plotData = append(plotData, xy.PlotterXY())
	}
	plotNetwork(plotData, plotFilePath)

	return adjacencyList
}

func NewRandomTopology(idList []string) *randomTopology {
	initPos := position{0, 0}
	return &randomTopology{idList, map[position]string{}, *set.NewSet(initPos)}
}

func (rt *randomTopology) create() error {
	for _, id := range rt.idList {
		rt.occupyRandomSpace(id)
	}
	return nil
}

func (rt *randomTopology) occupyRandomSpace(id string) position {
	if rt.spaces.Size() == 0 {
		log.Fatal("No space left")
	}
	pos := rt.spaces.RandomChoice()

	// 割り当て
	rt.nodesPositionMap[pos] = id

	// 隣接する新たなspaceを登録
	adjacentPositions := pos.adjacentPos()
	for _, adjacentPosition := range adjacentPositions {
		if _, ok := rt.nodesPositionMap[adjacentPosition]; !ok {
			// もし占有されていなかったら
			rt.spaces.Add(adjacentPosition)
		}
	}
	// 削除
	rt.spaces.Remove(pos)

	return pos
}

func (rt *randomTopology) buildNetwork() map[string][]string {
	adjacencyList := map[string][]string{}
	for pos, id := range rt.nodesPositionMap {
		// このposに隣接するノードを探す
		resList := []string{}
		// 周囲4方向を探す
		adjacentPositions := pos.adjacentPos()

		for _, ap := range adjacentPositions {
			// その方向に存在していればリストに追加
			if aId, ok := rt.nodesPositionMap[ap]; ok {
				resList = append(resList, aId)
			}
		}
		adjacencyList[id] = resList
	}
	return adjacencyList
}
