package game

import (
	"encoding/json"
	"fmt"
	"os"

	"game3/assets"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState int

const (
	MainMenu GameState = iota
	Playing
	TimeStop
	Paused
	Editing
)

const (
	GRAVITY                float32 = 900
	FALL_TERMINAL_VELOCITY float32 = 600
)

var gameStateName = map[GameState]string {
	MainMenu: "MainMenu",
	Playing:  "Playing",
	TimeStop: "TimeStop",
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
	World                   *World
	CurrentLevel            *Level
	Textures                map[string]rl.Texture2D
}

type World struct {
	ID     string   `json:"iid"`
	Levels []*Level `json:"levels"`
}

func InitGame() *Game {
	world, loadWorldErr := LoadWorld()
	if loadWorldErr != nil {
		panic(fmt.Sprintf("error loading world: %s", loadWorldErr.Error()))
	}

	groundImage := rl.LoadImageFromMemory(".png", assets.GROUND_SPRITE_DATA, int32(len(assets.GROUND_SPRITE_DATA)))
	var textures = map[string]rl.Texture2D{
		"ground": rl.LoadTextureFromImage(groundImage),
	}

	pc := InitPlayer()
	game := Game{
		Player: pc,
		State: Playing,
		World: world,
		Textures: textures,
	}

	game.LoadLevel("Level_0")

	return &game
}

func (game *Game) SetState(state GameState) {
	if state == Playing {
		game.LoadLevel("Level_0")
	}

	game.State = state
}

func (g *Game) Tick(delta float32) {
	if g.State == MainMenu {
		return
	}

	g.Player.Tick(delta, g.CurrentLevel)

	g.CurrentLevel.Draw()
	g.Player.Draw()

	g.IncreaseFrameCount()
	g.LogState()
}

func LoadWorld() (*World, error) {
	worldFile, jsonErr := os.ReadFile("levels/game3.ldtk")
	if jsonErr != nil {
		return nil, jsonErr
	}

	var world World
	json.Unmarshal(worldFile, &world)

	return &world, nil
}

func (g *Game) LoadLevel(levelID string) {
	var currentLevel *Level
	for _, level := range g.World.Levels {
		if level.Name == levelID {
			currentLevel = level
			break
		}
	}

	currentLevel.Load()
	g.CurrentLevel = currentLevel
}

func (game *Game) LogState() {
	moveLeft := rl.IsKeyDown(rl.KeyA)
	moveRight := rl.IsKeyDown(rl.KeyD)
	jump := rl.IsKeyPressed(rl.KeySpace)

	rl.TraceLog(rl.LogInfo, "=======")
	rl.TraceLog(rl.LogInfo, "frame: %d", game.AbsoluteFrame)
	rl.TraceLog(rl.LogInfo, "player.Velocity: %f", game.Player.Velocity)
	rl.TraceLog(rl.LogInfo, "player.FramesCounter: %d", game.Player.FramesCounter)
	rl.TraceLog(rl.LogInfo, "input.moveLeft: %v", moveLeft)
	rl.TraceLog(rl.LogInfo, "input.moveRight: %v", moveRight)
	rl.TraceLog(rl.LogInfo, "input.jump: %v", jump)
	rl.TraceLog(rl.LogInfo, "input.buttonPressed: %v", rl.GetGamepadButtonPressed())
	rl.TraceLog(rl.LogInfo, "input.isGamePad0Present: %v, %v", rl.IsGamepadAvailable(0), rl.GetGamepadName(0))
	rl.TraceLog(rl.LogInfo, "input.isGamePad1Present: %v, %v", rl.IsGamepadAvailable(1), rl.GetGamepadName(1))
	rl.TraceLog(rl.LogInfo, "input.isGamePad2Present: %v", rl.IsGamepadAvailable(2))
	rl.TraceLog(rl.LogInfo, "input.isGamePad3Present: %v", rl.IsGamepadAvailable(3))
	rl.TraceLog(rl.LogInfo, "input.isGamePad4Present: %v", rl.IsGamepadAvailable(4))
	rl.TraceLog(rl.LogInfo, "game.State: %s", game.State)
}

func (game *Game) IncreaseFrameCount() {
	game.AbsoluteFrame += 1
}
