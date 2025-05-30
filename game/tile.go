package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	Type        string
	Variant     string
	IsSolid     bool
	Sprite      rl.Texture2D
	TextureRect rl.Rectangle
	BaseColor   rl.Color
}

func (tile Tile) Draw() {
	rl.DrawRectangle(
		tile.TextureRect.ToInt32().X,
		tile.TextureRect.ToInt32().Y,
		tile.TextureRect.ToInt32().Width,
		tile.TextureRect.ToInt32().Height,
		tile.BaseColor,
	)
}

func BasicTestLevelLayout() []Tile {
	var layout []Tile

	// floor first row
	for y := 0; y < 40; y++ {
		tileRect := rl.NewRectangle(float32(8 * y), 172, 8, 8)
		tile := Tile{
			Type:        "grass",
			IsSolid:     true,
			TextureRect: tileRect,
			BaseColor:   grassRegularGreen,
		}

		layout = append(layout, tile)
    }

	// floor second row
	for y := 0; y < 40; y++ {
		tileRect := rl.NewRectangle(float32(8 * y), 164, 8, 8)
		tile := Tile{
			Type:        "grass",
			IsSolid:     true,
			TextureRect: tileRect,
			BaseColor:   grassRegularGreen,
		}

		layout = append(layout, tile)
    }

	// close floating platform
	for y := 0; y < 7; y++ {
		tileRect := rl.NewRectangle(200 + float32(8 * y), 124, 8, 8)
		tile := Tile{
			Type:        "grass",
			IsSolid:     true,
			TextureRect: tileRect,
			BaseColor:   grassRegularGreen,
		}

		layout = append(layout, tile)
    }

	// middle floating platform
	for y := 0; y < 7; y++ {
		tileRect := rl.NewRectangle(110 + float32(8 * y), 105, 8, 8)
		tile := Tile{
			Type:        "grass",
			IsSolid:     true,
			TextureRect: tileRect,
			BaseColor:   grassRegularGreen,
		}

		layout = append(layout, tile)
    }

	// far floating platform
	for y := 0; y < 7; y++ {
		tileRect := rl.NewRectangle(20 + float32(8 * y), 90, 8, 8)
		tile := Tile{
			Type:        "grass",
			IsSolid:     true,
			TextureRect: tileRect,
			BaseColor:   grassRegularGreen,
		}

		layout = append(layout, tile)
    }

	return layout
}
