package gen

import (
	"math/rand"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type position struct {
	y int
	x int
}

func (pos *position) adjacentPos() []position {
	return []position{
		{pos.y - 1, pos.x},
		{pos.y, pos.x + 1},
		{pos.y + 1, pos.x},
		{pos.y, pos.x - 1},
	}
}

func (pos *position) PlotterXY() plotter.XY {
	return plotter.XY{X: float64(pos.x), Y: float64(pos.y)}
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

func plotNetwork(xys plotter.XYs, plotFilePath string) {
	p := plot.New()
	s, _ := plotter.NewScatter(xys)
	p.Add(s, plotter.NewGrid())
	p.Save(4*vg.Inch, 4*vg.Inch, plotFilePath)
}
