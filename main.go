package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/Cycl0o0/QueensGO/game"
	"github.com/Cycl0o0/QueensGO/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
)

type App struct {
	state  GameState
	menu   *ui.Menu
	gameUI *ui.GameUI
	face   text.Face

	screenW, screenH int
}

func NewApp() *App {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	face := &text.GoTextFace{
		Source: src,
		Size:   24,
	}

	return &App{
		state:   StateMenu,
		menu:    ui.NewMenu(),
		face:    face,
		screenW: 600,
		screenH: 700,
	}
}

func (a *App) Update() error {
	switch a.state {
	case StateMenu:
		sel := a.menu.Update(a.screenW, a.screenH)
		if sel >= 0 {
			size := ui.Difficulties[sel].Size
			regions := game.GeneratePuzzle(size)
			board := game.NewBoard(size, regions)
			a.gameUI = ui.NewGameUI(board)
			a.state = StatePlaying
		}

	case StatePlaying:
		a.gameUI.Update(a.screenW, a.screenH)
		if a.gameUI.WantMenu {
			a.state = StateMenu
			a.menu = ui.NewMenu()
		}
		if a.gameUI.WantNew {
			size := a.gameUI.Board.Size
			regions := game.GeneratePuzzle(size)
			board := game.NewBoard(size, regions)
			a.gameUI = ui.NewGameUI(board)
		}
	}
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(colorBg)

	switch a.state {
	case StateMenu:
		a.menu.Draw(screen, a.face)
	case StatePlaying:
		a.gameUI.Draw(screen, a.face)
	}
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	a.screenW = outsideWidth
	a.screenH = outsideHeight
	return outsideWidth, outsideHeight
}

var colorBg = color.RGBA{0xF5, 0xF5, 0xF0, 0xFF}

func main() {
	ebiten.SetWindowSize(600, 700)
	ebiten.SetWindowTitle("QueensGO")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	app := NewApp()
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
