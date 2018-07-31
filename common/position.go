package common

import (
	"math"
)

type Position struct {
	Index int
	Layer int
}

func NewPosition(index, layer int) Position {
	return Position{index, layer}
}

func (pos Position) GetLeftNode() Position {
	return Position{pos.Index, pos.Layer - 1}
}

func (pos Position) GetRightNode() Position {
	index := pos.Index + int(math.Exp2(float64(pos.Layer-1)))
	layer := pos.Layer - 1
	return Position{index, layer}
}

func (pos Position) ComputeNodeValue() float64 {
	return float64(pos.Index) + math.Exp2(float64(pos.Layer-1))
}
