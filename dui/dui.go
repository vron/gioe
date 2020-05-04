// Package dui implements a UI usefull for debugging a 2D physics engine
package dui

import (
	"sync"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/vron/gioe/cleantheme"
	"github.com/vron/gioe/cleantheme/card"
	"github.com/vron/gioe/cleantheme/polyscene"
	"github.com/vron/v2d"
)

type DUI struct {
	w     *app.Window
	theme *cleantheme.Theme
	gtx   *layout.Context

	mu            sync.Mutex
	measuresOrder []string
	measures      map[string]*card.Card
	bounds        v2d.Rect
	shapes        [][][3]v2d.Vec
}

func New(title string) *DUI {
	gofont.Register()
	ui := &DUI{
		w:             app.NewWindow(app.Title(title)),
		theme:         cleantheme.New(),
		measuresOrder: []string{},
		measures:      make(map[string]*card.Card),
		shapes:        make([][][3]v2d.Vec, 0),
	}
	ui.gtx = layout.NewContext(ui.w.Queue())
	return ui
}

// Draw blocks and starts drawing for ever.
func (ui *DUI) Draw() (err error) {
	go func() {
		err = ui.loopDraw()
	}()
	app.Main()
	return
}

// SetShapes uses copies the data provided and will use that for displaying. Not that
// the argument will not be retained after the call returns. It is thus safe for the caller
// to re-use the same buffer for multiple frames.
func (ui *DUI) SetShapes(shapes [][][3]v2d.Vec) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	if len(shapes) == len(ui.shapes) {
		for i := range shapes {
			copy(ui.shapes[i], shapes[i])
		}
		return
	}
	ui.shapes = make([][][3]v2d.Vec, len(shapes))
	for i := range shapes {
		ui.shapes[i] = make([][3]v2d.Vec, len(shapes[i]))
		copy(ui.shapes[i], shapes[i])
	}
}

// SetBounds sets the simulation region that should be shown.
func (ui *DUI) SetBounds(bnd v2d.Rect) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	ui.bounds = bnd
}

// SetMeasure registers a measure (if the string measure is new) or updates
// an previously registered one that should be shown to the user.
func (ui *DUI) SetMeasure(measure string, value float64) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	m, ok := ui.measures[measure]
	if ok {
		m.UpdateValue(value)
		return
	}
	m = card.New(measure, ui.theme)
	ui.measures[measure] = m
	ui.measuresOrder = append(ui.measuresOrder, measure)
}

func (ui *DUI) loopDraw() error {
	gtx := ui.gtx
	scene := polyscene.New(ui.theme)

	for {
		e := <-ui.w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)
			paint.ColorOp{Color: ui.theme.Color.Background}.Add(gtx.Ops)
			paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(e.Size.X), Y: float32(e.Size.Y)}}}.Add(gtx.Ops)

			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func() {

						// show all the measures that have been registered
						ui.mu.Lock()
						children := make([]layout.FlexChild, 0, 10)
						for _, m := range ui.measuresOrder {
							card := ui.measures[m]
							children = append(children, layout.Rigid(func() {
								card.Layout(gtx)
							}))
						}
						ui.mu.Unlock()
						layout.Flex{
							Axis:    layout.Vertical,
							Spacing: layout.SpaceEnd}.Layout(gtx, children...)
					})
				}),
				layout.Flexed(1.0, func() {
					layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceAround}.Layout(gtx,
						layout.Flexed(1.0, func() {
							layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func() {
								ui.mu.Lock()
								defer ui.mu.Unlock()
								scene.SetBounds(ui.bounds)
								scene.Layout(gtx, ui.shapes)
							})
						}))
				}),
			)

			op.InvalidateOp{}.Add(gtx.Ops)
			e.Frame(gtx.Ops)
		}
	}
}
