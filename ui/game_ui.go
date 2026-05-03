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
	btnWidth     = 110
	btnHeight    = 36
	btnGap       = 15
)

// GameUI handles rendering and input for the puzzle board.
type GameUI struct {
	Board       *game.Board
	cellSize    int
	offsetX     int
	offsetY     int
	justClicked bool
	Won         bool
	WantMenu    bool
	WantNew     bool
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

	var mx, my int
	var inputActive bool
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my = ebiten.CursorPosition()
		inputActive = true
	} else if touches := ebiten.AppendTouchIDs(nil); len(touches) > 0 {
		mx, my = ebiten.TouchPosition(touches[0])
		inputActive = true
	}

	if !inputActive {
		g.justClicked = false
		return
	}
	if g.justClicked {
		return
	}
	g.justClicked = true

	// Bottom buttons — three buttons centered
	totalBtnsW := 3*btnWidth + 2*btnGap
	startX := (screenW - totalBtnsW) / 2
	bY := screenH - botBarHeight + 12

	if image.Pt(mx, my).In(image.Rect(startX, bY, startX+btnWidth, bY+btnHeight)) {
		g.WantNew = true
		return
	}
	resetX := startX + btnWidth + btnGap
	if image.Pt(mx, my).In(image.Rect(resetX, bY, resetX+btnWidth, bY+btnHeight)) {
		g.Board.Reset()
		g.Won = false
		return
	}
	menuX := startX + 2*(btnWidth+btnGap)
	if image.Pt(mx, my).In(image.Rect(menuX, bY, menuX+btnWidth, bY+btnHeight)) {
		g.WantMenu = true
		return
	}

	if g.Won {
		return
	}

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
	totalSize := float32(g.Board.Size * g.cellSize)

	// Draw region cells and thick region borders
	for r := 0; r < g.Board.Size; r++ {
		for c := 0; c < g.Board.Size; c++ {
			x := float32(g.offsetX + c*g.cellSize)
			y := float32(g.offsetY + r*g.cellSize)
			reg := g.Board.Regions[r][c]
			vector.DrawFilledRect(screen, x, y, cs, cs, RegionColors[reg%len(RegionColors)], true)

			borderColor := color.RGBA{0x33, 0x33, 0x33, 0xFF}
			const bw = float32(3)
			if c < g.Board.Size-1 && g.Board.Regions[r][c+1] != reg {
				vector.StrokeLine(screen, x+cs, y, x+cs, y+cs, bw, borderColor, true)
			}
			if r < g.Board.Size-1 && g.Board.Regions[r+1][c] != reg {
				vector.StrokeLine(screen, x, y+cs, x+cs, y+cs, bw, borderColor, true)
			}
			if c == 0 || g.Board.Regions[r][c-1] != reg {
				vector.StrokeLine(screen, x, y, x, y+cs, bw, borderColor, true)
			}
			if r == 0 || g.Board.Regions[r-1][c] != reg {
				vector.StrokeLine(screen, x, y, x+cs, y, bw, borderColor, true)
			}
		}
	}

	// Tint entire rows and columns containing conflicting queens
	conflictRows := map[int]bool{}
	conflictCols := map[int]bool{}
	for pos := range conflicts {
		conflictRows[pos[0]] = true
		conflictCols[pos[1]] = true
	}
	tint := color.RGBA{0xFF, 0x44, 0x44, 0x28}
	for row := range conflictRows {
		ry := float32(g.offsetY + row*g.cellSize)
		vector.DrawFilledRect(screen, float32(g.offsetX), ry, totalSize, cs, tint, true)
	}
	for col := range conflictCols {
		cx := float32(g.offsetX + col*g.cellSize)
		vector.DrawFilledRect(screen, cx, float32(g.offsetY), cs, totalSize, tint, true)
	}

	// Outer border and thin grid lines
	vector.StrokeRect(screen, float32(g.offsetX), float32(g.offsetY), totalSize, totalSize, 3, color.RGBA{0x33, 0x33, 0x33, 0xFF}, true)
	for i := 1; i < g.Board.Size; i++ {
		x := float32(g.offsetX + i*g.cellSize)
		y := float32(g.offsetY + i*g.cellSize)
		vector.StrokeLine(screen, x, float32(g.offsetY), x, float32(g.offsetY)+totalSize, 1, color.RGBA{0x88, 0x88, 0x88, 0x88}, true)
		vector.StrokeLine(screen, float32(g.offsetX), y, float32(g.offsetX)+totalSize, y, 1, color.RGBA{0x88, 0x88, 0x88, 0x88}, true)
	}

	// Draw cell marks
	for r := 0; r < g.Board.Size; r++ {
		for c := 0; c < g.Board.Size; c++ {
			x := float32(g.offsetX + c*g.cellSize)
			y := float32(g.offsetY + r*g.cellSize)
			switch g.Board.Marks[r][c] {
			case game.MarkQueen:
				g.drawQueenSymbol(screen, face, x, y, cs, conflicts[[2]int{r, c}])
			case game.MarkX:
				g.drawXMark(screen, x, y, cs)
			}
		}
	}

	// Top bar: centered title
	title := fmt.Sprintf("Queens %dx%d", g.Board.Size, g.Board.Size)
	tw, _ := text.Measure(title, face, 0)
	{
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(screenW)/2-tw/2, 18)
		opts.ColorScale.ScaleWithColor(color.RGBA{0x33, 0x33, 0x33, 0xFF})
		text.Draw(screen, title, face, opts)
	}

	// Top bar: queen counter (top-right)
	queensPlaced := len(g.Board.Queens())
	counter := fmt.Sprintf("%d / %d", queensPlaced, g.Board.Size)
	cw, _ := text.Measure(counter, face, 0)
	{
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(screenW)-float64(cw)-20, 18)
		counterColor := color.RGBA{0x77, 0x77, 0x77, 0xFF}
		if g.Won {
			counterColor = color.RGBA{0x00, 0x99, 0x00, 0xFF}
		} else if len(conflicts) > 0 {
			counterColor = color.RGBA{0xBB, 0x22, 0x22, 0xFF}
		}
		opts.ColorScale.ScaleWithColor(counterColor)
		text.Draw(screen, counter, face, opts)
	}

	// Win overlay panel
	if g.Won {
		panelW := float32(360)
		panelH := float32(110)
		px := float32(screenW)/2 - panelW/2
		py := float32(screenH)/2 - panelH/2

		vector.DrawFilledRect(screen, px+4, py+4, panelW, panelH, color.RGBA{0, 0, 0, 0x50}, true)
		vector.DrawFilledRect(screen, px, py, panelW, panelH, color.RGBA{0xF0, 0xFF, 0xF2, 0xF8}, true)
		vector.StrokeRect(screen, px, py, panelW, panelH, 3, color.RGBA{0x00, 0xAA, 0x44, 0xFF}, true)

		msg1 := "Puzzle Solved!"
		mw1, _ := text.Measure(msg1, face, 0)
		{
			const scale = 1.3
			opts := &text.DrawOptions{}
			opts.GeoM.Scale(scale, scale)
			opts.GeoM.Translate(float64(screenW)/2-mw1*scale/2, float64(py)+14)
			opts.ColorScale.ScaleWithColor(color.RGBA{0x00, 0x99, 0x00, 0xFF})
			text.Draw(screen, msg1, face, opts)
		}

		msg2 := "All queens placed correctly!"
		mw2, _ := text.Measure(msg2, face, 0)
		{
			opts := &text.DrawOptions{}
			opts.GeoM.Translate(float64(screenW)/2-mw2/2, float64(py)+70)
			opts.ColorScale.ScaleWithColor(color.RGBA{0x44, 0x44, 0x44, 0xFF})
			text.Draw(screen, msg2, face, opts)
		}
	}

	// Bottom buttons
	totalBtnsW := 3*btnWidth + 2*btnGap
	startX := (screenW - totalBtnsW) / 2
	bY := screenH - botBarHeight + 12

	g.drawButton(screen, face, "New Game", startX, bY, btnWidth, btnHeight,
		color.RGBA{0xDD, 0xDD, 0xDD, 0xFF}, color.RGBA{0xC5, 0xC5, 0xC5, 0xFF})
	g.drawButton(screen, face, "Reset", startX+btnWidth+btnGap, bY, btnWidth, btnHeight,
		color.RGBA{0xF5, 0xE8, 0xCC, 0xFF}, color.RGBA{0xE5, 0xD0, 0xAA, 0xFF})
	g.drawButton(screen, face, "Menu", startX+2*(btnWidth+btnGap), bY, btnWidth, btnHeight,
		color.RGBA{0xDD, 0xDD, 0xDD, 0xFF}, color.RGBA{0xC5, 0xC5, 0xC5, 0xFF})
}

func (g *GameUI) drawQueenSymbol(screen *ebiten.Image, face text.Face, x, y, cs float32, isConflict bool) {
	cx := x + cs/2
	cy := y + cs/2
	var col color.RGBA
	if isConflict {
		col = color.RGBA{0xCC, 0x00, 0x00, 0xFF}
	} else {
		col = color.RGBA{0x22, 0x22, 0x22, 0xFF}
	}

	radius := cs * 0.34
	vector.DrawFilledCircle(screen, cx, cy, radius, color.RGBA{0xFF, 0xFF, 0xFF, 0xC0}, true)
	vector.StrokeCircle(screen, cx, cy, radius, 2.5, col, true)

	sym := "Q"
	tw, th := text.Measure(sym, face, 0)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(float64(cx)-tw/2, float64(cy)-th/2)
	opts.ColorScale.ScaleWithColor(col)
	text.Draw(screen, sym, face, opts)
}

func (g *GameUI) drawXMark(screen *ebiten.Image, x, y, cs float32) {
	col := color.RGBA{0x77, 0x77, 0x77, 0xCC}
	p := cs * 0.3
	vector.StrokeLine(screen, x+p, y+p, x+cs-p, y+cs-p, 1.5, col, true)
	vector.StrokeLine(screen, x+cs-p, y+p, x+p, y+cs-p, 1.5, col, true)
}

func (g *GameUI) drawButton(screen *ebiten.Image, face text.Face, label string, x, y, w, h int, bg, bgHover color.RGBA) {
	mx, my := ebiten.CursorPosition()
	hovered := mx >= x && mx <= x+w && my >= y && my <= y+h

	bgCol := bg
	if hovered {
		bgCol = bgHover
	}
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), bgCol, true)

	borderCol := color.RGBA{0x66, 0x66, 0x66, 0xFF}
	if hovered {
		borderCol = color.RGBA{0x33, 0x33, 0x33, 0xFF}
	}
	vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 2, borderCol, true)

	tw, th := text.Measure(label, face, 0)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(float64(x)+float64(w)/2-tw/2, float64(y)+float64(h)/2-th/2)
	opts.ColorScale.ScaleWithColor(color.RGBA{0x22, 0x22, 0x22, 0xFF})
	text.Draw(screen, label, face, opts)
}
