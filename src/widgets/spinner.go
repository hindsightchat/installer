package widgets

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type Spinner struct {
	widget.BaseWidget
	angle float64
	size  float32
	color color.RGBA
	stop  chan bool
}

func NewSpinner(size float32, c color.RGBA) *Spinner {
	s := &Spinner{
		size:  size,
		color: c,
		stop:  make(chan bool),
	}
	s.ExtendBaseWidget(s)
	return s
}

func (s *Spinner) Start() {
	go func() {
		tick := time.NewTicker(16 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case <-s.stop:
				return
			case <-tick.C:
				s.angle += 5
				if s.angle >= 360 {
					s.angle = 0
				}
				s.Refresh()
			}
		}
	}()
}

func (s *Spinner) Stop() {
	select {
	case s.stop <- true:
	default:
	}
}

func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	return &spinnerRenderer{s: s}
}

type spinnerRenderer struct {
	s *Spinner
}

func (r *spinnerRenderer) Layout(_ fyne.Size) {}

func (r *spinnerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.s.size, r.s.size)
}

func (r *spinnerRenderer) Refresh() {
	canvas.Refresh(r.s)
}

func (r *spinnerRenderer) Objects() []fyne.CanvasObject {
	size := r.s.size
	cx, cy := size/2, size/2
	radius := size/2 - 4
	stroke := float32(3)

	var objs []fyne.CanvasObject
	segments := 12
	arc := 270.0

	for i := 0; i < segments; i++ {
		prog := float64(i) / float64(segments)
		alpha := uint8(255 * prog)

		start := r.s.angle + (arc * prog)
		end := start + (arc / float64(segments))

		startRad := start * math.Pi / 180
		endRad := end * math.Pi / 180

		x1 := cx + float32(math.Cos(startRad))*radius
		y1 := cy + float32(math.Sin(startRad))*radius
		x2 := cx + float32(math.Cos(endRad))*radius
		y2 := cy + float32(math.Sin(endRad))*radius

		line := canvas.NewLine(color.RGBA{
			R: r.s.color.R,
			G: r.s.color.G,
			B: r.s.color.B,
			A: alpha,
		})
		line.StrokeWidth = stroke
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)
		objs = append(objs, line)
	}

	return objs
}

func (r *spinnerRenderer) Destroy() {
	r.s.Stop()
}
