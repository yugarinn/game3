package game

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ParticleSource struct {
	ThermalStrength float32
	Turbulence      float32
	Drift           float32
}

type Particle struct {
	Position     rl.Vector2
	Velocity     rl.Vector2
	FramesToLive int
	Seed         float32
	Type         string // TODO: dust, rain, fire...
	Source       *ParticleSource
}

func NewParticle() *Particle {
	positionX := rand.Intn(320)
	positionY := rand.Intn(180)

	source := &ParticleSource{
		ThermalStrength: 15,
		Turbulence: 100,
		Drift: 2,
	}

	return &Particle{
		Position: rl.NewVector2(float32(positionX), float32(positionY)),
		Velocity: rl.NewVector2(0, 0),
		FramesToLive: 600 + rand.Intn(600),
		Seed: float32(rl.GetRandomValue(0, 1000)) * 0.01,
		Source: source,
	}
}

// TODO: no particles seem to be moving leftwards...
// TODO: I should have several methods for each particle type...
func (p *Particle) UpdatePosition(delta float32) {
	thermalStrength := float32(15.0)
	turbulence := float32(8.0)
	drift := float32(5.0)
	driftDirection := rand.Intn(2)

	posX := p.Position.X * 0.01
	posY := p.Position.Y * 0.01

	thermal := float32(math.Sin(float64(posY*2.1 + p.Seed))) * thermalStrength
	turbulenceX := float32(math.Sin(float64(posX*3.7 + p.Seed*1.3))) * turbulence
	turbulenceY := float32(math.Cos(float64(posY*2.9 + p.Seed*0.7))) * turbulence * 0.6
	horizontalDrift := drift + float32(math.Sin(float64(posX*1.5 + p.Seed))) * drift * 0.3

	if driftDirection == 0 {
		horizontalDrift = -horizontalDrift
	}

	totalForceX := horizontalDrift + turbulenceX
	totalForceY := thermal + turbulenceY

	resistance := float32(0.92)
	p.Velocity.X = p.Velocity.X*resistance + totalForceX*delta
	p.Velocity.Y = p.Velocity.Y*resistance + totalForceY*delta

	p.Position.X += p.Velocity.X * delta
	p.Position.Y += p.Velocity.Y * delta

	p.FramesToLive--
}
