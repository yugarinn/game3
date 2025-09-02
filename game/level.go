package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Level struct {
	ID         string            `json:"iid"`
	Name       string            `json:"identifier"`
	Neighbours []*LevelNeighbour `json:"__neighbours"`
	Layers     []*LevelLayer     `json:"layerInstances"`
	Props      []Prop
}

type LevelLayer struct {
	ID          string      `json:"iid"`
	Name        string      `json:"__identifier"`
	TilesetPath []*Tile     `json:"__tilesetRelPath"`
	RawLayout   []*LDtkTile `json:"gridTiles"`
	Layout      []*Tile
}

type LevelNeighbour struct {
	LevelID   string `json:"levelIid"`
	Direction string `json:"dir"`
}

func (l *Level) Load() {
	l.GetBackgroundLayer().LoadLayout()
	l.GetGroundLayer().LoadLayout()
}

func (l *Level) Draw(r *Renderer) {
	l.DrawBackground(r)
	l.DrawGround(r)
}

func (l *Level) DrawBackground(r *Renderer) {
	for _, tile := range l.GetBackgroundLayer().Layout {
		// TODO: could I maybe do just r.DrawGroundTile(tile)????/
		rec := rl.NewRectangle(tile.SpritePosition.X, tile.SpritePosition.Y, 8, 8)
		r.DrawSprite("ground", rec, tile.Position)
	}
}

func (l *Level) DrawGround(r *Renderer) {
	for _, tile := range l.GetGroundLayer().Layout {
		// TODO: could I maybe do just r.DrawGroundTile(tile)????/
		rec := rl.NewRectangle(tile.SpritePosition.X, tile.SpritePosition.Y, 8, 8)
		r.DrawSprite("ground", rec, tile.Position)
	}
}

func (l *Level) GetBackgroundLayer() *LevelLayer {
	for _, layer := range l.Layers {
		if layer.Name == "Background" {
			return layer
		}
	}

	return nil
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
