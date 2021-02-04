package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	SpriteTileX = 48
	SpriteTileY = 48

	PlayerAccel     = 5.0
	PlayerJumpSpeed = 8.0
	Gravity         = 15.0
	Friction        = 0.25
	MaxSpeed        = 3
)

type GeorgeLevel struct {
	state  GameState
	player Player
	space  *resolv.Space
}

type Player struct {
	sprite   Sprite
	velX     float64
	velY     float64
	rect     *resolv.Rectangle
	grounded bool
}

func NewGeorgeLevel() *GeorgeLevel {
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

	georgeLevel := &GeorgeLevel{
		state:  Playing,
		player: player,
		space:  space,
	}

	return georgeLevel
}

func GetTileRect(tilex, tiley int) image.Rectangle {
	return image.Rect(tilex*SpriteTileX,
		tiley*SpriteTileY,
		(tilex+1)*SpriteTileX,
		(tiley+1)*SpriteTileY)
}

func (g *GeorgeLevel) Update() error {
	delta := 1.0 / float64(ebiten.MaxTPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.velX -= PlayerAccel*delta + Friction
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.velX += PlayerAccel*delta + Friction
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.state = Quitting
	}

	if res := g.space.Resolve(g.player.rect, 0, 3); res.Colliding() {
		g.player.velY = 0
		g.player.grounded = true
	} else {
		g.player.grounded = false
	}

	if g.player.grounded && (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyUp)) {
		g.player.velY -= PlayerJumpSpeed
	}

	if res := g.space.Resolve(g.player.rect, int32(g.player.velX), 0); res.Colliding() {
		fmt.Println("collide x")
		g.player.rect.X += res.ResolveX
		g.player.velX = 0
	}

	if res := g.space.Resolve(g.player.rect, 0, int32(g.player.velY)); res.Colliding() {
		fmt.Println("collide y")
		g.player.rect.Y += res.ResolveY
		g.player.velY = 0
	}
	g.player.velY += delta * Gravity

	if g.player.velX > Friction {
		g.player.velX -= Friction
	} else if g.player.velX < -Friction {
		g.player.velX += Friction
	} else {
		g.player.velX = 0
	}

	if g.player.velX > MaxSpeed {
		g.player.velX = MaxSpeed
	}
	if g.player.velX < -MaxSpeed {
		g.player.velX = -MaxSpeed
	}

	g.player.rect.X += int32(g.player.velX)
	g.player.rect.Y += int32(g.player.velY)

	// wrapping
	if g.player.rect.X > ScreenX {
		g.player.rect.X = 0
	}
	if g.player.rect.X < 0 {
		g.player.rect.X = ScreenX
	}
	switch g.state {
	case Playing:
	case Quitting:
		os.Exit(0)
	}

	return nil
}

func (g *GeorgeLevel) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// draw floor
	ebitenutil.DrawRect(screen, 0, ScreenY-10, ScreenX, 10, color.White)

	// ebitenutil.DrawRect(screen, 0, 0, ScreenX, ScreenY, color.White)
	op.GeoM.Translate(float64(g.player.rect.X), float64(g.player.rect.Y))
	screen.DrawImage(g.player.sprite.image, op)

}
