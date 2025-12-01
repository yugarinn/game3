package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type Ray2D struct {
	Origin    rl.Vector2
	Direction rl.Vector2
}

func CheckRay2DRectangleCollision(ray Ray2D, targetRect rl.Rectangle, playerSize rl.Vector2) (float64, bool) {
	length := float32(math.Sqrt(float64(ray.Direction.X*ray.Direction.X + ray.Direction.Y*ray.Direction.Y)))
	if length == 0 {
		return 0, false
	}

	normalizedDir := rl.Vector2{X: ray.Direction.X / length, Y: ray.Direction.Y / length}
	expandedRect := rl.Rectangle{
		X:      targetRect.X - playerSize.X/2,
		Y:      targetRect.Y - playerSize.Y/2,
		Width:  targetRect.Width + playerSize.X,
		Height: targetRect.Height + playerSize.Y,
	}
}
