package history

import "math"

type position struct {
	index int
	layer int
}

func getLeftNode(pos position) position {
	return position{pos.index, pos.layer - 1}
}

func getRightNode(pos position) position {
	index := pos.index + int(math.Exp2(float64(pos.layer-1)))
	layer := pos.layer - 1
	return position{index, layer}
}
