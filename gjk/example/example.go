package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/imdraw"

	"aragno/zero"
	"aragno/gjk"

	"fmt"
	"time"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1600, 1000),
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

		st := time.Now()

		mpos := win.MousePosition()
		poly := gjk.Polygon{}

		poly.Pnts = append(poly.Pnts, zero.Vector2D{mpos.X - 1, mpos.Y - 1})
		poly.Pnts = append(poly.Pnts, zero.Vector2D{mpos.X + 1, mpos.Y - 1})
		poly.Pnts = append(poly.Pnts, zero.Vector2D{mpos.X + 1, mpos.Y + 1})
		poly.Pnts = append(poly.Pnts, zero.Vector2D{mpos.X - 1, mpos.Y + 1})

		collisions = []bool{}
		for _, shape := range shapes {
			spoly := gjk.Polygon{}
			for _, vec := range shape {
				spoly.Pnts = append(spoly.Pnts, zero.Vector2D{vec.X, vec.Y})
			}

			report := gjk.CheckCollision(poly, spoly)
			collisions = append(collisions, report.Collision)

			imd.Color = pixel.RGB(0, 1, 0)
			imd.Push(pixel.V(report.ClosestPointShapeB.X, report.ClosestPointShapeB.Y))
			imd.Circle(10, 0)
			imd.Color = pixel.RGB(.5, .5, 0)
			imd.Push(pixel.V(report.ClosestPointShapeA.X, report.ClosestPointShapeA.Y))
			imd.Circle(10, 0)
		}

		fmt.Println(time.Since(st))

		for idx, shape := range shapes {
			if len(collisions) > 0 && collisions[idx] {
				imd.Color = pixel.RGB(1, 0, 0)
			} else {
				imd.Color = pixel.RGB(0, 0, 1)
			}
			for _, vec := range shape {
				imd.Push(vec)
			}
			imd.Polygon(0)
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