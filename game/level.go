package game

import (
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Level struct {
	ID         string            `json:"iid"`
	Name       string            `json:"identifier"`
	Neighbours []*LevelNeighbour `json:"__neighbours"`
	Layers     []*LevelLayer     `json:"layerInstances"`
	Background string            `json:"bgRelPath"`
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
	l.GetGroundLayer().LoadLayout()

	if l.Background != "" {
		l.Background = strings.TrimSuffix(filepath.Base(l.Background), filepath.Ext(l.Background))
	}
}

func (l *Level) Draw(r *Renderer) {
	l.DrawBackground(r)
	l.DrawGround(r)
}

func (l *Level) DrawBackground(r *Renderer) {
	r.DrawBackground(l.Background)
}

func (l *Level) DrawGround(r *Renderer) {
	for _, tile := range l.GetGroundLayer().Layout {
		// TODO: could I maybe do just r.DrawGroundTile(tile)????/
		rec := rl.NewRectangle(tile.SpritePosition.X, tile.SpritePosition.Y, 8, 8)
		r.DrawSprite("tileset_ground", rec, tile.Position)
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
