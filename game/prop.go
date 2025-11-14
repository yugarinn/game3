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
	HitboxRect             rl.Rectangle
	IsSelected             bool
}

type LDtkEntity struct {
	Width  int 
	Height int 
    Px     []float32
}

func NewPropFromLDtk(entity *LDtkEntity) *Prop {
	return &Prop{
		Position: entity.Px,
		HitboxRect: rl.NewRectangle(entity.Px[0], entity.Px[1], 8, 8),
		Pickable: true,
	}
}
