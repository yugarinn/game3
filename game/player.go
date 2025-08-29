package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"game3/assets"
)

const (
	PLAYER_WIDTH        float32 = 24
	PLAYER_HEIGHT       float32 = 24
	PLAYER_MOVE_SPEED   float32 = 100
	PLAYER_ACCELERATION float32 = 500
	PLAYER_DECELERATION float32 = 700
	PLAYER_JUMP_FORCE   float32 = -280
)

type Player struct {
	Position          rl.Vector2
	Velocity          rl.Vector2
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
	LastStateChangeAt float64
	AttackTime        float64
	AttackDirection   string
	OnGround          bool
	IsRunning         bool
	IsJumping         bool
	IsFalling         bool
	CanJump           bool
}

func InitPlayer() *Player {
	playerImage := rl.LoadImageFromMemory(".png", assets.PLAYER_SPRITE_DATA, int32(len(assets.PLAYER_SPRITE_DATA)))

	player := Player{
		Position:        rl.NewVector2(10, 140),
		Velocity:        rl.NewVector2(0, 0),
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
		OnGround:        false,
		CanJump:         true,
		IsJumping:       false,
		IsFalling:       false,
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
	// rl.DrawRectangleLinesEx(player.HitboxRect, 1, rl.NewColor(230, 41, 55, 100))
	// rl.DrawPixel(int32(player.Position.X), int32(player.Position.Y), rl.Green)
}

func (player *Player) Tick(delta float32, level *Level) {
	player.ProcessInput(delta)
	player.CalculateVelocity(delta)
	player.UpdatePosition(delta, level)
	player.UpdateState()
	player.UpdateAnimation()
}

func (player *Player) ProcessInput(delta float32) {
	moveLeft := rl.IsKeyDown(rl.KeyA) || rl.IsGamepadButtonDown(1, rl.GamepadButtonLeftFaceLeft)
	moveRight := rl.IsKeyDown(rl.KeyD) || rl.IsGamepadButtonDown(1, rl.GamepadButtonLeftFaceRight)
	jump := rl.IsKeyDown(rl.KeySpace)
	jumpReleased := rl.IsKeyReleased(rl.KeySpace)
	reset := rl.IsKeyPressed(rl.KeyR)

	// prevents spamming jumps
	if jumpReleased {
		player.CanJump = true
	}

	if reset {
		player.Position = rl.NewVector2(10, 10)
		player.Velocity.Y = 0
		return
	}

	if moveLeft && !moveRight {
		player.FacingDirection = "LEFT"

		if player.Velocity.X > -PLAYER_MOVE_SPEED {
			player.Velocity.X -= PLAYER_ACCELERATION * delta

			if player.Velocity.X < -PLAYER_MOVE_SPEED {
				player.Velocity.X = -PLAYER_MOVE_SPEED
			}
		}
	}

	if moveRight && !moveLeft {
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

	if jump && player.OnGround && player.CanJump {
		player.Velocity.Y = PLAYER_JUMP_FORCE
		player.OnGround = false
		player.CanJump = false
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

func (player *Player) UpdatePosition(delta float32, level *Level) {
	player.Position.X += player.Velocity.X * delta
	player.Position.Y += player.Velocity.Y * delta

	player.UpdateHitbox()
	player.HandleTileCollisions(level.GetGroundLayer().Layout)
}

func (player *Player) UpdateState() {
	var isRunning bool
	var isJumping bool
	var isFalling bool

	if player.Velocity.X != 0 {
		isRunning = true
		player.FramesSpeed = 6
	}

	if player.Velocity.X == 0 {
		isRunning = false
		player.FramesSpeed = 2
	}

	if player.Velocity.Y < 0 {
		isJumping = true
	}

	if player.Velocity.Y > 0 {
		isFalling = true
	}

	player.IsRunning = isRunning
	player.IsJumping = isJumping
	player.IsFalling = isFalling
}

func (player *Player) UpdateHitbox() {
	player.HitboxRect.X = player.Position.X
	player.HitboxRect.Y = player.Position.Y
}

func (player *Player) UpdateAnimation() {
	player.FramesCounter++

	if player.FramesCounter >= (60 / player.FramesSpeed) {
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

		if player.IsJumping {
			player.TextureRect.X = 16
			player.TextureRect.Y = 80

			return
		}

		if player.IsFalling {
			player.TextureRect.X = 48
			player.TextureRect.Y = 80

			return
		}

		offset := float32(player.CurrentFrame) * 16
		player.TextureRect.X = offset
	}
}

func (player *Player) HandleTileCollisions(layout []*Tile) {
	player.OnGround = false

	for _, tile := range layout {
		if rl.CheckCollisionRecs(tile.HitboxRect, player.HitboxRect) {
			overlapLeft := (player.HitboxRect.X + player.HitboxRect.Width) - tile.HitboxRect.X
			overlapRight := (tile.HitboxRect.X + tile.HitboxRect.Width) - player.HitboxRect.X
			overlapTop := (player.HitboxRect.Y + player.HitboxRect.Height) - tile.Position.Y
			overlapBottom := (tile.Position.Y + tile.HitboxRect.Height) - player.HitboxRect.Y

			minOverlap := overlapLeft
			collisionSide := "LEFT"

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

			// This is my poor man's solution to ghost colliding
			// https://briansemrau.github.io/dealing-with-ghost-collisions/
			if (collisionSide == "LEFT" || collisionSide == "RIGHT") && overlapBottom > 0 && overlapBottom > overlapTop {
				continue
			}

			switch collisionSide {
			case "TOP":
				player.Position.Y = tile.Position.Y - player.HitboxRect.Height
				player.Velocity.Y = 0
				player.OnGround = true
			case "BOTTOM":
				player.Position.Y = tile.Position.Y + tile.HitboxRect.Height
				player.Velocity.Y = 0
			case "LEFT":
				player.Position.X = tile.Position.X - player.HitboxRect.Width
				player.Velocity.X = 0
			case "RIGHT":
				player.Position.X = tile.Position.X + tile.HitboxRect.Width
				player.Velocity.X = 0
			}
		}
	}
}
