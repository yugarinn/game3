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

	var earlierXCrossing, laterXCrossing float32
	if normalizedDir.X != 0 {
		leftCrossing := (expandedRect.X - ray.Origin.X) / normalizedDir.X
		rightCrossing := (expandedRect.X + expandedRect.Width - ray.Origin.X) / normalizedDir.X

		earlierXCrossing = float32(math.Min(float64(leftCrossing), float64(rightCrossing)))
		laterXCrossing = float32(math.Max(float64(leftCrossing), float64(rightCrossing)))
	} else {
		if ray.Origin.X < expandedRect.X || ray.Origin.X > expandedRect.X+expandedRect.Width {
			return 0, false
		}

		earlierXCrossing = float32(math.Inf(-1))
		laterXCrossing = float32(math.Inf(1))
	}

	var earlierYCrossing, laterYCrossing float32
	if normalizedDir.Y != 0 {
		topCrossing := (expandedRect.Y - ray.Origin.Y) / normalizedDir.Y
		bottomCrossing := (expandedRect.Y + expandedRect.Height - ray.Origin.Y) / normalizedDir.Y

		earlierYCrossing = float32(math.Min(float64(topCrossing), float64(bottomCrossing)))
		laterYCrossing = float32(math.Max(float64(topCrossing), float64(bottomCrossing)))
	} else {
		if ray.Origin.Y < expandedRect.Y || ray.Origin.Y > expandedRect.Y+expandedRect.Height {
			return 0, false
		}

		earlierYCrossing = float32(math.Inf(-1))
		laterYCrossing = float32(math.Inf(1))
	}

	timeToCrossStart := float32(math.Max(float64(earlierXCrossing), float64(earlierYCrossing)))
	timeToCrossEnd := float32(math.Min(float64(laterXCrossing), float64(laterYCrossing)))

	if timeToCrossEnd < 0 || timeToCrossStart > timeToCrossEnd {
		return 0, false
	}

	if timeToCrossStart < 0 {
		return 0, true
	}

	return float64(timeToCrossStart), true
}
