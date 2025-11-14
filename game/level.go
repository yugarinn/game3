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
	Props      []*Prop
	Particles  []*Particle
}

type LevelLayer struct {
	ID          string        `json:"iid"`
	Name        string        `json:"__identifier"`
	TilesetPath []*Tile       `json:"__tilesetRelPath"`
	RawLayout   []*LDtkTile   `json:"gridTiles"`
	RawEntities []*LDtkEntity `json:"entityInstances"`
	Entities    []*Prop
	Layout      []*Tile
}

type LevelNeighbour struct {
	LevelID   string `json:"levelIid"`
	Direction string `json:"dir"`
}

func (l *Level) Load() {
	l.GetGroundLayer().LoadLayout()
	l.LoadProps()
	l.LoadParticles()

	if l.Background != "" {
		l.Background = strings.TrimSuffix(filepath.Base(l.Background), filepath.Ext(l.Background))
	}
}

func (l *Level) Unload() {
	l.Particles = []*Particle{}
}

func (l *Level) Tick(delta float32) {
	for i, particle := range(l.Particles) {
		particle.UpdatePosition(delta)

		if particle.FramesToLive < 0 {
			l.Particles[i] = NewParticle()
		}
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

func (l *Level) DrawProps(r *Renderer) {
	for _, prop := range l.Props {
		r.DrawProp(prop)
	}
}

func (l *Level) DrawParticles(r *Renderer) {
	for _, particle := range(l.Particles) {
		r.DrawParticle(particle)
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

func (l *Level) GetEntitiesLayer() *LevelLayer {
	for _, layer := range l.Layers {
		if layer.Name == "Entities" {
			return layer
		}
	}

	return nil
}

func (l *Level) LoadParticles() {
	// TODO: the particles density should be determined by a level prop
	for range(20) {
		l.Particles = append(l.Particles, NewParticle())
	}
}

func (ll *LevelLayer) LoadLayout() {
	var tiles []*Tile
	for _, rawTile := range ll.RawLayout {
		tiles = append(tiles, NewTileFromLDtk(rawTile))
	}

	ll.Layout = tiles
}

func (l *Level) LoadProps() {
	propsLayer := l.GetEntitiesLayer()
	if propsLayer == nil {
		return
	}

	var props []*Prop
	for _, entity := range propsLayer.RawEntities {
		props = append(props, NewPropFromLDtk(entity))
	}

	l.Props = props
}
