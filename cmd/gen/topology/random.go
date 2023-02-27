package gen

import (
	"log"
	"math/rand"
	"strconv"
)

type randomTopology struct {
	idList        []string
	nodesPositionMap map[position]string
	spaces        []position
}

type position struct {
	y int
	x int
}

func RandomNetwork(nodeNum int) map[string][]string {
	idList := nodeIdGen(nodeNum)
	rt := NewRandomTopology(idList)
	rt.create()
	adjacencyList := rt.buildNetwork()
	return adjacencyList
}

func nodeIdGen(nodeNum int) []string {
	idList := make([]string, nodeNum)

	for i := 0; i < nodeNum; i++ {
		idList[i] = strconv.Itoa(i)
	}
	rand.Shuffle(nodeNum, func(i, j int) {
		idList[i], idList[j] = idList[j], idList[i]
	})
	return idList
}

func NewRandomTopology(idList []string) *randomTopology {
	initPos := position{0, 0}
	return &randomTopology{idList, map[position]string{}, []position{initPos}}
}

func (rt *randomTopology) create() error {
	for _, id := range rt.idList {
		rt.occupyRandomSpace(id)
	}
	return nil
}

func (rt *randomTopology) occupyRandomSpace(id string) {
	if len(rt.spaces) == 0 {
		log.Fatal("No space left")
	}
	occupySpaceIndex := rand.Intn(len(rt.spaces))

	// 割り当て
	rt.nodesPositionMap[rt.spaces[occupySpaceIndex]] = id

	// 隣接する新たなspaceを登録
	adjacentPositions := rt.spaces[occupySpaceIndex].adjacentPos()
	for _, adjacentPosition := range adjacentPositions {
		if _, ok := rt.nodesPositionMap[adjacentPosition]; !ok {
			// もし占有されていなかったら
			rt.spaces = append(rt.spaces, adjacentPosition)
		}
	}
	// 削除
	lastIndex := len(rt.spaces) - 1
	rt.spaces[occupySpaceIndex] = rt.spaces[lastIndex]
	rt.spaces = rt.spaces[:lastIndex]

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

func (pos *position) adjacentPos() []position {
	return []position{
		{pos.y - 1, pos.x},
		{pos.y, pos.x + 1},
		{pos.y + 1, pos.x},
		{pos.y, pos.x - 1},
	}
}