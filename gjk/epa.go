package gjk

import (
	"aragno/zero"
	"math"
)

type EPANode struct {
	Supp SupportComponents
	Next *EPANode
}

func (node *EPANode) Insert(supp SupportComponents) {
	insert := &EPANode{supp, node.Next}
	node.Next = insert
}

type Edge struct {
	Normal   zero.Vector2D
	Distance float64
	Parent   *EPANode
}

type EPA struct {
	Root *EPANode
}

type PenetrationInfo struct {
	Depth         float64
	Normal        zero.Vector2D
	ContactShapeA zero.Vector2D
	ContactShapeB zero.Vector2D
}

func GetContactPoints(epa EPA, edge Edge) (zero.Vector2D, zero.Vector2D) {
	ab := edge.Parent.Next.Supp.Pnt.Sub(edge.Parent.Supp.Pnt)
	if ab.IsZero(zero.Tolerance) {
		return edge.Parent.Supp.ShapeAPnt, edge.Parent.Supp.ShapeBPnt
	}

	ao := edge.Parent.Supp.Pnt.Inverse()
	along := ao.Dot(ab) / ab.Dot(ab)

	if along < 0.0 {
		along = 0.0
	} else if along > 1.0 {
		along = 1.0
	}

	contacta := edge.Parent.Next.Supp.ShapeAPnt.Sub(edge.Parent.Supp.ShapeAPnt).Mult(along).Add(edge.Parent.Supp.ShapeAPnt)
	contactb := edge.Parent.Next.Supp.ShapeBPnt.Sub(edge.Parent.Supp.ShapeBPnt).Mult(along).Add(edge.Parent.Supp.ShapeBPnt)
	return contacta, contactb
}

func InitEPA(pnta SupportComponents, pntb SupportComponents, pntc SupportComponents) EPA {
	a := &EPANode{pnta, nil}
	b := &EPANode{pntb, nil}
	c := &EPANode{pntc, nil}

	a.Next = b
	b.Next = c
	c.Next = a

	return EPA{a}
}

func ClosestEdge(epa EPA) Edge {
	edge := Edge{}
	edge.Distance = math.MaxFloat64

	node := epa.Root
	for {
		e := node.Next.Supp.Pnt.Sub(node.Supp.Pnt)
		n := node.Supp.Pnt.Mult(e.Dot(e)).Sub(e.Mult(node.Supp.Pnt.Dot(e))).Normalize()
		distance := n.Dot(node.Supp.Pnt)

		if distance < edge.Distance {
			edge.Distance = distance
			edge.Normal = n
			edge.Parent = node
		}

		if node.Next == epa.Root {
			break
		}
		node = node.Next
	}
	return edge
}

func GetPenetrationInfo(shapea Collidable, shapeb Collidable, pnta SupportComponents, pntb SupportComponents, pntc SupportComponents) PenetrationInfo {
	epa := InitEPA(pnta, pntb, pntc)

	// for i := 0; i < MaxIter; i++ {
	for {
		edge := ClosestEdge(epa)

		support := Support(shapea, shapeb, edge.Normal)
		depth := support.Pnt.Dot(edge.Normal)
		if depth-edge.Distance < zero.Tolerance {
			contacta, contactb := GetContactPoints(epa, edge)
			return PenetrationInfo{depth, edge.Normal, contacta, contactb}
		}

		edge.Parent.Insert(support)
	}
	return PenetrationInfo{}
}
