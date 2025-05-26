package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState int

const (
	MainMenu GameState = iota
	Playing
	Paused
	Editing
)

var gameStateName = map[GameState]string {
	MainMenu: "MainMenu",
	Playing:  "Playing",
	Paused:   "Paused",
	Editing:  "Editing",
}

func (gs GameState) String() string {
	return gameStateName[gs]
}

type Game struct {
	CurrentFrame            uint32
	AbsoluteFrame           uint32
	LastActionAbsoluteFrame uint32
	State                   GameState
}

func InitGame() *Game {
	game := Game{
		State: MainMenu,
	}

	return &game
}

func (game *Game) Tick() {
	rl.TraceLog(rl.LogDebug, "frame: %d", game.AbsoluteFrame)
	game.AbsoluteFrame += 1

	if game.State == MainMenu {
		return
	}
}
