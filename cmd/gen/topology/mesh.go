package gen

import (
	"github.com/Neccolini/RecSimu/cmd/debug"
	"gonum.org/v1/plot/plotter"
)

type meshTopology struct {
	rows             int
	columns          int
	idList           []string
	nodesPositionMap map[position]string
}

func MeshNetwork(rows int, columns int, plotFilePath string) map[string][]string {
	idList := nodeIdGen(rows * columns)
	rt := NewMeshTopology(rows, columns, idList)
	rt.create()
	adjacencyList := rt.buildNetwork()

	plotData := plotter.XYs{}
	for xy := range rt.nodesPositionMap {
		plotData = append(plotData, xy.PlotterXY())
	}
	plotNetwork(plotData, plotFilePath)

	return adjacencyList
}

func NewMeshTopology(rows int, columns int, idList []string) *meshTopology {
	return &meshTopology{rows, columns, idList, map[position]string{}}
}

func (rt *meshTopology) create() error {
	cur := 0
	for _, id := range rt.idList {
		pos := position{cur / rt.columns, cur % rt.columns}
		rt.nodesPositionMap[pos] = id
		cur++
	}
	debug.Debug.Printf("map: %v\n", rt.nodesPositionMap)
	return nil
}

func (rt *meshTopology) buildNetwork() map[string][]string {
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
