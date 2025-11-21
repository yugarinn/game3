package main

import (
	"flag"
	"game3/game"
	"game3/ui"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	VIRTUAL_WINDOW_WIDTH  int32  = 320
	VIRTUAL_WINDOW_HEIGHT int32  = 180
	GAME_TITLE            string = "game3"
	FLOOR_TILE_SIZE       int32  = 16
	UI_HEIGHT             int32  = 5
	GRID_WIDTH            int32  = VIRTUAL_WINDOW_WIDTH / FLOOR_TILE_SIZE
	GRID_HEIGHT           int32  = VIRTUAL_WINDOW_HEIGHT / (FLOOR_TILE_SIZE + UI_HEIGHT)
	TARGET_FPS            int32  = 60
)

var (
	virtualScreen rl.RenderTexture2D
	scale         float32 = 1.0
	offset        rl.Vector2
)

func main() {
	slowMotionScale := 1
	rl.InitWindow(VIRTUAL_WINDOW_WIDTH*3, VIRTUAL_WINDOW_HEIGHT*3, GAME_TITLE)
	defer rl.CloseWindow()

	virtualScreen = rl.LoadRenderTexture(VIRTUAL_WINDOW_WIDTH, VIRTUAL_WINDOW_HEIGHT)
	defer rl.UnloadRenderTexture(virtualScreen)

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetTargetFPS(TARGET_FPS)
	rl.SetExitKey(0)

	debugMode := flag.Bool("debug", false, "init the game in debug mode")
	flag.Parse()

	instance := game.InitGame(*debugMode)

	for !rl.WindowShouldClose() {
		updateScreenScale()

		delta := rl.GetFrameTime()

		rl.BeginTextureMode(virtualScreen)
		{
			rl.ClearBackground(rl.Black)

			if rl.IsKeyPressed(rl.KeyEscape) && instance.State == game.Playing {
				instance.State = game.Paused
			}

			if rl.IsKeyPressed(rl.KeyEscape) && instance.State == game.Paused {
				instance.State = game.Playing
			}

			if rl.IsKeyPressed(rl.KeyR) {
				instance.Player.Position = rl.NewVector2(60, 60)
				instance.Player.Velocity.Y = 0
				instance.LoadLevel("Level_1")
			}

			if instance.State == game.MainMenu {
				ui.ShowMainMenu(instance)
			}

			if instance.State == game.Playing {
				instance.Tick(delta / float32(slowMotionScale))
				time.Sleep(time.Millisecond * time.Duration(slowMotionScale))
			}
		}
		rl.EndTextureMode()

		rl.BeginDrawing()
		{
			projectVirtualScreenToWindow()
		}
		rl.EndDrawing()
	}
}

func updateScreenScale() {
	windowWidth := rl.GetScreenWidth()
	windowHeight := rl.GetScreenHeight()

	scaleX := float32(windowWidth) / float32(VIRTUAL_WINDOW_WIDTH)
	scaleY := float32(windowHeight) / float32(VIRTUAL_WINDOW_HEIGHT)

	if scaleX < scaleY {
		scale = scaleX
	} else {
		scale = scaleY
	}

	offset.X = (float32(windowWidth) - (float32(VIRTUAL_WINDOW_WIDTH) * scale)) * 0.5
	offset.Y = (float32(windowHeight) - (float32(VIRTUAL_WINDOW_HEIGHT) * scale)) * 0.5

	rl.SetMouseScale(1/scale, 1/scale)
}

func projectVirtualScreenToWindow() {
	rl.ClearBackground(rl.Black)

	source := rl.NewRectangle(0, 0, float32(VIRTUAL_WINDOW_WIDTH), float32(-VIRTUAL_WINDOW_HEIGHT))
	destination := rl.NewRectangle(
		offset.X,
		offset.Y,
		float32(VIRTUAL_WINDOW_WIDTH)*scale,
		float32(VIRTUAL_WINDOW_HEIGHT)*scale,
	)
	origin := rl.NewVector2(0, 0)

	rl.DrawTexturePro(
		virtualScreen.Texture,
		source,
		destination,
		origin,
		0.0,
		rl.White,
	)
}
