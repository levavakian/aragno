package zero

import (
	"math"
)

var TwoPi = 2 * math.Pi

// NormalizeBipolar clamps angle in range [-pi, pi)
func NormalizeBipolar(angle float64) float64 {
	angle = math.Mod(angle+math.Pi, TwoPi)
	if angle < 0 {
		angle = angle + TwoPi
	}
	return angle - math.Pi
}
