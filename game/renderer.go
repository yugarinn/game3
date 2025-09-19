package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Renderer struct {
	Textures map[string]rl.Texture2D
}

func (r *Renderer) DrawSprite(textureID string, rec rl.Rectangle, position rl.Vector2) {
	texture := r.Textures[textureID]
	rl.DrawTextureRec(texture, rec, position, rl.White)

	if DEBUG {
		rl.DrawRectangleLines(int32(position.X), int32(position.Y), int32(rec.Width), int32(rec.Height), rl.Red)
	}
}

func (r *Renderer) DrawBackground(textureID string) {
	fmt.Println(textureID)
	texture := r.Textures[textureID]
	rl.DrawTexture(texture, 0, 0, rl.White)
}
