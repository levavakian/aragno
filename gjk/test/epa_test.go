package gjk

import (
	"aragno/gjk"
	"aragno/zero"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestTriPolyCollisionEpacd(t *testing.T) {
	quadA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{4, 11},
		zero.Vector2D{9, 9},
		zero.Vector2D{4, 5}}}

	quadB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{5, 7},
		zero.Vector2D{12, 7},
		zero.Vector2D{10, 2},
		zero.Vector2D{7, 3}}}

	pnta := gjk.SupportComponents{zero.Vector2D{}, zero.Vector2D{}, zero.Vector2D{4, 2}}
	pntb := gjk.SupportComponents{zero.Vector2D{}, zero.Vector2D{}, zero.Vector2D{-8, -2}}
	pntc := gjk.SupportComponents{zero.Vector2D{}, zero.Vector2D{}, zero.Vector2D{-1, -2}}

	gjk.GetPenetrationInfo(quadA, quadB, pnta, pntb, pntc)
}
