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
	Renderer                *Renderer

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

	renderer := &Renderer{
		Textures: textures,
	}

	pc := InitPlayer()

	game := Game{
		Player: pc,
		State: Playing,
		World: world,
		Renderer: renderer,
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

	// TODO: remove duplication
	if g.Player.WentNorth {
		for _, neighbour := range g.CurrentLevel.Neighbours {
			if neighbour.Direction == "n" {
				levelName := g.FindLevelNameFromID(neighbour.LevelID)

				g.LoadLevel(levelName)
		        g.Player.WentNorth = false

				break
			}
		}
	}

	if g.Player.WentEast {
		for _, neighbour := range g.CurrentLevel.Neighbours {
			if neighbour.Direction == "e" {
				levelName := g.FindLevelNameFromID(neighbour.LevelID)

				g.LoadLevel(levelName)
		        g.Player.WentEast = false

				break
			}
		}
	}

	if g.Player.WentSouth {
		for _, neighbour := range g.CurrentLevel.Neighbours {
			if neighbour.Direction == "s" {
				levelName := g.FindLevelNameFromID(neighbour.LevelID)

				g.LoadLevel(levelName)
		        g.Player.WentSouth = false

				break
			}
		}
	}

	if g.Player.WentWest {
		for _, neighbour := range g.CurrentLevel.Neighbours {
			if neighbour.Direction == "w" {
				levelName := g.FindLevelNameFromID(neighbour.LevelID)

				g.LoadLevel(levelName)
		        g.Player.WentWest = false

				break
			}
		}
	}

	g.Player.Tick(delta, g.CurrentLevel)

	g.Render()
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

func (g *Game) LoadLevel(levelName string) {
	var currentLevel *Level
	for _, level := range g.World.Levels {
		if level.Name == levelName {
			currentLevel = level
			break
		}
	}

	currentLevel.Load()
	g.CurrentLevel = currentLevel
}

func (g *Game) FindLevelNameFromID(levelID string) string {
	for _, level := range g.World.Levels {
		if level.ID  == levelID {
			return level.Name
		}
	}

	return ""
}

func (g *Game) Render() {
	g.CurrentLevel.Draw(g.Renderer)
	g.Player.Draw()
}

func (game *Game) LogState() {
	moveLeft := rl.IsKeyDown(rl.KeyA)
	moveRight := rl.IsKeyDown(rl.KeyD)
	jump := rl.IsKeyPressed(rl.KeySpace)

	rl.TraceLog(rl.LogInfo, "=======")
	rl.TraceLog(rl.LogInfo, "frame: %d", game.AbsoluteFrame)
	rl.TraceLog(rl.LogInfo, "player.Velocity: %f", game.Player.Velocity)
	rl.TraceLog(rl.LogInfo, "player.FramesCounter: %d", game.Player.FramesCounter)
	rl.TraceLog(rl.LogInfo, "player.Position.X: %f", game.Player.Position.X)
	rl.TraceLog(rl.LogInfo, "player.Position.Y: %f", game.Player.Position.Y)
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
