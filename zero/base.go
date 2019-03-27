package zero

import (
	"math"
)

const (
	Tolerance = 1.0e-10
)

type Vector2D struct {
	X float64
	Y float64
}

type Tf2D struct {
	Translation Vector2D
	Rotation    float64
}

func (vec Vector2D) IsZero(tolerance float64) bool {
	return math.Abs(vec.X) < tolerance && math.Abs(vec.Y) < tolerance
}

func (vec Vector2D) Inverse() Vector2D {
	return Vector2D{-vec.X, -vec.Y}
}

func (tf Tf2D) Inverse() Tf2D {
	return Tf2D{Vector2D{-tf.Translation.X, -tf.Translation.Y}, -tf.Rotation}
}

func (v1 Vector2D) Mult(mult float64) Vector2D {
	return Vector2D{v1.X * mult, v1.Y * mult}
}

func (v1 Vector2D) Div(div float64) Vector2D {
	return Vector2D{v1.X / div, v1.Y / div}
}

func (v1 Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{v1.X + v2.X, v1.Y + v2.Y}
}

func (v1 Vector2D) Sub(v2 Vector2D) Vector2D {
	return Vector2D{v1.X - v2.X, v1.Y - v2.Y}
}

func (v1 Vector2D) Magnitude() float64 {
	return math.Sqrt(v1.Dot(v1))
}

func (v1 Vector2D) Normalize() Vector2D {
	imag := 1.0 / v1.Magnitude()
	return Vector2D{v1.X * imag, v1.Y * imag}
}

func (vec Vector2D) Rotate(rot float64) Vector2D {
	cost := math.Cos(rot)
	sint := math.Sin(rot)

	return Vector2D{cost*vec.X - sint*vec.Y,
		sint*vec.X + cost*vec.Y}
}

func (tf Tf2D) Transform(vec Vector2D) Vector2D {
	cost := math.Cos(tf.Rotation)
	sint := math.Sin(tf.Rotation)

	return Vector2D{cost*vec.X - sint*vec.Y + tf.Translation.X,
		sint*vec.X + cost*vec.Y + tf.Translation.Y}
}

func (v1 Vector2D) Dot(v2 Vector2D) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func (v1 Vector2D) Cross(z float64) Vector2D {
	return Vector2D{v1.Y * z, -v1.X * z}
}
