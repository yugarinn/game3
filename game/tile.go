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
    Px []float32
	T  int
}

const (
    TileSize = 8
    TilesPerRow = 17
)

func NewTileFromLDtk(ldtkTile *LDtkTile) *Tile {
	positionX := float32(ldtkTile.Px[0])
	positionY := float32(ldtkTile.Px[1])

	// This assumes that all tilemaps have a width of 136px
	// there are 17 8x8 tiles in a tileset of width 136

	// But... actually, LDtk already stores the position
	// in the sprite values under the `src` field...
	//
	// I'm leaving this here because I felt pretty
	// smart having figured this out, even though
	// I'll probably have to remove it when I
	// inevitabilly run into tilesets with different
	// widths... ¯\_(ツ)_/¯
	spriteXPosition := float32(ldtkTile.T % TilesPerRow * TileSize)
	spriteYPosition := float32(int(ldtkTile.T / TilesPerRow) * TileSize)

	return &Tile{
		Position:       rl.NewVector2(positionX, positionY),
		SpritePosition: rl.NewVector2(spriteXPosition, spriteYPosition),
		HitboxRect:     rl.NewRectangle(positionX, positionY, TileSize, TileSize),
	}
}
