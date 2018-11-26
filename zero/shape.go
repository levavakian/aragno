package zero

type Pose struct {
	X     float64
	Y     float64
	Theta float64
}

func (p *Pose) Sum(p2 Pose) {
	p.X = p.X + p2.X
	p.Y = p.Y + p2.Y
	p.Theta = NormalizeBipolar(p.Theta + p2.Theta)
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Rectangle struct {
	P0 Point `json:"p0"`
	P1 Point `json:"p1"`
	P2 Point `json:"p2"`
}
