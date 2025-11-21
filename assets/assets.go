package assets

import _ "embed"

/** ENTITIES SPRITES */
//go:embed entities/player.png
var PLAYER_SPRITE_DATA []byte

/** GROUND SPRITES */
//go:embed ground.png
var GROUND_SPRITE_DATA []byte

/** BACKGROUNDS */
//go:embed background-daylight-sky.png
var BACKRGOUND_DAYLIGHT_SKY []byte

//go:embed background-underground.png
var BACKRGOUND_UNDERGROUND []byte

//go:embed tilemap.png
var TILEMAP []byte
