package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Level struct {
	ID     string `json:"iid"`
	Name   string `json:"identifier"`
	Layers []*LevelLayer
	Props  []Prop
}

type LevelLayer struct {
	ID          string  `json:"iid"`
	Name        string  `json:"__identifier"`
	TilesetPath []*Tile `json:"__tilesetRelPath"`
	RawLayout   []*LDtkTile `json:"gridTiles"`
	Layout      []*Tile
}

func (l *Level) Draw() {
	l.DrawGround()
}

func (l *Level) DrawGround() {
	for _, tile := range l.GetGround().Layout {
		fmt.Println("TILE")
		fmt.Println(tile)
		rl.DrawRectangle(int32(tile.Position.X), int32(tile.Position.Y), 8, 8, rl.Pink)
	}
}

func (l *Level) GetGround() *LevelLayer {
	var groundLayer LevelLayer

	for _, layer := range l.Layers {
		if layer.Name == "Ground" {
			groundLayer = *layer
			break
		}
	}

	return &groundLayer
}

func (ll *LevelLayer) LoadLayout() {
	var tiles []*Tile
	for _, rawTile := range ll.RawLayout {
		tiles = append(tiles, NewTileFromLDtk(rawTile))
	}

	ll.Layout = tiles
}
