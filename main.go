package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenX = 640
	ScreenY = 480
)

type GameState int

const (
	Playing GameState = iota
	Quitting
)

type Sprite struct {
	image *ebiten.Image
}

type Level interface {
	Update() error
	Draw(*ebiten.Image)
}

type Game struct {
	levels       map[string]Level
	currentLevel string
}

func (g *Game) Update() error {
	level := g.levels[g.currentLevel]
	return level.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	level := g.levels[g.currentLevel]
	level.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenX, ScreenY
}

func main() {

	fmt.Println("starting")
	ebiten.SetWindowSize(2*ScreenX, 2*ScreenY)
	ebiten.SetWindowTitle("Reed's Game!")

	georgeLevel := NewGeorgeLevel()

	game := &Game{
		levels: map[string]Level{
			"GEORGE": georgeLevel,
		},
		currentLevel: "GEORGE",
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)

	}

}
