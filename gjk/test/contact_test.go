package gjk

import (
  "aragno/gjk"
  "aragno/zero"
  "math"
  "testing"
  "fmt"
  "github.com/stretchr/testify/assert"
)

const (
  Tolerance = 1e-3
)

func TestGetCandidate(t *testing.T) {
  quadA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{8, 4},
                                       zero.Vector2D{14, 4},
                                       zero.Vector2D{14, 9},
                                       zero.Vector2D{8, 9}}}

  quadB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{4, 2},
                                                     zero.Vector2D{12, 2},
                                                     zero.Vector2D{12, 5},
                                                     zero.Vector2D{4, 5}}}

  normal := zero.Vector2D{0, -1}

  candidateA := gjk.GetCandidateEdge(quadA, normal)
  candidateB := gjk.GetCandidateEdge(quadB, normal.Inverse())

  expected_A_pntA := zero.Vector2D{8, 4}
  expected_A_pntB := zero.Vector2D{14, 4}
  expected_A_max_proj := expected_A_pntA

  expected_B_pntA := zero.Vector2D{12, 5}
  expected_B_pntB := zero.Vector2D{4, 5}
  expected_B_max_proj := expected_B_pntA

  assert.True(t, candidateA.PntA.Sub(expected_A_pntA).Magnitude() < Tolerance)
  assert.True(t, candidateA.PntB.Sub(expected_A_pntB).Magnitude() < Tolerance)
  assert.True(t, candidateA.MaximumProjection.Sub(expected_A_max_proj).Magnitude() < Tolerance)

  assert.True(t, candidateB.PntA.Sub(expected_B_pntA).Magnitude() < Tolerance)
  assert.True(t, candidateB.PntB.Sub(expected_B_pntB).Magnitude() < Tolerance)
  assert.True(t, candidateB.MaximumProjection.Sub(expected_B_max_proj).Magnitude() < Tolerance)
}

func TestClip(t *testing.T) {
  reference := gjk.CandidateEdge{zero.Vector2D{8, 4}, zero.Vector2D{8, 4}, zero.Vector2D{14, 4}}
  incident := gjk.CandidateEdge{zero.Vector2D{12, 5}, zero.Vector2D{12, 5}, zero.Vector2D{4, 5}}
  refv := reference.Edge().Normalize()

  assert.Equal(t, refv.X, 1.0)
  assert.Equal(t, refv.Y, 0.0)

  clip_start_distance := refv.Dot(reference.PntA)
  assert.Equal(t, clip_start_distance, 8.0)

  points := gjk.Clip(incident.PntA, incident.PntB, refv, clip_start_distance)

  assert.Equal(t, len(points), 2)
  assert.Equal(t, 12.0, points[0].X)
  assert.Equal(t, 5.0, points[0].Y)
  assert.Equal(t, 8.0, points[1].X)
  assert.Equal(t, 5.0, points[1].Y)
}

func TestManifold(t *testing.T) {
  quadA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{8, 4},
                                       zero.Vector2D{14, 4},
                                       zero.Vector2D{14, 9},
                                       zero.Vector2D{8, 9}}}

  quadB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{4, 2},
                                                     zero.Vector2D{12, 2},
                                                     zero.Vector2D{12, 5},
                                                     zero.Vector2D{4, 5}}}

  normal := zero.Vector2D{0, -1}

  points := gjk.GetContactManifold(quadA, quadB, normal)

  assert.Equal(t, 2, len(points))
  assert.Equal(t, 12.0, points[0].Pnt.X)
  assert.Equal(t, 5.0, points[0].Pnt.Y)
  assert.Equal(t, 8.0, points[1].Pnt.X)
  assert.Equal(t, 5.0, points[1].Pnt.Y)
}

func TestManifoldSinglePoint(t *testing.T) {
  fmt.Println("START0")
  quadA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{2, 8},
                                       zero.Vector2D{6, 4},
                                       zero.Vector2D{9, 7},
                                       zero.Vector2D{5, 11}}}

  quadB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{4, 2},
                                                     zero.Vector2D{12, 2},
                                                     zero.Vector2D{12, 5},
                                                     zero.Vector2D{4, 5}}}

  normal := zero.Vector2D{0.0, -1.0}

  points := gjk.GetContactManifold(quadA, quadB, normal)

  assert.Equal(t, 1, len(points))
  assert.Equal(t, 6.0, points[0].Pnt.X)
  assert.Equal(t, 4.0, points[0].Pnt.Y)
  fmt.Println("END0")
}

func TestManifoldAdjustDepth(t *testing.T) {
  fmt.Println("START")
  quadA := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{9, 4},
                                       zero.Vector2D{13, 3},
                                       zero.Vector2D{14, 7},
                                       zero.Vector2D{10, 8}}}

  quadB := gjk.Polygon{[]zero.Vector2D{zero.Vector2D{4, 2},
                                                     zero.Vector2D{12, 2},
                                                     zero.Vector2D{12, 5},
                                                     zero.Vector2D{4, 5}}}

  normal := zero.Vector2D{-0.243, -0.970}

  points := gjk.GetContactManifold(quadA, quadB, normal)

  assert.Equal(t, 2, len(points))
  assert.Equal(t, 12.0, points[0].Pnt.X)
  assert.Equal(t, 5.0, points[0].Pnt.Y)
  assert.True(t, math.Abs(1.69 - points[0].Depth) < 1e-3)
  assert.Equal(t, 8.0, points[1].Pnt.X)
  assert.Equal(t, 5.0, points[1].Pnt.Y)
  assert.True(t, math.Abs(1.04 - points[1].Depth) < 1e-3)
  fmt.Println("END")
}