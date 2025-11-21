package ui

import (
	"game3/game"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	VIRTUAL_WINDOW_WIDTH  int = 320
	VIRTUAL_WINDOW_HEIGHT int = 180
)

var greenishBlack rl.Color = rl.NewColor(88, 68, 34, 255)
var dirtyYellow rl.Color = rl.NewColor(212, 210, 155, 255)
var regularGreen rl.Color = rl.NewColor(94, 133, 73, 255)
var regularGreenHover rl.Color = rl.NewColor(100, 200, 100, 255)

func ShowMainMenu(instance *game.Game) {
	menu := NewUiElement(NewUiElementInput{
		Width:           float32(VIRTUAL_WINDOW_WIDTH),
		Height:          float32(VIRTUAL_WINDOW_HEIGHT),
		BackgroundColor: greenishBlack,
		BorderColor:     regularGreen,
		BorderWidth:     2,
		HPosition:       HCentered,
		VPosition:       VCentered,
	})

	resumeButton := NewUiElement(NewUiElementInput{
		Width:           100,
		Height:          20,
		BackgroundColor: greenishBlack,
		BorderColor:     regularGreen,
		BorderWidth:     1,
		HPosition:       Top,
		VPosition:       VCentered,
		Margin:          UiMargin{Top: 20},
		Text:            "Resume",
	})

	optionsButton := NewUiElement(NewUiElementInput{
		Width:           100,
		Height:          20,
		BackgroundColor: greenishBlack,
		BorderColor:     regularGreen,
		BorderWidth:     1,
		HPosition:       Top,
		VPosition:       VCentered,
		Margin:          UiMargin{Top: 50},
		Text:            "Options",
	})

	quitButton := NewUiElement(NewUiElementInput{
		Width:           100,
		Height:          20,
		BackgroundColor: greenishBlack,
		BorderColor:     regularGreen,
		BorderWidth:     1,
		HPosition:       Bottom,
		VPosition:       VCentered,
		Margin:          UiMargin{Bottom: 10},
		Text:            "Quit",
	})

	resumeButton.AddEventListener("click", func() {
		instance.SetState(game.Playing)
	})

	resumeButton.AddEventListener("hover", func() {
		resumeButton.SetBackgroundColor(dirtyYellow)
	})

	optionsButton.AddEventListener("click", func() {
		rl.TraceLog(rl.LogInfo, "options button clicked")
	})

	optionsButton.AddEventListener("hover", func() {
		optionsButton.SetBackgroundColor(dirtyYellow)
	})

	quitButton.AddEventListener("click", func() {
		rl.TraceLog(rl.LogInfo, "quit button clicked")
	})

	quitButton.AddEventListener("hover", func() {
		quitButton.SetBackgroundColor(dirtyYellow)
	})

	menu.AddChild(&resumeButton)
	menu.AddChild(&optionsButton)
	menu.AddChild(&quitButton)

	menu.Tick()
}
