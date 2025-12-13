package collisions

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ray2D struct {
	Origin    rl.Vector2
	Direction rl.Vector2
}

type HitFace int

const (
	FaceTop HitFace = iota
	FaceRight
	FaceBottom
	FaceLeft
	NoFace
)

func CheckRectanglesCollision(movingRec rl.Rectangle, staticRec rl.Rectangle) (bool, HitFace) {
	hits := false
	hitFace := NoFace

	if movingRec.X < (staticRec.X+staticRec.Width) && (movingRec.X+movingRec.Width) > staticRec.X &&
		movingRec.Y < (staticRec.Y+staticRec.Height) && (movingRec.Y+movingRec.Height) > staticRec.Y {
		hits = true
	}

	if !hits {
		return false, hitFace
	}

	topOverlap := math.Abs(float64(movingRec.Y + movingRec.Height - staticRec.Y))
	bottomOverlap := math.Abs(float64(staticRec.Y + staticRec.Height - movingRec.Y))
	leftOverlap := math.Abs(float64(movingRec.X + movingRec.Width - staticRec.X))
	rightOverlap := math.Abs(float64(staticRec.X + staticRec.Width - movingRec.X))

	minOverlap := topOverlap
	hitFace = FaceTop

	if bottomOverlap < minOverlap {
		minOverlap = bottomOverlap
		hitFace = FaceBottom
	}

	if leftOverlap < minOverlap {
		minOverlap = leftOverlap
		hitFace = FaceLeft
	}

	if rightOverlap < minOverlap {
		hitFace = FaceRight
	}

	// This is my poor man's solution to ghost colliding
	// https://briansemrau.github.io/dealing-with-ghost-collisions/
	if (hitFace == FaceLeft || hitFace == FaceRight) && bottomOverlap > 0 && bottomOverlap > topOverlap {
		return false, NoFace
	}

	return hits, hitFace
}

func CheckRay2DRectangleCollision(ray Ray2D, targetRect rl.Rectangle, movingRecDimensions rl.Vector2) (float32, bool, HitFace) {
	return 0, true, NoFace
}
