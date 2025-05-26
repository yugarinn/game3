package ui

import (
	"game3/game"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var greenishBlack     rl.Color = rl.NewColor(88, 68, 34, 255)
var dirtyYellow       rl.Color = rl.NewColor(212, 210, 155, 255)
var regularGreen      rl.Color = rl.NewColor(94, 133, 73, 255)
var regularGreenHover rl.Color = rl.NewColor(100, 200, 100, 255)

func ShowMainMenu(game *game.Game) {
	container := NewUiElement(NewUiElementInput{
		Width: 250,
		Height: 260,
		BackgroundColor: greenishBlack,
		BorderColor: regularGreen,
		BorderWidth: 2,
		HPosition: HCentered,
		VPosition: VCentered,
	})

	container.Tick()
}
