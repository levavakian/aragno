package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"aragno/gjk"
	"aragno/zero"

	"fmt"
	"time"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1600, 800),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	shapes := []([]pixel.Vec){}
	collisions := []bool{}
	accumPoints := []pixel.Vec{}
	for !win.Closed() {
		imd := imdraw.New(nil)

		if win.Pressed(pixelgl.KeyLeftShift) {
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				accumPoints = append(accumPoints, win.MousePosition())
			}
		} else if len(accumPoints) > 0 {
			shapes = append(shapes, accumPoints)
			accumPoints = []pixel.Vec{}
		}

		mpos := win.MousePosition()
		poly := gjk.Polygon{}

		poly.Pnts = append(poly.Pnts, zero.Vector2D{mpos.X - 50, mpos.Y})
		poly.Pnts = append(poly.Pnts, zero.Vector2D{mpos.X, mpos.Y + 50})
		poly.Pnts = append(poly.Pnts, zero.Vector2D{mpos.X + 50, mpos.Y})
		poly.Pnts = append(poly.Pnts, zero.Vector2D{mpos.X, mpos.Y - 50})

		imd.Color = pixel.RGB(1, 0, 0)
		for _, pnt := range poly.Pnts {
			imd.Push(pixel.V(pnt.X, pnt.Y))
		}
		imd.Polygon(5)


		collisions = []bool{}
		for _, shape := range shapes {
			spoly := gjk.Polygon{}
			for _, vec := range shape {
				spoly.Pnts = append(spoly.Pnts, zero.Vector2D{vec.X, vec.Y})
			}

			report := gjk.CheckCollision(poly, spoly)
			collisions = append(collisions, report.Collision)

			if report.Collision {
				imd.Color = pixel.RGB(.2, .2, .8)
				imd.Push(pixel.V(report.Penetration.ContactShapeB.X, report.Penetration.ContactShapeB.Y))
				imd.Circle(10, 0)

				imd.Color = colornames.Blueviolet
				imd.EndShape = imdraw.RoundEndShape
				imd.Push(pixel.V(report.Penetration.ContactShapeB.X, report.Penetration.ContactShapeB.Y))
				imd.Push(pixel.V(report.Penetration.ContactShapeB.X+report.Penetration.Normal.X*report.Penetration.Depth,
					report.Penetration.ContactShapeB.Y+report.Penetration.Normal.Y*report.Penetration.Depth))
				imd.Line(5)
			}

			imd.Color = pixel.RGB(0, 1, 0)
			imd.Push(pixel.V(report.ClosestPointShapeB.X, report.ClosestPointShapeB.Y))
			imd.Circle(10, 0)
		}

		for idx, shape := range shapes {
			if len(collisions) > 0 && collisions[idx] {
				imd.Color = pixel.RGB(1, 0, 0)
			} else {
				imd.Color = pixel.RGB(0, 0, 1)
			}
			for _, vec := range shape {
				imd.Push(vec)
			}
			imd.Polygon(5)
		}

		imd.Color = pixel.RGB(0, 1, 0)
		imd.EndShape = imdraw.RoundEndShape
		for _, vec := range accumPoints {
			imd.Push(vec)
		}
		imd.Line(5)

		win.Clear(colornames.Black)
		imd.Draw(win)
		win.Update()
		time.Sleep(time.Millisecond * 10)
	}
}

func main() {
	fmt.Println("Collision Example Exe")
	pixelgl.Run(run)
}
