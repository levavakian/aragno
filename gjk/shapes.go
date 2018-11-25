package gjk

import (
	"aragno/zero"
	"math"
)

type Polygon struct {
	Pnts []zero.Vector2D
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
