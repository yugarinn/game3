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
	tilemapPositionX, tilemapPositionY := getPropTilemapPosition(prop.Type, prop.IsOpen)
	rec := rl.NewRectangle(tilemapPositionX, tilemapPositionY, prop.Width, prop.Height)
	rl.DrawTextureRec(r.Textures["tilemap"], rec, prop.Position, rl.White)

	if r.DebugMode {
		rl.DrawRectangleLines(int32(prop.HitboxRect.X), int32(prop.HitboxRect.Y), prop.HitboxRect.ToInt32().Width, prop.HitboxRect.ToInt32().Height, rl.Purple)
	}
}

func getPropTilemapPosition(propType PropType, isOpen bool) (float32, float32) {
	positions := map[PropType][]float32{
		PropKey:   {48, 32},
		PropDoor:  {88, 48},
		PropGrass: {40, 32},
	}

	if position, ok := positions[propType]; ok {
		positionX := position[0]
		positionY := position[1]

		if isOpen && propType == PropDoor {
			positionX += 16
		}

		return positionX, positionY
	}

	return 0, 0
}
