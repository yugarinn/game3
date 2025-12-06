package game

import (
	"encoding/json"
	"fmt"

	"game3/assets"
	"game3/levels"

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
	GRAVITY                float32 = 50
	FALL_TERMINAL_VELOCITY float32 = 250
	DEBUG                  bool    = false
)

var gameStateName = map[GameState]string{
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
	FrameInspectorMode      bool
	ActiveGamepad           int32
	CurrentVFXs             []*VFX
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
	tilemap := rl.LoadImageFromMemory(".png", assets.TILEMAP, int32(len(assets.TILEMAP)))
	var textures = map[string]rl.Texture2D{
		"tilemap": rl.LoadTextureFromImage(tilemap),
	}

	renderer := &Renderer{
		Textures:  textures,
		DebugMode: debugMode,
	}

	pc := InitPlayer()

	game := Game{
		Player:    pc,
		State:     Playing,
		World:     world,
		Renderer:  renderer,
		DebugMode: debugMode,
	}

	game.LoadLevel("Level_4")

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

	g.Render()
	g.IncreaseFrameCount()

	g.DetectActiveGamepad()
	g.ProcessInput()
	g.CheckRoomChange()

	if g.CanTick() {
		if g.Player.LastAction == Jump {
			g.PlayVFX(PlayerJumpVFX, g.Player.Position)
		}

		if g.Player.IsDead {
			g.PlayVFX(PlayerDeathVFX, g.Player.Position)
			g.Reset()
			g.Player.IsDead = false
		}

		g.Player.Tick(delta, g.CurrentLevel, g.ActiveGamepad)
		g.CurrentLevel.Tick(delta)
	}

	if g.DebugMode {
		g.LogState()
	}
}

func (g *Game) ProcessInput() {
	if rl.IsKeyReleased(rl.KeyI) {
		g.FrameInspectorMode = !g.FrameInspectorMode
	}
}

func LoadWorld() (*World, error) {
	var world World
	json.Unmarshal(levels.LEVELS, &world)

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
		if level.ID == levelID {
			return level.Name
		}
	}

	return ""
}

func (g *Game) Render() {
	g.CurrentLevel.DrawLayer("Background", g.Renderer)
	g.CurrentLevel.DrawLayer("BackgroundProps", g.Renderer)
	g.CurrentLevel.DrawLayer("Ground", g.Renderer)
	g.CurrentLevel.DrawProps(g.Renderer)
	g.Player.Draw(g.Renderer)
	g.DrawCurrentVFXs()
	g.CurrentLevel.DrawParticles(g.Renderer)
	g.CurrentLevel.DrawLayer("ForegroundProps", g.Renderer)

	if g.DebugMode {
		g.Player.DrawHitbox()
	}
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
			g.Player.Path = make([]rl.Vector2, 20)

			break
		}
	}

	g.Player.WentNorth = false
	g.Player.WentEast = false
	g.Player.WentSouth = false
	g.Player.WentWest = false
}

func (game *Game) LogState() {
	// rl.TraceLog(rl.LogInfo, "=======")
	// rl.TraceLog(rl.LogInfo, "frame: %d", game.AbsoluteFrame)
	rl.TraceLog(rl.LogInfo, "player.Velocity: %f", game.Player.Velocity)
	rl.TraceLog(rl.LogInfo, "collisionable.Position: %f", game.CurrentLevel.CollisionableHitboxes[game.CurrentLevel.PlayerCollisionIndex])
	// rl.TraceLog(rl.LogInfo, "player.FramesCounter: %d", game.Player.FramesCounter)
	// rl.TraceLog(rl.LogInfo, "player.Position.X: %f", game.Player.Position.X)
	// rl.TraceLog(rl.LogInfo, "player.Position.Y: %f", game.Player.Position.Y)
	// rl.TraceLog(rl.LogInfo, "player.WentNorth: %t", game.Player.WentNorth)
	// rl.TraceLog(rl.LogInfo, "player.WentEast: %t", game.Player.WentEast)
	// rl.TraceLog(rl.LogInfo, "player.WentSouth: %t", game.Player.WentSouth)
	// rl.TraceLog(rl.LogInfo, "player.WentWest: %t", game.Player.WentWest)
	// rl.TraceLog(rl.LogInfo, "input.moveLeft: %v", moveLeft)
	// rl.TraceLog(rl.LogInfo, "input.moveRight: %v", moveRight)
	// rl.TraceLog(rl.LogInfo, "input.jump: %v", jump)
	// rl.TraceLog(rl.LogInfo, "game.State: %s", game.State)

	// for _, prop := range game.CurrentLevel.Props {
	// 	rl.TraceLog(rl.LogInfo, "level.Props: %#v", prop)
	// }

	/// rl.TraceLog(rl.LogInfo, "IsGamepadAvailable: %t", rl.IsGamepadAvailable(0))
	/// rl.TraceLog(rl.LogInfo, "IsGamepadAvailable: %t", rl.IsGamepadAvailable(0))
	/// rl.TraceLog(rl.LogInfo, "GamepadName: %s", rl.GetGamepadName(0))
	/// rl.TraceLog(rl.LogInfo, "IsGamepadButtonDown: %t", rl.IsGamepadButtonDown(1, rl.GamepadButtonRightFaceDown))
	/// rl.TraceLog(rl.LogInfo, "FrameInspectorMode: %t", game.FrameInspectorMode)
}

func (g *Game) IncreaseFrameCount() {
	g.AbsoluteFrame += 1
}

func (g *Game) CanTick() bool {
	if g.FrameInspectorMode && rl.IsKeyReleased(rl.KeyN) {
		return true
	}

	if g.FrameInspectorMode {
		return false
	}

	return true
}

func (game *Game) Reset() {
	game.LoadLevel("Level_2")
	game.Player.Position = rl.NewVector2(147, 82)
	game.Player.Velocity.Y = 0
}

func (game *Game) DetectActiveGamepad() {
	for i := range int32(4) {
		if rl.IsGamepadAvailable(i) {
			game.ActiveGamepad = i
			return
		}
	}
}

func (game *Game) PlayVFX(vfxType VFXType, position rl.Vector2) {
	vfx := NewVFX(vfxType, position)
	game.CurrentVFXs = append(game.CurrentVFXs, &vfx)
}

func (game *Game) DrawCurrentVFXs() {
	for i, vfx := range game.CurrentVFXs {
		game.Renderer.DrawVFX(vfx)
		vfx.AnimationCurrentFrame += 1

		if vfx.AnimationCurrentFrame%vfx.AnimationFramesPerPosition == 0 {
			vfx.AnimationCurrentPosition += 1
			vfx.AnimationCurrentFrame = 0
		}

		if vfx.AnimationCurrentPosition >= vfx.AnimationPositionsCounter {
			if vfx.Loops {
				vfx.AnimationCurrentPosition = 0
			} else {
				game.CurrentVFXs = append(game.CurrentVFXs[:i], game.CurrentVFXs[i+1:]...)
			}
		}
	}
}
