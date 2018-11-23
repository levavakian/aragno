package gjk

import (
	"aragno/gjk"
	"aragno/zero"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestClosestPoint(t *testing.T) {
	pnta := zero.Vector2D{-4, -1}
	pntb := zero.Vector2D{1, 3}
	closest := gjk.ClosestPointToOrigin(pnta, pntb)

	assert.True(t, closest.Sub(zero.Vector2D{-1.07, 1.34}).IsZero(.01))
}

func TestTwoTriangleNoCollision(t *testing.T) {
	triangleA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{-1, 1},
	                 	                     zero.Vector2D{-1, -1},
	                     	                 zero.Vector2D{-2, 0}}}

	triangleB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{1, 1},
	                 	                     zero.Vector2D{1, -1},
	                     	                 zero.Vector2D{2, 0}}}

	reportAB := gjk.CheckCollision(triangleA, triangleB)
	reportBA := gjk.CheckCollision(triangleB, triangleA)
	assert.False(t, reportAB.Collision)
	assert.False(t, reportBA.Collision)
}

func TestTwoTriangleCollision(t *testing.T) {
	triangleA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{-1, 1},
	                 	                     zero.Vector2D{-1, -1},
	                     	                 zero.Vector2D{2, 0}}}

	triangleB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{1, 1},
	                 	                     zero.Vector2D{1, -1},
	                     	                 zero.Vector2D{-2, 0}}}

	reportAB := gjk.CheckCollision(triangleA, triangleB)
	reportBA := gjk.CheckCollision(triangleB, triangleA)
	assert.True(t, reportAB.Collision)
	assert.True(t, reportBA.Collision)
}

func TestTwoTriangleTouchCollision(t *testing.T) {
	triangleA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{-1, 1},
	                 	                     zero.Vector2D{-1, -1},
	                     	                 zero.Vector2D{0, 0}}}

	triangleB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{1, 1},
	                 	                     zero.Vector2D{1, -1},
	                     	                 zero.Vector2D{0, 0}}}

	reportAB := gjk.CheckCollision(triangleA, triangleB)
	reportBA := gjk.CheckCollision(triangleB, triangleA)
	assert.True(t, reportAB.Collision)
	assert.True(t, reportBA.Collision)
}

func TestQuadrilateralCollision(t *testing.T) {
	quadA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{-2, 0},
	                 	                 zero.Vector2D{0, 2},
	                     	             zero.Vector2D{2, 0},
	                     	             zero.Vector2D{0, -2}}}

	quadB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{-3, 0},
	                 	                 zero.Vector2D{-1, 2},
	                     	             zero.Vector2D{1, 0},
	                     	             zero.Vector2D{-1, -2}}}

	reportAB := gjk.CheckCollision(quadA, quadB)
	reportBA := gjk.CheckCollision(quadB, quadA)
	assert.True(t, reportAB.Collision)
	assert.True(t, reportBA.Collision)
}

func TestTriPolyNoCollision(t *testing.T) {
	quadA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{4, 11},
	                 	                 zero.Vector2D{9, 9},
	                     	             zero.Vector2D{4, 5}}}

	quadB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{8, 6},
	                 	                 zero.Vector2D{15, 6},
	                     	             zero.Vector2D{13, 1},
	                     	             zero.Vector2D{10, 2}}}

	reportAB := gjk.CheckCollision(quadA, quadB)
	reportBA := gjk.CheckCollision(quadB, quadA)
	assert.False(t, reportAB.Collision)
	assert.False(t, reportBA.Collision)
}

func TestQuadrilateralNoCollision(t *testing.T) {
	quadA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{-1, 0},
	                 	                 zero.Vector2D{0, 1},
	                     	             zero.Vector2D{1, 0},
	                     	             zero.Vector2D{0, -1}}}

	quadB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{-10, 0},
	                 	                 zero.Vector2D{-9, 1},
	                     	             zero.Vector2D{-8, 0},
	                     	             zero.Vector2D{-9, -1}}}

	reportAB := gjk.CheckCollision(quadA, quadB)
	reportBA := gjk.CheckCollision(quadB, quadA)
	assert.False(t, reportAB.Collision)
	assert.False(t, reportBA.Collision)
}
