package zero

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Rectangle struct {
	P0 Point `json:"p0"`
	P1 Point `json:"p1"`
	P2 Point `json:"p2"`
}
