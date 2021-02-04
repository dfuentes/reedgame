package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/SolarLune/resolv"
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

	spritesfd, err := os.Open("./assets/monkeysprites.png")
	if err != nil {
		panic(err)
	}
	spritesImage, err := png.Decode(spritesfd)
	if err != nil {
		panic(err)
	}

	spritesEbitenImage := ebiten.NewImageFromImage(spritesImage)

	player := Player{
		sprite: Sprite{
			image: spritesEbitenImage.SubImage(GetTileRect(0, 2)).(*ebiten.Image),
		},
		rect: resolv.NewRectangle(ScreenX/2, ScreenY/2, SpriteTileX, SpriteTileY),
		velX: float64(0),
		velY: float64(0),
	}

	space := resolv.NewSpace()

	floor := resolv.NewRectangle(0, ScreenY-10, ScreenX, 10)
	space.Add(floor)

	fmt.Println("starting")
	ebiten.SetWindowSize(2*ScreenX, 2*ScreenY)
	ebiten.SetWindowTitle("test")

	georgeLevel := &GeorgeLevel{
		state:  Playing,
		player: player,
		space:  space,
	}

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
