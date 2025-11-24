package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PropType int

const (
	PropKey PropType = iota
	PropDoor
	PropGrass
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
	Width                  float32
	Height                 float32
	IsAnimated             bool
	CurrentFrame           int32
	AnimationFramesCounter int32
	FramesCounter          int32
	FramesSpeed            int32
	HitboxRect             rl.Rectangle
	IsOpen                 bool
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

	propType := entity.GetPropType()
	width, height := getDimensionsForType(propType)
	position := rl.NewVector2(entity.Px[0], entity.Px[1])
	hitbox := rl.NewRectangle(position.X, position.Y, width, height)

	return &Prop{
		Type:       propType,
		Walkable:   walkable,
		Pickable:   pickable,
		Position:   position,
		HitboxRect: hitbox,
		Width:      width,
		Height:     height,
	}
}

func (entity *LDtkEntity) GetPropType() PropType {
	types := map[string]PropType{
		"Key":    PropKey,
		"Door":   PropDoor,
		"Grass1": PropGrass,
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

func getDimensionsForType(propType PropType) (float32, float32) {
	dimensions := map[PropType][]float32{
		PropKey:     {8, 8},
		PropDoor:    {16, 16},
		PropGrass:   {8, 8},
		PropGeneral: {8, 8},
	}

	if typeDimensions, ok := dimensions[propType]; ok {
		return typeDimensions[0], typeDimensions[1]
	}

	return dimensions[PropGeneral][0], dimensions[PropGeneral][1]
}
