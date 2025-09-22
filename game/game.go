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
	DEBUG                     bool = false
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
	DebugMode               bool
}

type World struct {
	ID     string   `json:"iid"`
	Levels []*Level `json:"levels"`
}

func InitGame(debugMode bool) *Game {
	world, loadWorldErr := LoadWorld()
	if loadWorldErr != nil {
		panic(fmt.Sprintf("error loading world: %s", loadWorldErr.Error()))
	}

	// TODO: the keys from the map can probably be infered from the assets file data
	tilesetGroundImage := rl.LoadImageFromMemory(".png", assets.GROUND_SPRITE_DATA, int32(len(assets.GROUND_SPRITE_DATA)))
	backgroundDaylightSky := rl.LoadImageFromMemory(".png", assets.BACKRGOUND_DAYLIGHT_SKY, int32(len(assets.BACKRGOUND_DAYLIGHT_SKY)))
	backgroundUnderground := rl.LoadImageFromMemory(".png", assets.BACKRGOUND_UNDERGROUND, int32(len(assets.BACKRGOUND_UNDERGROUND)))
	var textures = map[string]rl.Texture2D{
		"tileset_ground": rl.LoadTextureFromImage(tilesetGroundImage),
		"background-daylight-sky": rl.LoadTextureFromImage(backgroundDaylightSky),
		"background-underground": rl.LoadTextureFromImage(backgroundUnderground),
	}

	renderer := &Renderer{
		Textures: textures,
		DebugMode: debugMode,
	}

	pc := InitPlayer()

	game := Game{
		Player: pc,
		State: Playing,
		World: world,
		Renderer: renderer,
		DebugMode: debugMode,
	}

	game.LoadLevel("Level_1")

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

	g.CheckRoomChange()

	g.Player.Tick(delta, g.CurrentLevel)
	g.CurrentLevel.Tick(delta)

	g.Render()
	g.IncreaseFrameCount()

	if g.DebugMode {
		g.LogState()
	}
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

	if g.DebugMode {
		g.Player.DrawHitbox()
	}

	g.CurrentLevel.DrawParticles(g.Renderer)
}

func (g *Game) CheckRoomChange() {
	var direction string

	if g.Player.WentNorth {
		direction = "n"
	}

	if g.Player.WentEast {
		direction = "e"
	}

	if g.Player.WentSouth {
		direction = "s"
	}

	if g.Player.WentWest {
		direction = "w"
	}

	for _, neighbour := range g.CurrentLevel.Neighbours {
		if neighbour.Direction == direction {
			levelName := g.FindLevelNameFromID(neighbour.LevelID)

			g.CurrentLevel.Unload()
			g.LoadLevel(levelName)

			break
		}
	}

	g.Player.WentNorth = false
	g.Player.WentEast = false
	g.Player.WentSouth = false
	g.Player.WentWest = false
}

func (game *Game) LogState() {
	moveLeft := rl.IsKeyDown(rl.KeyA)
	moveRight := rl.IsKeyDown(rl.KeyD)
	jump := rl.IsKeyPressed(rl.KeySpace)

	rl.TraceLog(rl.LogInfo, "=======")
	rl.TraceLog(rl.LogInfo, "frame: %d", game.AbsoluteFrame)
	rl.TraceLog(rl.LogInfo, "player.Velocity: %f", game.Player.Velocity)
	// rl.TraceLog(rl.LogInfo, "player.FramesCounter: %d", game.Player.FramesCounter)
	rl.TraceLog(rl.LogInfo, "player.Position.X: %f", game.Player.Position.X)
	rl.TraceLog(rl.LogInfo, "player.Position.Y: %f", game.Player.Position.Y)
	rl.TraceLog(rl.LogInfo, "player.WentNorth: %t", game.Player.WentNorth)
	rl.TraceLog(rl.LogInfo, "player.WentEast: %t", game.Player.WentEast)
	rl.TraceLog(rl.LogInfo, "player.WentSouth: %t", game.Player.WentSouth)
	rl.TraceLog(rl.LogInfo, "player.WentWest: %t", game.Player.WentWest)
	rl.TraceLog(rl.LogInfo, "input.moveLeft: %v", moveLeft)
	rl.TraceLog(rl.LogInfo, "input.moveRight: %v", moveRight)
	rl.TraceLog(rl.LogInfo, "input.jump: %v", jump)
	rl.TraceLog(rl.LogInfo, "game.State: %s", game.State)
}

func (game *Game) IncreaseFrameCount() {
	game.AbsoluteFrame += 1
}
