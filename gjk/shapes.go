package gjk

import (
	"aragno/zero"
	"math"
)

type Polygon struct {
	Pnts []zero.Vector2D
}

func (poly *Polygon) NextIdx(idx int) int {
	new_i := idx + 1
	if new_i >= len(poly.Pnts) {
		return 0
	}
	return new_i
}

func (poly *Polygon) PrevIdx(idx int) int {
	new_i := idx - 1
	if new_i < 0 {
		return len(poly.Pnts) - 1
	}
	return new_i
}

func (poly Polygon) FurthestPoint(dir zero.Vector2D) zero.Vector2D {
	maxDot := -math.MaxFloat64
	cpnt := zero.Vector2D{}
	for _, pnt := range poly.Pnts {
		if dot := pnt.Dot(dir); dot >= maxDot {
			cpnt = pnt
			maxDot = dot
		}
	}
	return cpnt
}
