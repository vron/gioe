package cleantheme

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

// A cleantheme represents styling for a very simple and clean, grayscale UI.
type Theme struct {
	Shaper text.Shaper
	Color  struct {
		Primary          color.RGBA
		Secondary        color.RGBA
		Background       color.RGBA
		ShadedBackground color.RGBA
	}
	Radius   unit.Value
	Spacing  unit.Value
	TextSize unit.Value
	Font     text.Font
}

func New() *Theme {
	t := &Theme{
		Shaper:   font.Default(),
		Radius:   unit.Dp(8),
		Spacing:  unit.Dp(20),
		TextSize: unit.Sp(16),
	}
	t.Color.Primary = rgb(0xD4D4D4)
	t.Color.Secondary = rgb(0x595959)
	t.Color.Background = rgb(0x000000)
	t.Color.ShadedBackground = rgb(0x0F0F0F)
	return t
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func Fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
