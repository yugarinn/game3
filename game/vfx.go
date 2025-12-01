package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type VFXType int

const (
	PlayerJumpVFX VFXType = iota
	PlayerDeathVFX
)

type VFX struct {
	Type                       VFXType
	Position                   rl.Vector2
	AnimationCurrentFrame      int32
	AnimationCurrentPosition   int32
	AnimationFramesPerPosition int32
	AnimationFramesCounter     int32
	AnimationPositionsCounter  int32
	Loops                      bool
}

func NewVFX(vfxType VFXType, position rl.Vector2) VFX {
	inventory := map[VFXType]VFX{
		PlayerJumpVFX: {
			Type:                       PlayerJumpVFX,
			Position:                   position,
			AnimationFramesPerPosition: 4,
			AnimationPositionsCounter:  4,
			Loops:                      false,
		},
		PlayerDeathVFX: {
			Type:                       PlayerDeathVFX,
			Position:                   position,
			AnimationFramesPerPosition: 10,
			AnimationPositionsCounter:  3,
			Loops:                      false,
		},
	}

	if vfx, ok := inventory[vfxType]; ok {
		return vfx
	}

	return VFX{}
}
