// Command example shows a simple example of using dui to display debug info for a 2D physics simulation.
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/vron/gioe/dui"
	"github.com/vron/v2d"
)

var (
	fNo  int
	fRes int
)

func init() {
	flag.IntVar(&fNo, "no", 100, "number of circles to draw")
	flag.IntVar(&fRes, "res", 16, "number of triangles per circle")
}

func main() {
	fmt.Println("Running an example of debugging of a 2D simulation engine")
	ui := dui.New("DUI example")
	ui.SetBounds(v2d.R(v2d.V(-20, -20), v2d.V(20, 20)))
	go runDummySimulation(ui)

	if err := ui.Draw(); err != nil {
		log.Fatal(err)
	}
}

func runDummySimulation(ui *dui.DUI) {
	// simply perform a random walk that resets to the center if the center
	// point of a circle is out of bounds.
	rnd := rand.New(rand.NewSource(0))
	shapes := makeCircles(fNo, fRes)
	dirs := make([]float64, len(shapes))
	for i := range dirs {
		dirs[i] = rnd.Float64() * 2 * math.Pi
	}
	tick := time.NewTicker(time.Second / 120)
	c := tick.C
	for {

		accx, accy := 0.0, 0.0
		accReset, accDist := 0.0, 0.0

		for i := range shapes {
			s := float32(0.1)
			d := rnd.Float32() * 0.02
			dirs[i] += (rnd.Float64() - 0.5) * float64(s)
			dx, dy := d*float32(math.Cos(dirs[i])), d*float32(math.Sin(dirs[i]))
			accDist += float64(d)
			rx, ry := float32(0), float32(0)
			if shapes[i][0][0].X < -20 || shapes[i][0][0].X > 20 ||
				shapes[i][0][0].Y < -20 || shapes[i][0][0].Y > 20 {
				rx, ry = shapes[i][0][0].X, shapes[i][0][0].Y
				accReset++
			}
			for j := range shapes[i] {
				shapes[i][j][0].X += -rx + dx
				shapes[i][j][0].Y += -ry + dy
				shapes[i][j][1].X += -rx + dx
				shapes[i][j][1].Y += -ry + dy
				shapes[i][j][2].X += -rx + dx
				shapes[i][j][2].Y += -ry + dy
			}
			accx += float64(shapes[i][0][0].X)
			accy += float64(shapes[i][0][0].Y)
		}

		ui.SetShapes(shapes)
		ui.SetMeasure("Mean x", accx/float64(len(shapes)))
		ui.SetMeasure("Mean y", accy/float64(len(shapes)))
		ui.SetMeasure("No reset", float64(accReset))
		ui.SetMeasure("Step distance", accDist)

		<-c // to not move them to often...
	}
}

func makeCircles(no, res int) [][][3]v2d.Vec {
	outline := make([][3]v2d.Vec, res)
	for i := range outline {
		a1 := 2 * float64(i) * math.Pi / float64(res)
		a2 := 2 * float64(i+1) * math.Pi / float64(res)
		outline[i][1] = v2d.Vec{float32(math.Cos(a1)), float32(math.Sin(a1))}
		outline[i][2] = v2d.Vec{float32(math.Cos(a2)), float32(math.Sin(a2))}
	}
	shapes := make([][][3]v2d.Vec, no)
	for i := range shapes {
		shapes[i] = make([][3]v2d.Vec, res)
		copy(shapes[i], outline)
	}

	return shapes
}
