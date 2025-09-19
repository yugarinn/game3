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
