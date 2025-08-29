package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Level struct {
	ID     string        `json:"iid"`
	Name   string        `json:"identifier"`
	Layers []*LevelLayer `json:"layerInstances"`
	Props  []Prop
}

type LevelLayer struct {
	ID          string  `json:"iid"`
	Name        string  `json:"__identifier"`
	TilesetPath []*Tile `json:"__tilesetRelPath"`
	RawLayout   []*LDtkTile `json:"gridTiles"`
	Layout      []*Tile
}

func (l *Level) Load() {
	l.GetGroundLayer().LoadLayout()
}

func (l *Level) Draw() {
	l.DrawGround()
}

func (l *Level) DrawGround() {
	for _, tile := range l.GetGroundLayer().Layout {
		rl.DrawRectangle(int32(tile.Position.X), int32(tile.Position.Y), 8, 8, rl.Pink)
	}
}

func (l *Level) GetGroundLayer() *LevelLayer {
	for _, layer := range l.Layers {
		if layer.Name == "Ground" {
			return layer
		}
	}

	return nil
}

func (ll *LevelLayer) LoadLayout() {
	var tiles []*Tile
	for _, rawTile := range ll.RawLayout {
		tiles = append(tiles, NewTileFromLDtk(rawTile))
	}

	ll.Layout = tiles
}
