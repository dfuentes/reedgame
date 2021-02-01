package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenX = 320
	ScreenY = 240
)

type GameState int

const (
	Playing GameState = iota
	Quitting
)

type Sprite struct {
	image *ebiten.Image
	posX  float64
	posY  float64
}

type Game struct {
	state  GameState
	player Sprite
}

func (g *Game) Update() error {
	// Handle input
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.posY -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.posY += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.posX -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.posX += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.state = Quitting
	}

	switch g.state {
	case Playing:
	case Quitting:
		os.Exit(0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(g.player.posX, g.player.posY)
	screen.DrawImage(g.player.image, op)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenX, ScreenY
}

func main() {

	gfd, err := os.Open("./assets/george.png")
	if err != nil {
		panic(err)
	}
	gim, err := png.Decode(gfd)
	if err != nil {
		panic(err)
	}

	geim := ebiten.NewImageFromImage(gim)

	gsprite := Sprite{
		image: geim,
		posX:  float64(ScreenX / 2),
		posY:  float64(ScreenY / 2),
	}

	fmt.Println("starting")
	ebiten.SetWindowSize(2*ScreenX, 2*ScreenY)
	ebiten.SetWindowTitle("test")

	game := &Game{
		state:  Playing,
		player: gsprite,
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)

	}

}
