package polyscene

import (
	"image/color"
	"sync"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/vron/gioe/cleantheme"
	"github.com/vron/v2d"
)

// A Polyscen shows a slice of triangulated polygons in a 2D environment.
type Polyscene struct {
	Theme *cleantheme.Theme

	mu  sync.Mutex
	bnd v2d.Rect
}

func New(theme *cleantheme.Theme) Polyscene {
	return Polyscene{Theme: theme}
}

func (c Polyscene) Layout(gtx *layout.Context, shapes [][][3]v2d.Vec) {
	layout.Inset{Top: c.Theme.Spacing, Bottom: c.Theme.Spacing}.Layout(gtx, func() {
		rr := float32(gtx.Px(c.Theme.Radius))
		clip.Rect{
			Rect: f32.Rectangle{Max: f32.Point{
				X: float32(gtx.Constraints.Width.Min),
				Y: float32(gtx.Constraints.Height.Min),
			}},
			NE: rr, NW: rr, SE: rr, SW: rr,
		}.Op(gtx.Ops).Add(gtx.Ops)
		cleantheme.Fill(gtx, c.Theme.Color.ShadedBackground)

		for _, shape := range shapes {
			drawShape(gtx, shape, c.Bounds(), c.Theme.Color.Secondary)
		}
	})
}

func drawShape(gtx *layout.Context, s [][3]v2d.Vec, bnd v2d.Rect, col color.RGBA) {
	h := float32(gtx.Constraints.Height.Max)
	w := float32(gtx.Constraints.Width.Max)
	scale := h / bnd.H()
	if w/bnd.W() < scale {
		scale = w / bnd.W()
	}
	tx := func(x float32) float32 {
		x -= 0.5 * (bnd.Max.X + bnd.Min.X)
		x *= scale
		x += 0.5 * w
		return x
	}
	ty := func(y float32) float32 {
		y -= 0.5 * (bnd.Max.Y + bnd.Min.Y)
		y *= scale
		y += 0.5 * h
		return y
	}

	// draw each of the triangles
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	for _, tri := range s {
		var stack op.StackOp
		stack.Push(gtx.Ops)
		p := clip.Path{}
		p.Begin(gtx.Ops)
		p.Move(f32.Point{
			X: tx(tri[0].X),
			Y: ty(tri[0].Y),
		})
		p.Line(f32.Point{
			X: tx(tri[1].X) - tx(tri[0].X),
			Y: ty(tri[1].Y) - ty(tri[0].Y),
		})
		p.Line(f32.Point{
			X: tx(tri[2].X) - tx(tri[1].X),
			Y: ty(tri[2].Y) - ty(tri[1].Y),
		})
		p.Line(f32.Point{
			X: tx(tri[0].X) - tx(tri[2].X),
			Y: ty(tri[0].Y) - ty(tri[2].Y),
		})
		p.End().Add(gtx.Ops)
		paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(w), Y: float32(h)}}}.Add(gtx.Ops)
		stack.Pop()

	}
}
