package gjk

import (
	"aragno/zero"
)

type Collidable interface {
	FurthestPoint(zero.Vector2D) zero.Vector2D
}

type CollisionReport struct {
	ClosestPointShapeA zero.Vector2D
	ClosestPointShapeB zero.Vector2D
	Collision          bool
	Distance           float64
}

type SupportComponents struct {
	ShapeAPnt zero.Vector2D
	ShapeBPnt zero.Vector2D
	Pnt       zero.Vector2D
}

type Simplex struct {
	PntA SupportComponents
	PntB SupportComponents
}

func Support(shapeA Collidable, shapeB Collidable, dir zero.Vector2D) SupportComponents {
	vecA := shapeA.FurthestPoint(dir)
	vecB := shapeB.FurthestPoint(dir.Inverse())
	return SupportComponents{vecA, vecB, vecA.Sub(vecB)}
}

func ClosestPointToOrigin(pnta zero.Vector2D, pntb zero.Vector2D) zero.Vector2D {
	AB := pntb.Sub(pnta)
	if AB.IsZero(zero.Tolerance) {
		return pnta
	}

	AO := pnta.Inverse()

	AOABproj := AB.Dot(AO)
	ABL2 := AB.Dot(AB)

	along := AOABproj / ABL2

	if along <= 0.0 {
		return pnta
	} else if along >= 1.0 {
		return pntb
	}

	return AB.Mult(along).Add(pnta)
}

func OriginInTriangle(a zero.Vector2D, b zero.Vector2D, c zero.Vector2D) bool {
	denom := (b.Y - c.Y)*(a.X - c.X) + (c.X - b.X)*(a.Y - c.Y)
	s := ((b.Y - c.Y)*(-c.X) + (c.X - b.X)*(-c.Y)) / denom
	if s > 1.0 || s < 0.0 {
		return false
	}

	t := ((c.Y - a.Y)*(-c.X) + (a.X - c.X)*(-c.Y)) / denom
	if t > 1.0 || t < 0.0 {
		return false
	}

	u := 1 - s - t
	if u > 1.0 || u < 0.0 {
		return false
	}

	return true
}

func ClosestShapePnts(pnta SupportComponents, pntb SupportComponents) (zero.Vector2D, zero.Vector2D) {
	L := pntb.Pnt.Sub(pnta.Pnt)
	if L.IsZero(zero.Tolerance) {
		return pnta.Pnt, pntb.Pnt
	}
	LL2 := L.Dot(L)
	LA := L.Dot(pnta.Pnt)

	lambda2 := -LA / LL2
	if lambda2 > 1.0 {
		lambda2 = 1.0
	}
	if lambda2 < 0.0 {
		lambda2 = 0.0
	}
	lambda1 := 1 - lambda2

	shapeAClosest := pnta.ShapeAPnt.Mult(lambda1).Add(pntb.ShapeAPnt.Mult(lambda2))
	shapeBClosest := pnta.ShapeBPnt.Mult(lambda1).Add(pntb.ShapeBPnt.Mult(lambda2))
	return shapeAClosest, shapeBClosest
}

func SameExact(veca zero.Vector2D, vecb zero.Vector2D) bool {
	return veca.X == vecb.X && veca.Y == vecb.Y
}

func CheckCollision(shapeA Collidable, shapeB Collidable) CollisionReport {
	// TODO: do line of centers of shapes
	d := zero.Vector2D{1, 1}
	simplex := Simplex{}
	simplex.PntA = Support(shapeA, shapeB, d)
	simplex.PntB = Support(shapeA, shapeB, d.Inverse())
	d = ClosestPointToOrigin(simplex.PntA.Pnt, simplex.PntB.Pnt)

	for i := 1; i <= 20; i++ {
		d = d.Inverse()

		if (d.IsZero(zero.Tolerance)) {
			sApnt, sBpnt := ClosestShapePnts(simplex.PntA, simplex.PntB)
			return CollisionReport{sApnt, sBpnt, true, 0}
		}

		c := Support(shapeA, shapeB, d)
		dc := c.Pnt.Dot(d)
		da := simplex.PntA.Pnt.Dot(d)

		if dc - da < zero.Tolerance || SameExact(simplex.PntA.Pnt, c.Pnt) || SameExact(simplex.PntB.Pnt, c.Pnt) {
			sApnt, sBpnt := ClosestShapePnts(simplex.PntA, simplex.PntB)
			return CollisionReport{sApnt, sBpnt, false, sApnt.Sub(sBpnt).Magnitude()}
		}

		p1 := ClosestPointToOrigin(simplex.PntA.Pnt, c.Pnt)
		p2 := ClosestPointToOrigin(c.Pnt, simplex.PntB.Pnt)

		collision := OriginInTriangle(simplex.PntA.Pnt, simplex.PntB.Pnt, c.Pnt)
		if (p1.Magnitude() < p2.Magnitude()) {
			simplex.PntB = c
			d = p1
		} else {
			simplex.PntA = c
			d = p2
		}

		if collision {
			sApnt, sBpnt := ClosestShapePnts(simplex.PntA, simplex.PntB)
			return CollisionReport{sApnt, sBpnt, true, 0}
		}
	}

	return CollisionReport{zero.Vector2D{}, zero.Vector2D{}, false, 0}
}
