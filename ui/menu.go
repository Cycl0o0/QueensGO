package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Difficulty struct {
	Name string
	Size int
}

var Difficulties = []Difficulty{
	{"Easy", 5},
	{"Normal", 6},
	{"Hard", 7},
	{"Expert", 8},
	{"Extreme", 9},
}

// Menu renders the difficulty selection screen.
type Menu struct {
	Selected int // -1 if nothing selected yet
}

func NewMenu() *Menu {
	return &Menu{Selected: -1}
}

func (m *Menu) Update(screenW, screenH int) int {
	m.Selected = -1

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		btnW := 240
		btnH := 50
		startY := 200
		gap := 15

		for i := range Difficulties {
			bx := (screenW - btnW) / 2
			by := startY + i*(btnH+gap)
			rect := image.Rect(bx, by, bx+btnW, by+btnH)
			if image.Pt(mx, my).In(rect) {
				m.Selected = i
				return m.Selected
			}
		}
	}
	return -1
}

func (m *Menu) Draw(screen *ebiten.Image, face text.Face) {
	screenW := screen.Bounds().Dx()

	// Title
	title := "Queens"
	tw, th := text.Measure(title, face, 0)
	_ = th
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(float64(screenW)/2-tw/2, 80)
	opts.ColorScale.ScaleWithColor(color.RGBA{0x33, 0x33, 0x33, 0xFF})
	text.Draw(screen, title, face, opts)

	subtitle := "Select Difficulty"
	sw, _ := text.Measure(subtitle, face, 0)
	opts2 := &text.DrawOptions{}
	opts2.GeoM.Translate(float64(screenW)/2-sw/2, 140)
	opts2.ColorScale.ScaleWithColor(color.RGBA{0x66, 0x66, 0x66, 0xFF})
	text.Draw(screen, subtitle, face, opts2)

	// Buttons
	btnW := float32(240)
	btnH := float32(50)
	startY := float32(200)
	gap := float32(15)

	mx, my := ebiten.CursorPosition()

	for i, diff := range Difficulties {
		bx := (float32(screenW) - btnW) / 2
		by := startY + float32(i)*(btnH+gap)

		// Check hover
		hovered := float32(mx) >= bx && float32(mx) <= bx+btnW &&
			float32(my) >= by && float32(my) <= by+btnH

		bgColor := RegionColors[i]
		if hovered {
			bgColor = RegionColorsDark[i]
		}

		vector.DrawFilledRect(screen, bx, by, btnW, btnH, bgColor, true)
		vector.StrokeRect(screen, bx, by, btnW, btnH, 2, color.RGBA{0x44, 0x44, 0x44, 0xFF}, true)

		// Button label
		label := diff.Name
		lw, lh := text.Measure(label, face, 0)
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(bx)+float64(btnW)/2-lw/2, float64(by)+float64(btnH)/2-lh/2)
		opts.ColorScale.ScaleWithColor(color.RGBA{0x22, 0x22, 0x22, 0xFF})
		text.Draw(screen, label, face, opts)
	}
}
