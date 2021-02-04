package main

import (
	"image/color"
	"os"

	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Facing int

const (
	Left Facing = iota
	Right
)

type KatieLevel struct {
	player   KatiePlayer
	space    *resolv.Space
	snow     []*Snow
	drawSnow bool
	bg       *ebiten.Image
}

type KatiePlayer struct {
	// posX   float64
	// posY   float64
	image  *ebiten.Image
	facing Facing
	rect   *resolv.Rectangle
}

type Snow struct {
	rect   *resolv.Rectangle
	plowed bool
}

func NewKatieLevel() *KatieLevel {
	ktImage := loadPNG("./assets/kt.png")
	w, h := ktImage.Size()

	k := &KatieLevel{
		player: KatiePlayer{
			image:  ktImage,
			facing: Left,
			rect:   resolv.NewRectangle(ScreenX/2, ScreenY/2, int32(w)/10, int32(h)/10),
		},
		space:    resolv.NewSpace(),
		snow:     []*Snow{},
		drawSnow: true,
		bg:       loadPNG("./assets/geopolis.png"),
	}
	k.player.rect.SetData(k.player)

	// Generate snow!
	numX := int32((ScreenX / 32) + 1)
	numY := int32((ScreenY / 32) + 1)

	var i, j int32
	for i = 0; i < numX; i++ {
		for j = 0; j < numY; j++ {
			s := &Snow{
				rect:   resolv.NewRectangle(i*32, j*32, 32, 32),
				plowed: false,
			}
			s.rect.SetData(s)
			s.rect.AddTags("snow")
			k.space.Add(s.rect)
			k.snow = append(k.snow, s)
		}
	}
	return k
}

func (k *KatieLevel) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		k.player.rect.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		k.player.rect.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		k.player.rect.X -= 1
		k.player.facing = Left
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		k.player.rect.X += 1
		k.player.facing = Right
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		k.drawSnow = !k.drawSnow
	}

	if k.player.rect.X > ScreenX {
		k.player.rect.X = 0
	}
	if k.player.rect.X < 0 {
		k.player.rect.X = ScreenX
	}
	if k.player.rect.Y > ScreenY {
		k.player.rect.Y = 0
	}
	if k.player.rect.Y < 0 {
		k.player.rect.Y = ScreenY
	}

	collisions := k.space.GetCollidingShapes(k.player.rect)
	for _, shape := range *collisions {
		s := shape.GetData().(*Snow)
		s.plowed = true
	}

	return nil
}

func (k *KatieLevel) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// draw background
	screen.DrawImage(k.bg, op)

	// Draw katie
	_, w := k.player.image.Size()
	if k.player.facing == Right {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(w), 0)
	}
	op.GeoM.Scale(0.10, 0.10)

	op.GeoM.Translate(float64(k.player.rect.X), float64(k.player.rect.Y))
	screen.DrawImage(k.player.image, op)

	// Draw snow
	for _, s := range k.snow {
		if !s.plowed && k.drawSnow {
			ebitenutil.DrawRect(screen, float64(s.rect.X), float64(s.rect.Y), 32, 32, color.White)
		}
	}
}
