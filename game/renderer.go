package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Renderer struct {
	Textures  map[string]rl.Texture2D
	DebugMode bool
}

func (r *Renderer) DrawSprite(textureID string, rec rl.Rectangle, position rl.Vector2) {
	texture := r.Textures[textureID]
	rl.DrawTextureRec(texture, rec, position, rl.White)

	if r.DebugMode {
		rl.DrawRectangleLines(int32(position.X), int32(position.Y), int32(rec.Width), int32(rec.Height), rl.Red)
	}
}

func (r *Renderer) DrawBackground(textureID string) {
	texture := r.Textures[textureID]
	rl.DrawTexture(texture, 0, 0, rl.White)
}

func (r *Renderer) DrawParticle(particle *Particle) {
	rl.DrawPixel(int32(particle.Position.X), int32(particle.Position.Y), rl.White)
}

func (r *Renderer) DrawProp(prop *Prop) {
	rl.DrawRectangle(int32(prop.Position[0]), int32(prop.Position[1]), 8, 8, rl.Green)

	if r.DebugMode {
		rl.DrawRectangleLines(int32(prop.Position[0]), int32(prop.Position[1]), 8, 8, rl.Red)
	}
}
