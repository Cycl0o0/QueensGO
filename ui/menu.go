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
	Selected int
}

func NewMenu() *Menu {
	return &Menu{Selected: -1}
}

func (m *Menu) Update(screenW, screenH int) int {
	m.Selected = -1

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		const bW, bH, startY, gap = 240, 50, 210, 15

		for i := range Difficulties {
			bx := (screenW - bW) / 2
			by := startY + i*(bH+gap)
			if image.Pt(mx, my).In(image.Rect(bx, by, bx+bW, by+bH)) {
				m.Selected = i
				return m.Selected
			}
		}
	}
	return -1
}

func (m *Menu) Draw(screen *ebiten.Image, face text.Face) {
	screenW := screen.Bounds().Dx()

	// Large title
	title := "Queens"
	tw, _ := text.Measure(title, face, 0)
	{
		const scale = 1.9
		opts := &text.DrawOptions{}
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(float64(screenW)/2-tw*scale/2, 42)
		opts.ColorScale.ScaleWithColor(color.RGBA{0x33, 0x33, 0x33, 0xFF})
		text.Draw(screen, title, face, opts)
	}

	// Tagline
	tagline := "One queen per row, column & color. No adjacency."
	tlw, _ := text.Measure(tagline, face, 0)
	{
		const scale = 0.78
		opts := &text.DrawOptions{}
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(float64(screenW)/2-tlw*scale/2, 122)
		opts.ColorScale.ScaleWithColor(color.RGBA{0x77, 0x77, 0x77, 0xFF})
		text.Draw(screen, tagline, face, opts)
	}

	// Subtitle
	subtitle := "Select Difficulty"
	sw, _ := text.Measure(subtitle, face, 0)
	{
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(screenW)/2-sw/2, 165)
		opts.ColorScale.ScaleWithColor(color.RGBA{0x55, 0x55, 0x55, 0xFF})
		text.Draw(screen, subtitle, face, opts)
	}

	// Difficulty buttons
	const bW, bH, startY, gap = float32(240), float32(50), float32(210), float32(15)
	mx, my := ebiten.CursorPosition()

	for i, diff := range Difficulties {
		bx := (float32(screenW) - bW) / 2
		by := startY + float32(i)*(bH+gap)

		hovered := float32(mx) >= bx && float32(mx) <= bx+bW &&
			float32(my) >= by && float32(my) <= by+bH

		bgColor := RegionColors[i]
		if hovered {
			bgColor = RegionColorsDark[i]
		}

		vector.DrawFilledRect(screen, bx, by, bW, bH, bgColor, true)
		vector.StrokeRect(screen, bx, by, bW, bH, 2, color.RGBA{0x44, 0x44, 0x44, 0xFF}, true)

		label := diff.Name
		lw, lh := text.Measure(label, face, 0)
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(bx)+float64(bW)/2-lw/2, float64(by)+float64(bH)/2-lh/2)
		opts.ColorScale.ScaleWithColor(color.RGBA{0x22, 0x22, 0x22, 0xFF})
		text.Draw(screen, label, face, opts)
	}

	// Click tip at bottom
	tip := "Click: mark X  |  Click again: place Queen  |  Click again: clear"
	tipW, _ := text.Measure(tip, face, 0)
	{
		const scale = 0.70
		opts := &text.DrawOptions{}
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(float64(screenW)/2-tipW*scale/2, 590)
		opts.ColorScale.ScaleWithColor(color.RGBA{0x99, 0x99, 0x99, 0xFF})
		text.Draw(screen, tip, face, opts)
	}
}
