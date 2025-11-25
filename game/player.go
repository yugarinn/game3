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
	PLAYER_JUMP_FORCE   float32 = -160
)

type Player struct {
	Position             rl.Vector2
	Velocity             rl.Vector2
	MaxHealth            int8
	Health               int8
	Sprite               rl.Texture2D
	TextureRect          rl.Rectangle
	HitboxRect           rl.Rectangle
	InteractiveRect      rl.Rectangle
	CurrentFrame         int32
	FramesCounter        int32
	FramesSpeed          int32
	FacingDirection      string
	State                string
	LastStateChangeAt    float64
	OnGround             bool
	IsRunning            bool
	IsJumping            bool
	IsFalling            bool
	CanJump              bool
	WentNorth            bool
	WentWest             bool
	WentSouth            bool
	WentEast             bool
	Inventory            []*Prop
	ActivePropIndex      int
	IsDroppingActiveProp bool
	Path                 []rl.Vector2
}

func InitPlayer() *Player {
	playerImage := rl.LoadImageFromMemory(".png", assets.TILEMAP, int32(len(assets.TILEMAP)))

	player := Player{
		Position:        rl.NewVector2(10, 140),
		Velocity:        rl.NewVector2(0, 0),
		MaxHealth:       20,
		Health:          20,
		Sprite:          rl.LoadTextureFromImage(playerImage),
		TextureRect:     rl.NewRectangle(0, 56, 8, 8),
		CurrentFrame:    0,
		FramesCounter:   0,
		FramesSpeed:     2,
		FacingDirection: "RIGHT",
		State:           "IDLE",
		OnGround:        false,
		CanJump:         true,
		IsJumping:       false,
		IsFalling:       false,
		WentNorth:       false,
		WentWest:        false,
		WentSouth:       false,
		WentEast:        false,
	}

	player.HitboxRect = rl.NewRectangle(player.Position.X, player.Position.Y, 8, 8)

	interactiveRect := rl.Rectangle{
		X:      player.HitboxRect.X - 2,
		Y:      player.HitboxRect.Y - 2,
		Width:  player.HitboxRect.Width + 4,
		Height: player.HitboxRect.Height + 4,
	}
	player.InteractiveRect = interactiveRect

	return &player
}

func (player *Player) Draw(r *Renderer) {
	var spriteVector rl.Vector2

	if player.FacingDirection == "RIGHT" {
		spriteVector = player.Position
		player.TextureRect.Width = 8
	}

	if player.FacingDirection == "LEFT" {
		player.TextureRect.Width = -8

		// weird hack, the recommended way to flip a texture in raylib, negating the width, offsets it...
		spriteVector = rl.NewVector2(player.Position.X-1, player.Position.Y)
	}

	player.DrawInventory(r)
	rl.DrawTextureRec(player.Sprite, player.TextureRect, spriteVector, rl.White)
}

func (player *Player) DrawInventory(r *Renderer) {
	isMoving := rl.Vector2Length(player.Velocity) > 0.1

	for i := range player.Inventory {
		var targetDistance int
		if isMoving {
			targetDistance = 15 * (i + 1)
		} else {
			targetDistance = 5 * (i + 1)
		}

		var targetPos rl.Vector2
		pathIndex := len(player.Path) - targetDistance

		if pathIndex >= 0 && pathIndex < len(player.Path) {
			targetPos = player.Path[pathIndex]
		} else {
			targetPos = player.Position
		}

		if !isMoving && len(player.Path) > 0 {
			lastGroundPos := player.Path[len(player.Path)-1]
			offsetX := float32(-(i + 1) * 10)
			targetPos = rl.Vector2{
				X: lastGroundPos.X + offsetX,
				Y: lastGroundPos.Y,
			}
		}

		smoothing := float32(0.15)
		player.Inventory[i].Position = rl.Vector2Lerp(
			player.Inventory[i].Position,
			targetPos,
			smoothing,
		)

		r.DrawProp(player.Inventory[i])
	}
}

func (player *Player) DrawHitbox() {
	rl.DrawRectangleLinesEx(player.HitboxRect, 1, rl.Blue)
	rl.DrawRectangleLinesEx(player.InteractiveRect, 1, rl.Yellow)
	rl.DrawPixel(int32(player.Position.X), int32(player.Position.Y), rl.Green)
}

func (player *Player) Tick(delta float32, level *Level) {
	player.ProcessInput(delta)
	player.CalculateVelocity(delta)
	player.UpdatePosition(delta, level)
	player.UpdateState()
	player.UpdateAnimation()
	player.PickupCollidingProps(level)
	player.OpenCollidingClosedDoors(level)
	player.RecordPath()

	if player.IsDroppingActiveProp {
		player.DropActiveProp(level)
	}
}

func (player *Player) ProcessInput(delta float32) {
	moveLeft := rl.IsKeyDown(rl.KeyA) || rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceLeft)
	moveRight := rl.IsKeyDown(rl.KeyD) || rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceRight)
	jump := rl.IsKeyDown(rl.KeySpace) || rl.IsGamepadButtonDown(0, rl.GamepadButtonRightFaceDown)

	jumpReleased := rl.IsKeyReleased(rl.KeySpace) || rl.IsGamepadButtonReleased(0, rl.GamepadButtonRightFaceDown)
	dropProp := rl.IsKeyReleased(rl.KeyF)

	// prevents spamming jumps
	if jumpReleased {
		player.CanJump = true
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

	if dropProp {
		if len(player.Inventory) < 1 {
			return
		}

		player.IsDroppingActiveProp = true
	}
}

func (player *Player) CalculateVelocity(delta float32) {
	if !player.OnGround {
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
	player.HandleCollisions(level.CollisionableHitboxes)

	if player.Position.Y < -5 {
		player.Position.Y = 180
		player.WentNorth = true
	}

	if player.Position.X > 320 {
		player.Position.X = 0
		player.WentEast = true
	}

	if player.Position.Y > 185 {
		player.Position.Y = 0
		player.WentSouth = true
	}

	if player.Position.X < 0 {
		player.Position.X = 320
		player.WentWest = true
	}
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

	player.InteractiveRect.X = player.Position.X - 2
	player.InteractiveRect.Y = player.Position.Y - 2
}

func (player *Player) UpdateAnimation() {
	player.FramesCounter++

	if !player.IsRunning {
		player.CurrentFrame = 0
		player.TextureRect.X = 0

		return
	}

	if player.FramesCounter >= (60 / player.FramesSpeed) && !player.IsJumping {
		player.FramesCounter = 0
		player.CurrentFrame++

		if !player.IsRunning && player.CurrentFrame > 1 {
			player.CurrentFrame = 0
		}

		if player.IsRunning && player.CurrentFrame > 1 {
			player.CurrentFrame = 0
		}

		if !player.IsRunning {
			player.TextureRect.X = 0
		}

		if player.IsRunning {
			offset := (float32(player.CurrentFrame) * 8) + 8
			player.TextureRect.X = offset
		}

	}
}

func (player *Player) HandleCollisions(collisionableElements []*rl.Rectangle) {
	player.OnGround = false

	for _, hitbox := range collisionableElements {
		if rl.CheckCollisionRecs(*hitbox, player.HitboxRect) {
			overlapLeft := (player.HitboxRect.X + player.HitboxRect.Width) - hitbox.X
			overlapRight := (hitbox.X + hitbox.Width) - player.HitboxRect.X
			overlapTop := (player.HitboxRect.Y + player.HitboxRect.Height) - hitbox.Y
			overlapBottom := (hitbox.Y + hitbox.Height) - player.HitboxRect.Y

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
				player.Position.Y = hitbox.Y - player.HitboxRect.Height
				player.Velocity.Y = 0
				player.OnGround = true
			case "BOTTOM":
				player.Position.Y = hitbox.Y + hitbox.Height
				player.Velocity.Y = 0
			case "LEFT":
				player.Position.X = hitbox.X - player.HitboxRect.Width
				player.Velocity.X = 0
			case "RIGHT":
				player.Position.X = hitbox.X + hitbox.Width
				player.Velocity.X = 0
			}
		}
	}
}

func (player *Player) PickupCollidingProps(level *Level) {
	for i := len(level.Props) - 1; i >= 0; i-- {
		prop := level.Props[i]
		if !prop.Pickable {
			continue
		}

		if rl.CheckCollisionRecs(player.InteractiveRect, prop.HitboxRect) {
			player.Inventory = append(player.Inventory, prop)
			level.Props = append(level.Props[:i], level.Props[i+1:]...)

			if player.ActivePropIndex == -1 {
				player.ActivePropIndex = 0
			}
		}
	}
}

func (player *Player) DropActiveProp(level *Level) {
	if len(player.Inventory) == 0 || player.ActivePropIndex < 0 || player.ActivePropIndex >= len(player.Inventory) {
		return
	}

	drop := player.Inventory[player.ActivePropIndex]
	drop.HitboxRect.X = drop.Position.X
	drop.HitboxRect.Y = drop.Position.Y

	level.Props = append(level.Props, drop)
	player.Inventory = append(player.Inventory[:player.ActivePropIndex], player.Inventory[player.ActivePropIndex+1:]...)

	if player.ActivePropIndex >= len(player.Inventory) && len(player.Inventory) > 0 {
		player.ActivePropIndex = len(player.Inventory) - 1
	} else if len(player.Inventory) == 0 {
		player.ActivePropIndex = -1
	}

	player.IsDroppingActiveProp = false
}

func (player *Player) IsMoving() bool {
	return player.Velocity.X != 0 || player.Velocity.Y != 0
}

func (player *Player) RecordPath() {
	if len(player.Path) < 1 {
		player.Path = append(player.Path, player.Position)
		return
	}

	if player.IsMoving() {
		if len(player.Path) >= 100 {
			player.Path = player.Path[1:]
		}

		player.Path = append(player.Path, player.Position)
	}
}

func (player *Player) OpenCollidingClosedDoors(l *Level) {
	for i, prop := range l.Props {
		if prop.Type != PropDoor {
			continue
		}

		if rl.CheckCollisionRecs(player.InteractiveRect, prop.HitboxRect) && player.HasKeyInInventory() {
			l.Props[i].IsOpen = true
			l.Props[i].Walkable = true
			player.RemoveKeyFromInventory()
		}
	}

	l.LoadCollisionables()
}

func (player *Player) HasKeyInInventory() bool {
	for _, item := range player.Inventory {
		if item.Type == PropKey {
			return true
		}
	}

	return false
}

func (player *Player) RemoveKeyFromInventory() {
	for i, item := range player.Inventory {
		if item.Type == PropKey {
			player.Inventory = append(player.Inventory[:i], player.Inventory[i+1:]...)
		}
	}
}
