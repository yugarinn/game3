package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	Position   rl.Vector2
	TileNumber int `json:"t"`
	HitboxRect rl.Rectangle
}

type LDtkTile struct {
    Px []int
	T  int
}

func NewTileFromLDtk(ldtkTile *LDtkTile) *Tile {
	positionX := float32(ldtkTile.Px[0])
	positionY := float32(ldtkTile.Px[1])

	return &Tile{
		Position:   rl.NewVector2(positionX, positionY),
		TileNumber: ldtkTile.T,
		HitboxRect: rl.NewRectangle(positionX, positionY, 8, 8),
	}
}
