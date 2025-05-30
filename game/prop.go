package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Prop struct {
	Type                   string
	Pickable               bool
	Walkable               bool
	Pushable               bool
	Sprite                 rl.Texture2D
	TextureRect            rl.Rectangle
	Position               []float32
	IsAnimated             bool
	CurrentFrame           int32
	AnimationFramesCounter int32
	FramesCounter          int32
	FramesSpeed            int32
}
