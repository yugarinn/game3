package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"game3/assets"
)

type Tile struct {
	Position    rl.Vector2
	HitboxRect  rl.Rectangle
	Type        string
	Variant     string
	IsSolid     bool
	Sprite      rl.Texture2D
	TextureRect rl.Rectangle
	BaseColor   rl.Color
	HasTexture  bool
}

func (tile Tile) Draw() {
	if tile.HasTexture {
		rl.DrawTextureRec(tile.Sprite, tile.TextureRect, tile.Position, rl.White)
	} else {
		rl.DrawRectangle(
			int32(tile.Position.X),
			int32(tile.Position.Y),
			tile.TextureRect.ToInt32().Width,
			tile.TextureRect.ToInt32().Height,
			tile.BaseColor,
		)
	}
}

func BasicTestLevelLayout() []Tile {
	var layout []Tile
	groundImage := rl.LoadImageFromMemory(".png", assets.GROUND_SPRITE_DATA, int32(len(assets.GROUND_SPRITE_DATA)))

	// close floating platform
	for y := 0; y < 7; y++ {
		position := rl.NewVector2(200 + float32(8 * y), 124)

		tile := Tile{
			Position:    position,
			Type:        "grass",
			IsSolid:     true,
			TextureRect: rl.NewRectangle(200 + float32(8 * y), 124, 8, 8),
			HitboxRect:  rl.NewRectangle(position.X, position.Y, 8, 8),
			BaseColor:   grassRegularGreen,
		}

		layout = append(layout, tile)
    }

	// middle floating platform
	for y := 0; y < 7; y++ {
		position := rl.NewVector2(110 + float32(8 * y), 105)

		tile := Tile{
			Position:    position,
			Type:        "grass",
			IsSolid:     true,
			TextureRect: rl.NewRectangle(0, 0, 8, 8),
			HitboxRect:  rl.NewRectangle(position.X, position.Y, 8, 8),
			BaseColor:   grassRegularGreen,
		}

		layout = append(layout, tile)
    }

	// far floating platform
	for y := 0; y < 7; y++ {
		position := rl.NewVector2(20 + float32(8 * y), 90)

		tile := Tile{
			Position:    position,
			Type:        "grass",
			IsSolid:     true,
			TextureRect: rl.NewRectangle(0, 0, 8, 8),
			HitboxRect:  rl.NewRectangle(position.X, position.Y, 8, 8),
			BaseColor:   grassRegularGreen,
		}

		layout = append(layout, tile)
    }

	// right wall
	for y := 0; y < 21; y++ {
		position := rl.NewVector2(312, float32(y * 8))

		tile := Tile{
			Position:    position,
			Type:        "wall",
			IsSolid:     true,
			TextureRect: rl.NewRectangle(position.X, position.Y, 8, 8),
			HitboxRect:  rl.NewRectangle(position.X, position.Y, 8, 8),
			BaseColor:   rl.NewColor(0, 0, 0, 255),
		}

		layout = append(layout, tile)
    }

	// floor first row
	for y := 0; y < 40; y++ {
		position := rl.NewVector2(float32(8 * y), 172)

		tile := Tile{
			Position:    position,
			HitboxRect:  rl.NewRectangle(position.X, position.Y, 8, 8),
			Type:        "grass",
			IsSolid:     true,
			BaseColor:   grassRegularGreen,
			Sprite:      rl.LoadTextureFromImage(groundImage),
			TextureRect: rl.NewRectangle(float32(y % 2 * 8), 8, 8, 8),
			HasTexture:  true,
		}

		layout = append(layout, tile)
    }

	// floor second row
	for y := 0; y < 40; y++ {
		position := rl.NewVector2(float32(8 * y), 164)

		tile := Tile{
			Position:    position,
			HitboxRect:  rl.NewRectangle(position.X, position.Y, 8, 8),
			Type:        "grass",
			Sprite:      rl.LoadTextureFromImage(groundImage),
			TextureRect: rl.NewRectangle(float32(y % 2 * 8), 0, 8, 8),
			BaseColor:   grassRegularGreen,
			IsSolid:     true,
			HasTexture:  true,
		}

		layout = append(layout, tile)
    }

	return layout
}
