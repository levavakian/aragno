package gjk

import (
  "aragno/zero"
  "math"
)

type ContactPoint struct {
  Pnt   zero.Vector2D
  Depth float64
}

type CandidateEdge struct {
  MaximumProjection zero.Vector2D
  PntA zero.Vector2D
  PntB zero.Vector2D
}

func (candidate *CandidateEdge) Edge() zero.Vector2D {
  return candidate.PntB.Sub(candidate.PntA)
}

func GetContactManifold(shapea Polygon, shapeb Polygon, normal zero.Vector2D) []ContactPoint {
  // Get candidate edges for the manifold
  candidate_a := GetCandidateEdge(shapea, normal)
  candidate_b := GetCandidateEdge(shapeb, normal.Inverse())

  // Choose reference and incident edges
  reference := CandidateEdge{}
  incident := CandidateEdge{}

  if math.Abs(candidate_a.Edge().Dot(normal)) <= math.Abs(candidate_b.Edge().Dot(normal)) {
    reference = candidate_a
    incident = candidate_b
  } else {
    reference = candidate_b
    incident = candidate_a
  }
  ref_unit_vec := reference.Edge().Normalize()

  // Do the initial clip on one side of the clipping edge
  clipped_points := Clip(incident.PntA, incident.PntB, ref_unit_vec, ref_unit_vec.Dot(reference.PntA))
  if len(clipped_points) < 2 {
    return []ContactPoint{}
  }

  // Do the seconday clip on the other side of the clipping edge
  clipped_points = Clip(clipped_points[0], clipped_points[1], ref_unit_vec.Inverse(), -ref_unit_vec.Dot(reference.PntB))
  if len(clipped_points) < 2 {
    return []ContactPoint{}
  }

  // Do final fixup
  ref_norm := ref_unit_vec.Cross(-1.0)

  max_penetration := ref_norm.Dot(reference.MaximumProjection)
  final := []ContactPoint{}
  
  p0_depth := ref_norm.Dot(clipped_points[0]) - max_penetration
  if p0_depth >= 0.0 {
    final = append(final, ContactPoint{clipped_points[0], p0_depth})
  }
  
  p1_depth := ref_norm.Dot(clipped_points[1]) - max_penetration
  if p1_depth >= 0.0 {
    final = append(final, ContactPoint{clipped_points[1], p1_depth})
  }

  return final
}

func GetCandidateEdge(shape Polygon, normal zero.Vector2D) CandidateEdge {
  max_proj := normal.Dot(shape.Pnts[0])
  max_idx := 0

  // Find point with maximum distance along normal
  for i, pnt := range shape.Pnts {
    projection := normal.Dot(pnt)
    if projection > max_proj {
      max_proj = projection
      max_idx = i
    }
  }

  // Get points before and after max projection
  v := shape.Pnts[max_idx]
  v1 := shape.Pnts[shape.NextIdx(max_idx)]
  v0 := shape.Pnts[shape.PrevIdx(max_idx)]

  // Get edges pointing left and right
  left := v.Sub(v1).Normalize()
  right := v.Sub(v0).Normalize()

  // Choose edge that is more perpendicular to normal
  if right.Dot(normal) <= left.Dot(normal) {
    return CandidateEdge{v, v0, v}
  } else {
    return CandidateEdge{v, v, v1}
  }
}

func Clip(v0 zero.Vector2D, v1 zero.Vector2D, refv zero.Vector2D, clip_start_distance float64) []zero.Vector2D {
  points := []zero.Vector2D{}

  dist_v0 := refv.Dot(v0) - clip_start_distance
  dist_v1 := refv.Dot(v1) - clip_start_distance

  // If either point is greater than the clip start distance, keep them
  if dist_v0 >= 0.0 {
    points = append(points, v0)
  }
  if dist_v1 >= 0.0 {
    points = append(points, v1)
  }

  // Check if the points are on either side of the reference edge
  if dist_v0 * dist_v1 < 0.0 {
    // Distance ratio of v0 along v0 -> v1
    v0_dist_ratio := dist_v0 / (dist_v0 - dist_v1)
    clip_point := v1.Sub(v0).Mult(v0_dist_ratio).Add(v0)
    points = append(points, clip_point)
  }

  return points
}