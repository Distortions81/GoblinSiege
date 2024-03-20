package main

import "github.com/hajimehoshi/ebiten/v2/audio"

const (
	SND_ARROW_MISS = iota
	SND_ARROW_SHOOT
	SND_GOBLIN_DIE
	SND_GRASS_WALK
	SND_WIND
	SND_MAX
)

type soundData struct {
	file   string
	player *audio.Player
}

type varSoundData struct {
	numSounds int
	variants  []soundData
}

var varSounds []varSoundData

var sounds = [SND_MAX]soundData{
	{
		file: "arrow-miss.wav",
	},
	{
		file: "arrow-shoot.wav",
	},
	{
		file: "goblin-die.wav",
	},
	{
		file: "grass-walk.wav",
	},
	{
		file: "wind.wav",
	},
}
