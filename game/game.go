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

const (
	GRAVITY                float32 = 1000
	FALL_TERMINAL_VELOCITY float32 = 500
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
	Player                  *Player
	CurrentRoom             *Room
}

func InitGame() *Game {
	room, _ := LoadRoom("id")
	pc := InitPlayer()

	game := Game{
		Player: pc,
		State:  Playing,
		CurrentRoom: room,
	}

	return &game
}

func (game *Game) SetState(state GameState) {
	if state == Playing {
		room, _ := LoadRoom("id")
		game.CurrentRoom = room
	}

	game.State = state
}

func (game *Game) Tick(delta float32) {
	if game.State == MainMenu {
		return
	}

	game.Player.Tick(delta, game.CurrentRoom)

	game.CurrentRoom.Draw()
	game.Player.Draw()

	game.IncreaseFrameCount()
	game.LogState()
}

func (game *Game) LogState() {
	moveLeft := rl.IsKeyDown(rl.KeyA)
	moveRight := rl.IsKeyDown(rl.KeyD)
	jump := rl.IsKeyPressed(rl.KeySpace)

	rl.TraceLog(rl.LogInfo, "frame: %d", game.AbsoluteFrame)
	rl.TraceLog(rl.LogInfo, "player.Velocity: %f", game.Player.Velocity)
	rl.TraceLog(rl.LogInfo, "input.moveLeft: %v", moveLeft)
	rl.TraceLog(rl.LogInfo, "input.moveRight: %v", moveRight)
	rl.TraceLog(rl.LogInfo, "input.jump: %v", jump)
}

func (game *Game) IncreaseFrameCount() {
	game.AbsoluteFrame += 1
}
