package dynamo

import "fmt"
import "aragno/zero"

type Point struct {
	x float64
	y float64
}

type body struct {
	Position zero.Pose
	LinearV  zero.Pose
	AngularV zero.Pose
	Force    zero.Pose
	Mass     float64
}

type Circle struct {
	Center          Point
	Radius          float64
	momentOfInertia float64
	body
}

type shape interface {
	GetId() string
	MomentOfInertia() float64
	ComputeForceAndTorque(v zero.Pose)
	UpdateBodyProps(a zero.Pose, dt float64)
	Acceleration(moment float64) zero.Pose
	SetVelocity(setFunc)
	SetPosition(setFunc)
	Print()
}

type Joint struct {
	shapes []string // the IDs of the shapes this is joining
}

type Body struct {
	Shapes []shape
	Joints []Joint
}

type FTFunc func(zero.Pose)

type setFunc func(x float64, y float64, theta float64)

func (b *body) SetVelocity(setInfo setFunc) {
	setInfo(b.LinearV.X, b.LinearV.Y, b.LinearV.Theta)
}

func (b *body) SetPosition(setInfo setFunc) {
	setInfo(b.Position.X, b.Position.Y, b.Position.Theta)
}

func (b *body) Print() {
	fmt.Printf("X: %v Y: %v\n", b.Position.X, b.Position.Y)
}

func (b *body) UpdateBodyProps(a zero.Pose, dt float64) {
	b.LinearV.X += a.X * dt
	b.LinearV.Y += a.Y * dt
	b.Position.X += b.LinearV.X * dt
	b.Position.Y += b.LinearV.Y * dt
	b.LinearV.Theta += a.Theta * dt
	b.Position.Theta += b.LinearV.Theta * dt
}

func (b *body) Acceleration(momentOfInertia float64) zero.Pose {
	return zero.Pose{b.Force.X / b.Mass, b.Force.Y / b.Mass, b.Force.Theta / momentOfInertia}
}

func UpdateBody(b *Body, v zero.Pose, dt float64) {
	for _, sh := range b.Shapes {
		sh.ComputeForceAndTorque(v)
		sh.UpdateBodyProps(sh.Acceleration(sh.MomentOfInertia()), dt)
	}
}

func ComputeGravity(particle body) zero.Pose {
	var gravity float64 = -9.81
	return zero.Pose{0, particle.Mass * gravity, 0}
}
