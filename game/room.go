package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Room struct {
	ID     string
	Name   string
	Layout []Tile
	Props  []Prop
}

var grassRegularGreen rl.Color = rl.NewColor(99, 179, 140, 255)
var dawnUpperColor rl.Color = rl.NewColor(86, 85, 149, 255)
var dawnLowerColor rl.Color = rl.NewColor(239, 179, 165, 255)

func LoadRoom(id string) (*Room, error) {
	return &Room{
		ID:     "test_room",
		Name:   "Test Room",
		Layout: BasicTestLevelLayout(),
	}, nil
}

func (room *Room) Draw() {
	room.DrawBackground()
	room.DrawForeground()
	room.DrawTiles()
	room.DrawParticles()
}

func (room *Room) DrawBackground() {
	for x := 0; x < 40; x++ {
		for y := 0; y < 23; y++ {
			t := float32(y) / float32(23-1)

			tileTopR := uint8(float32(dawnUpperColor.R) + t * (float32(dawnLowerColor.R)-float32(dawnUpperColor.R)))
			tileTopG := uint8(float32(dawnUpperColor.G) + t * (float32(dawnLowerColor.G)-float32(dawnUpperColor.G)))
			tileTopB := uint8(float32(dawnUpperColor.B) + t * (float32(dawnLowerColor.B)-float32(dawnUpperColor.B)))
			tileTopColor := rl.NewColor(tileTopR, tileTopG, tileTopB, 255)

			tNext := float32(y + 1) / float32(23 - 1)
			if tNext > 1.0 {
				tNext = 1.0
			}

			tileBottomR := uint8(float32(dawnUpperColor.R) + tNext * (float32(dawnLowerColor.R)-float32(dawnUpperColor.R)))
			tileBottomG := uint8(float32(dawnUpperColor.G) + tNext * (float32(dawnLowerColor.G)-float32(dawnUpperColor.G)))
			tileBottomB := uint8(float32(dawnUpperColor.B) + tNext * (float32(dawnLowerColor.B)-float32(dawnUpperColor.B)))
			tileBottomColor := rl.NewColor(tileBottomR, tileBottomG, tileBottomB, 255)

			rl.DrawRectangleGradientV(int32(x * 8), int32(y * 8), 8, 8, tileTopColor, tileBottomColor)
		}
	}
}

func (room *Room) DrawForeground() {

}

func (room *Room) DrawParticles() {

}

func (room *Room) DrawTiles() {
	for _, tile := range(room.Layout) {
		tile.Draw()
	}
}
