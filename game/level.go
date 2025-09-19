package game

import (
	"math"
	"math/rand"
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
	Particles  []*Particle
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

type Particle struct {
	Position     rl.Vector2
	Velocity     rl.Vector2
	FramesToLive int
	Seed         float32
}

func (l *Level) Load() {
	l.GetGroundLayer().LoadLayout()
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
	l.DrawParticles(r)
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

func (l *Level) LoadParticles() {
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

func NewParticle() *Particle {
	positionX := rand.Intn(320)
	positionY := rand.Intn(180)

	return &Particle{
		Position: rl.NewVector2(float32(positionX), float32(positionY)),
		Velocity: rl.NewVector2(0, 0),
		FramesToLive: 240 + rand.Intn(600),
		Seed: float32(rl.GetRandomValue(0, 1000)) * 0.01,
	}
}

func (p *Particle) UpdatePosition(delta float32) bool {
	thermalStrength := float32(15.0)
	turbulence := float32(8.0)
	drift := float32(5.0)

	posX := p.Position.X * 0.01
	posY := p.Position.Y * 0.01

	thermal := float32(math.Sin(float64(posY*2.1 + p.Seed))) * thermalStrength
	turbulenceX := float32(math.Sin(float64(posX*3.7 + p.Seed*1.3))) * turbulence
	turbulenceY := float32(math.Cos(float64(posY*2.9 + p.Seed*0.7))) * turbulence * 0.6
	horizontalDrift := drift + float32(math.Sin(float64(posX*1.5 + p.Seed))) * drift * 0.3

	totalForceX := horizontalDrift + turbulenceX
	totalForceY := thermal + turbulenceY

	resistance := float32(0.92)
	p.Velocity.X = p.Velocity.X*resistance + totalForceX*delta
	p.Velocity.Y = p.Velocity.Y*resistance + totalForceY*delta

	p.Position.X += p.Velocity.X * delta
	p.Position.Y += p.Velocity.Y * delta

	p.FramesToLive--

	return p.FramesToLive > 0
}
