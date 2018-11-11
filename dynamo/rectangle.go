package dynamo

import "aragno/zero"

type Rectangle struct {
	width float64
	height float64
	momentOfInertia float64
	id string
	body
}

func NewRectangleBody(id string, width float64, height float64, mass float64, p zero.Pose ) *Body{
	return &Body{[]shape{NewRectangle(id, width, height, mass, p)}, []Joint{}}
}

func NewRectangle(id string, width float64, height float64, mass float64, p zero.Pose) *Rectangle{
	moment := mass * ( width * width + height * height) /12
	return &Rectangle {width, height,moment, id, body {p, zero.Pose{0,0,0},zero.Pose{0,0,0}, zero.Pose{0,0,0}, mass}}
}


func (r *Rectangle) GetId() string{
	return r.id
}

func (r *Rectangle) GetBody() *body{
	return &r.body
}

func (r *Rectangle) MomentOfInertia() float64{
    return r.body.Mass * ( r.width * r.width + r.height * r.height) /12
}

func (r *Rectangle) ComputeForceAndTorque(p zero.Pose) {
	r.body.Force = p
	// r is the 'arm vector' that goes from the center of mass to the point of force application
        var a zero.Pose = zero.Pose{r.width / 2, r.height / 2, 0}
	r.body.Force.Theta = r.body.Position.X * a.Y - r.body.Position.Y * a.X
	r.body.Force.Sum(ComputeGravity(r.body))
}
