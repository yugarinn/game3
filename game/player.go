package game

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"

	"game3/assets"
)

const (
	PLAYER_WIDTH        float32 = 24
	PLAYER_HEIGHT       float32 = 24
	PLAYER_MOVE_SPEED   float32 = 120
	PLAYER_ACCELERATION float32 = 700
	PLAYER_DECELERATION float32 = 1100
	PLAYER_JUMP_FORCE   float32 = -300
)

type Player struct {
	Position          rl.Vector2
	Velocity          rl.Vector2
	OnGround          bool
	MaxHealth         int8
	Health            int8
	Sprite            rl.Texture2D
	TextureRect       rl.Rectangle
	HitboxRect        rl.Rectangle
	CurrentFrame      int32
	FramesCounter     int32
	FramesSpeed       int32
	FacingDirection   string
	State             string
	IsRunning         bool
	LastStateChangeAt float64
	AttackTime        float64
	AttackDirection   string
}

func InitPlayer() *Player {
	playerImage := rl.LoadImageFromMemory(".png", assets.PLAYER_SPRITE_DATA, int32(len(assets.PLAYER_SPRITE_DATA)))

	player := Player{
		Position:        rl.NewVector2(10, 100),
		Velocity:        rl.NewVector2(0, 0),
		OnGround:        false,
		MaxHealth:       20,
		Health:          20,
		Sprite:          rl.LoadTextureFromImage(playerImage),
		TextureRect:     rl.NewRectangle(0, 0, 16, 17),
		CurrentFrame:    0,
		FramesCounter:   0,
		FramesSpeed:     2,
		FacingDirection: "RIGHT",
		State:           "IDLE",
		AttackTime:      0.2,
	}

	player.HitboxRect = rl.NewRectangle(player.Position.X, player.Position.Y, 9, 17)

	return &player
}

func (player *Player) Draw() {
	var spriteVector rl.Vector2

	if player.FacingDirection == "RIGHT" {
		spriteVector = player.Position
		player.TextureRect.Width = 16
	}

	if player.FacingDirection == "LEFT" {
		player.TextureRect.Width = -16
		// weird hack, the recommended way to flip a texture in raylib (negating the width) offsets it...
		spriteVector = rl.NewVector2(player.Position.X - 6, player.Position.Y)
	}

	rl.DrawTextureRec(player.Sprite, player.TextureRect, spriteVector, rl.White)
	rl.DrawRectangleLinesEx(player.HitboxRect, 1, rl.NewColor(230, 41, 55, 100))
	rl.DrawPixel(int32(player.Position.X), int32(player.Position.Y), rl.Green)
}

func (player *Player) Tick(delta float32, room *Room) {
	player.ProcessInput(delta)
	player.CalculateVelocity(delta)
	player.UpdateState()
	player.UpdatePosition(delta, room)
	player.UpdateAnimation()
}

func (player *Player) ProcessInput(delta float32) {
	moveLeft := rl.IsKeyDown(rl.KeyA)
	moveRight := rl.IsKeyDown(rl.KeyD)
	jump := rl.IsKeyDown(rl.KeySpace)
	reset := rl.IsKeyPressed(rl.KeyR)

	if reset {
		player.Position = rl.NewVector2(10, 10)
		return
	}

	if moveLeft && !moveRight {
		if player.Velocity.X > 0 {
			player.Velocity.X = 5
		}

		player.FacingDirection = "LEFT"

		if player.Velocity.X > -PLAYER_MOVE_SPEED {
			player.Velocity.X -= PLAYER_ACCELERATION * delta

			if player.Velocity.X < -PLAYER_MOVE_SPEED {
				player.Velocity.X = -PLAYER_MOVE_SPEED
			}
		}
	}

	if moveRight && !moveLeft {
		if player.Velocity.X < 0 {
			player.Velocity.X = -5
		}

		player.FacingDirection = "RIGHT"

		if player.Velocity.X < PLAYER_MOVE_SPEED {
			player.Velocity.X += PLAYER_ACCELERATION * delta

			if player.Velocity.X > PLAYER_MOVE_SPEED {
				player.Velocity.X = PLAYER_MOVE_SPEED
			}
		}
	}

	if !moveLeft && !moveRight {
		if player.Velocity.X > 0 {
			player.Velocity.X -= PLAYER_DECELERATION * delta

			if player.Velocity.X < 0 {
				player.Velocity.X = 0
			}
		} else if player.Velocity.X < 0 {
			player.Velocity.X += PLAYER_DECELERATION * delta

			if player.Velocity.X > 0 {
				player.Velocity.X = 0
			}
		}
	}

	if jump && player.OnGround {
		player.Velocity.Y = PLAYER_JUMP_FORCE
		player.OnGround = false
	}
}

func (player *Player) CalculateVelocity(delta float32) {
	if ! player.OnGround {
		player.Velocity.Y += GRAVITY * delta
	}

	if player.Velocity.Y > FALL_TERMINAL_VELOCITY {
		player.Velocity.Y = FALL_TERMINAL_VELOCITY
	}
}

func (player *Player) UpdatePosition(delta float32, room *Room) {
	player.Position.X += player.Velocity.X * delta
	player.Position.Y += player.Velocity.Y * delta

	player.UpdateHitbox()
	player.HandleTileCollisions(room.Layout)
}

func (player *Player) UpdateState() {
	var isRunning bool

	if player.Velocity.X != 0 {
		isRunning = true
		player.FramesSpeed = 6
	}

	if player.Velocity.X == 0 {
		isRunning = false
		player.FramesSpeed = 2
	}

	// This immediately resets the animaton frame
	if isRunning != player.IsRunning {
		player.CurrentFrame = 0
		player.IsRunning = isRunning
	}
}

func (player *Player) UpdateHitbox() {
	player.HitboxRect.X = player.Position.X
	player.HitboxRect.Y = player.Position.Y
}

func (player *Player) UpdateAnimation() {
	player.FramesCounter++

	if player.FramesCounter >= (60 / player.FramesSpeed) {
		log.Println("INSIDE")
		player.FramesCounter = 0
		player.CurrentFrame++

		if !player.IsRunning && player.CurrentFrame > 1 {
			player.CurrentFrame = 0
		}

		if player.IsRunning && player.CurrentFrame > 3 {
			player.CurrentFrame = 0
		}

		if !player.IsRunning {
			player.TextureRect.Y = 0
		}

		if player.IsRunning {
			player.TextureRect.Y = 31
		}

		offset := float32(player.CurrentFrame) * 16
		player.TextureRect.X = offset
	}
}

func (player *Player) HandleTileCollisions(layout []Tile) {
	player.OnGround = false

	for _, tile := range layout {
		if rl.CheckCollisionRecs(tile.TextureRect, player.HitboxRect) {
			overlapLeft := (player.HitboxRect.X + player.HitboxRect.Width) - tile.TextureRect.X
			overlapRight := (tile.TextureRect.X + tile.TextureRect.Width) - player.HitboxRect.X
			overlapTop := (player.HitboxRect.Y + player.HitboxRect.Height) - tile.TextureRect.Y
			overlapBottom := (tile.TextureRect.Y + tile.TextureRect.Height) - player.HitboxRect.Y

			minOverlap := overlapLeft
			collisionSide := "RIGHT"

			if overlapRight < minOverlap {
				minOverlap = overlapRight
				collisionSide = "RIGHT"
			}

			if overlapTop < minOverlap {
				minOverlap = overlapTop
				collisionSide = "TOP"
			}

			if overlapBottom < minOverlap {
				minOverlap = overlapBottom
				collisionSide = "BOTTOM"
			}

			switch collisionSide {
			case "TOP":
				player.Position.Y = tile.TextureRect.Y - player.TextureRect.Height
				player.Velocity.Y = 0
				player.OnGround = true
			case "BOTTOM":
				player.Position.Y = tile.TextureRect.Y + tile.TextureRect.Height
				player.Velocity.Y = 0
			case "LEFT":
				player.Position.X = tile.TextureRect.X - player.TextureRect.Width
				player.Velocity.X = 0
			case "RIGHT":
				player.Position.X = tile.TextureRect.X + tile.TextureRect.Width
				player.Velocity.X = 0
			}

			player.UpdateHitbox()
		}
	}
}
