package main

import (
	"fmt"
	"image"
	"os"

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

func loadPNG(path string) *ebiten.Image {
	fd, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	im, _, err := image.Decode(fd)
	if err != nil {
		panic(err)
	}

	eim := ebiten.NewImageFromImage(im)
	return eim
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
	katieLevel := NewKatieLevel()

	game := &Game{
		levels: map[string]Level{
			"GEORGE": georgeLevel,
			"KATIE":  katieLevel,
		},
		currentLevel: "KATIE",
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)

	}

}
