package card

import (
	"fmt"
	"math"
	"sync"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/dustin/go-humanize"
	"github.com/vron/gioe/cleantheme"
)

// A Card displays a card as in e.g. a BI application, i.e. a title with a value (and a min and max).
type Card struct {
	Theme *cleantheme.Theme
	Title string

	mu       sync.Mutex
	val      float64
	min, max float64
}

// New creates a new Card.
func New(title string, theme *cleantheme.Theme) *Card {
	return &Card{
		Theme: theme,
		Title: title,
		min:   math.MaxFloat64,
		max:   -math.MaxFloat64,
	}
}

// Layout draws the Card as a Gio widget
func (c *Card) Layout(gtx *layout.Context) {
	title, value, limits := c.getStrings()
	layout.Inset{Top: c.Theme.Spacing}.Layout(gtx, func() {
		layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func() {
				rr := float32(gtx.Px(c.Theme.Radius))
				clip.Rect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(gtx.Constraints.Width.Min),
						Y: float32(gtx.Constraints.Height.Min),
					}},
					NE: rr, NW: rr, SE: rr, SW: rr,
				}.Op(gtx.Ops).Add(gtx.Ops)
				cleantheme.Fill(gtx, c.Theme.Color.ShadedBackground)
			}),
			layout.Stacked(func() {
				layout.Inset{
					Top: unit.Dp(10), Bottom: unit.Dp(10),
					Left: unit.Dp(14), Right: unit.Dp(14),
				}.Layout(gtx, func() {
					gtx.Constraints.Width.Min = 440 // TODO: Base this on font metrics, but how?
					layout.Flex{
						Axis:      layout.Vertical,
						Spacing:   layout.SpaceAround,
						Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func() {
							paint.ColorOp{Color: c.Theme.Color.Secondary}.Add(gtx.Ops)
							widget.Label{Alignment: text.Middle}.Layout(gtx, c.Theme.Shaper, c.Theme.Font, c.Theme.TextSize, title)
						}),
						layout.Rigid(func() {
							paint.ColorOp{Color: c.Theme.Color.Primary}.Add(gtx.Ops)
							widget.Label{Alignment: text.Middle}.Layout(gtx, c.Theme.Shaper, c.Theme.Font, c.Theme.TextSize.Scale(3.5), value)
						}),
						layout.Rigid(func() {
							paint.ColorOp{Color: c.Theme.Color.Primary}.Add(gtx.Ops)
							widget.Label{Alignment: text.Middle}.Layout(gtx, c.Theme.Shaper, c.Theme.Font, c.Theme.TextSize, limits)
						}),
					)
				})
			}),
		)
	})
}

func (c *Card) getStrings() (title, value, limits string) {
	title = c.Title
	val, min, max := c.values()
	value = toPostfix(val)
	limits = "(" + toPostfix(min) + "   —   " + toPostfix(max) + ")"
	return
}

func toPostfix(v float64) string {
	va, s := humanize.ComputeSI(v)
	return fmt.Sprintf("%.1f%v", va, s)
}
