package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	Position       rl.Vector2
	SpritePosition rl.Vector2
	HitboxRect     rl.Rectangle
}

type LDtkTile struct {
	Px  []float32
	T   int
	Src []float32
}

const (
	TileSize    = 8
	TilesPerRow = 17
)

func NewTileFromLDtk(ldtkTile *LDtkTile) *Tile {
	positionX := float32(ldtkTile.Px[0])
	positionY := float32(ldtkTile.Px[1])

	return &Tile{
		Position:       rl.NewVector2(positionX, positionY),
		SpritePosition: rl.NewVector2(ldtkTile.Src[0], ldtkTile.Src[1]),
		HitboxRect:     rl.NewRectangle(positionX, positionY, TileSize, TileSize),
	}
}
