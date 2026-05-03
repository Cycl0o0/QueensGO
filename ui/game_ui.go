package ui

import (
	"fmt"
	"image"
	"image/color"

	"github.com/Cycl0o0/QueensGO/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	boardPadding = 40
	topBarHeight = 60
	botBarHeight = 60
)

// GameUI handles rendering and input for the puzzle board.
type GameUI struct {
	Board      *game.Board
	cellSize   int
	offsetX    int
	offsetY    int
	justClicked bool
	Won        bool
	WantMenu   bool
	WantNew    bool
}

func NewGameUI(board *game.Board) *GameUI {
	return &GameUI{Board: board}
}

func (g *GameUI) layout(screenW, screenH int) {
	available := screenH - topBarHeight - botBarHeight - boardPadding*2
	availW := screenW - boardPadding*2
	if availW < available {
		available = availW
	}
	g.cellSize = available / g.Board.Size
	totalSize := g.cellSize * g.Board.Size
	g.offsetX = (screenW - totalSize) / 2
	g.offsetY = topBarHeight + (screenH-topBarHeight-botBarHeight-totalSize)/2
}

func (g *GameUI) Update(screenW, screenH int) {
	g.layout(screenW, screenH)
	g.WantMenu = false
	g.WantNew = false

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.justClicked = false
		return
	}

	if g.justClicked {
		return
	}
	g.justClicked = true

	mx, my := ebiten.CursorPosition()

	// Check bottom buttons
	btnW := 140
	btnH := 36
	// "New Game" button
	newX := screenW/2 - btnW - 20
	newY := screenH - botBarHeight + 12
	if image.Pt(mx, my).In(image.Rect(newX, newY, newX+btnW, newY+btnH)) {
		g.WantNew = true
		return
	}
	// "Menu" button
	menuX := screenW/2 + 20
	menuY := newY
	if image.Pt(mx, my).In(image.Rect(menuX, menuY, menuX+btnW, menuY+btnH)) {
		g.WantMenu = true
		return
	}

	if g.Won {
		return
	}

	// Check board clicks
	col := (mx - g.offsetX) / g.cellSize
	row := (my - g.offsetY) / g.cellSize
	if row >= 0 && row < g.Board.Size && col >= 0 && col < g.Board.Size {
		g.Board.CycleMark(row, col)
		if g.Board.IsWon() {
			g.Won = true
		}
	}
}

func (g *GameUI) Draw(screen *ebiten.Image, face text.Face) {
	screenW := screen.Bounds().Dx()
	screenH := screen.Bounds().Dy()
	g.layout(screenW, screenH)

	conflicts := g.Board.Conflicts()
	cs := float32(g.cellSize)

	// Draw cells
	for r := 0; r < g.Board.Size; r++ {
		for c := 0; c < g.Board.Size; c++ {
			x := float32(g.offsetX + c*g.cellSize)
			y := float32(g.offsetY + r*g.cellSize)

			reg := g.Board.Regions[r][c]
			bgColor := RegionColors[reg%len(RegionColors)]

			vector.DrawFilledRect(screen, x, y, cs, cs, bgColor, true)

			// Region borders — draw thick line between cells of different regions
			borderColor := color.RGBA{0x33, 0x33, 0x33, 0xFF}
			borderW := float32(3)
			// Right border
			if c < g.Board.Size-1 && g.Board.Regions[r][c+1] != reg {
				vector.StrokeLine(screen, x+cs, y, x+cs, y+cs, borderW, borderColor, true)
			}
			// Bottom border
			if r < g.Board.Size-1 && g.Board.Regions[r+1][c] != reg {
				vector.StrokeLine(screen, x, y+cs, x+cs, y+cs, borderW, borderColor, true)
			}
			// Left border
			if c == 0 || g.Board.Regions[r][c-1] != reg {
				vector.StrokeLine(screen, x, y, x, y+cs, borderW, borderColor, true)
			}
			// Top border
			if r == 0 || g.Board.Regions[r-1][c] != reg {
				vector.StrokeLine(screen, x, y, x+cs, y, borderW, borderColor, true)
			}

			// Draw mark
			mark := g.Board.Marks[r][c]
			if mark == game.MarkQueen {
				isConflict := conflicts[[2]int{r, c}]
				sym := "Q"
				tw, th := text.Measure(sym, face, 0)
				opts := &text.DrawOptions{}
				opts.GeoM.Translate(float64(x)+float64(cs)/2-tw/2, float64(y)+float64(cs)/2-th/2)
				if isConflict {
					opts.ColorScale.ScaleWithColor(color.RGBA{0xCC, 0x00, 0x00, 0xFF})
				} else {
					opts.ColorScale.ScaleWithColor(color.RGBA{0x22, 0x22, 0x22, 0xFF})
				}
				text.Draw(screen, sym, face, opts)

				if isConflict {
					vector.StrokeRect(screen, x+2, y+2, cs-4, cs-4, 3, color.RGBA{0xCC, 0x00, 0x00, 0xFF}, true)
				}
			} else if mark == game.MarkX {
				sym := "X"
				tw, th := text.Measure(sym, face, 0)
				opts := &text.DrawOptions{}
				opts.GeoM.Translate(float64(x)+float64(cs)/2-tw/2, float64(y)+float64(cs)/2-th/2)
				opts.ColorScale.ScaleWithColor(color.RGBA{0x66, 0x66, 0x66, 0xAA})
				text.Draw(screen, sym, face, opts)
			}
		}
	}

	// Outer board border
	totalSize := float32(g.Board.Size * g.cellSize)
	vector.StrokeRect(screen, float32(g.offsetX), float32(g.offsetY), totalSize, totalSize, 3, color.RGBA{0x33, 0x33, 0x33, 0xFF}, true)

	// Grid lines (thin)
	for i := 1; i < g.Board.Size; i++ {
		x := float32(g.offsetX + i*g.cellSize)
		y := float32(g.offsetY + i*g.cellSize)
		vector.StrokeLine(screen, x, float32(g.offsetY), x, float32(g.offsetY)+totalSize, 1, color.RGBA{0x88, 0x88, 0x88, 0x88}, true)
		vector.StrokeLine(screen, float32(g.offsetX), y, float32(g.offsetX)+totalSize, y, 1, color.RGBA{0x88, 0x88, 0x88, 0x88}, true)
	}

	// Top bar — title
	title := fmt.Sprintf("Queens %dx%d", g.Board.Size, g.Board.Size)
	tw, _ := text.Measure(title, face, 0)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(float64(screenW)/2-tw/2, 20)
	opts.ColorScale.ScaleWithColor(color.RGBA{0x33, 0x33, 0x33, 0xFF})
	text.Draw(screen, title, face, opts)

	// Win message
	if g.Won {
		msg := "Congratulations! Puzzle Solved!"
		mw, _ := text.Measure(msg, face, 0)
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(screenW)/2-mw/2, float64(g.offsetY)-30)
		opts.ColorScale.ScaleWithColor(color.RGBA{0x00, 0x88, 0x00, 0xFF})
		text.Draw(screen, msg, face, opts)
	}

	// Bottom buttons
	g.drawButton(screen, face, "New Game", screenW/2-140-20, screenH-botBarHeight+12, 140, 36)
	g.drawButton(screen, face, "Menu", screenW/2+20, screenH-botBarHeight+12, 140, 36)
}

func (g *GameUI) drawButton(screen *ebiten.Image, face text.Face, label string, x, y, w, h int) {
	mx, my := ebiten.CursorPosition()
	hovered := mx >= x && mx <= x+w && my >= y && my <= y+h

	bg := color.RGBA{0xDD, 0xDD, 0xDD, 0xFF}
	if hovered {
		bg = color.RGBA{0xCC, 0xCC, 0xCC, 0xFF}
	}
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), bg, true)
	vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 2, color.RGBA{0x66, 0x66, 0x66, 0xFF}, true)

	tw, th := text.Measure(label, face, 0)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(float64(x)+float64(w)/2-tw/2, float64(y)+float64(h)/2-th/2)
	opts.ColorScale.ScaleWithColor(color.RGBA{0x22, 0x22, 0x22, 0xFF})
	text.Draw(screen, label, face, opts)
}
