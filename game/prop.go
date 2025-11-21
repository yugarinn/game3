package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PropType int

const (
	PropKey PropType = iota
	PropDoor
	PropGeneral
)

type Prop struct {
	Type                   PropType
	Pickable               bool
	Walkable               bool
	Pushable               bool
	Sprite                 rl.Texture2D
	TextureRect            rl.Rectangle
	Position               rl.Vector2
	IsAnimated             bool
	CurrentFrame           int32
	AnimationFramesCounter int32
	FramesCounter          int32
	FramesSpeed            int32
	HitboxRect             rl.Rectangle
	IsSelected             bool
}

type LDtkEntityCustomField struct {
	Identifier string `json:"__identifier"`
	Type       string `json:"__type"`
	Value      any    `json:"__value"`
}

type LDtkEntity struct {
	ID           string                  `json:"__identifier"`
	CustomFields []LDtkEntityCustomField `json:"fieldInstances"`
	Width        int
	Height       int
	Px           []float32
}

func NewPropFromLDtk(entity *LDtkEntity) *Prop {
	walkable := false
	if val := entity.GetCustomFieldValue("Walkable"); val != nil {
		if b, ok := val.(bool); ok {
			walkable = b
		}
	}

	pickable := false
	if val := entity.GetCustomFieldValue("Pickable"); val != nil {
		if b, ok := val.(bool); ok {
			pickable = b
		}
	}

	return &Prop{
		Type:       entity.GetPropType(),
		Walkable:   walkable,
		Pickable:   pickable,
		Position:   rl.NewVector2(entity.Px[0], entity.Px[1]),
		HitboxRect: rl.NewRectangle(entity.Px[0], entity.Px[1], 8, 8),
	}
}

func (entity *LDtkEntity) GetPropType() PropType {
	types := map[string]PropType{
		"Key":  PropKey,
		"Door": PropDoor,
	}

	if propType, ok := types[entity.ID]; ok {
		return propType
	}

	return PropGeneral
}

func (entity *LDtkEntity) GetCustomFieldValue(identifier string) any {
	for _, field := range entity.CustomFields {
		if field.Identifier == identifier {
			return field.Value
		}
	}

	return nil
}
