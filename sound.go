package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	SND_ARROW_SHOOT = iota
	SND_GOBLIN_DIE
	SND_GRASS_WALK
	SND_WIND
	SND_AXE
	SND_TOWER_DIE
	SND_GAMEOVER
	SND_TENSION
	SND_GAMEWON
	SND_MAX
)

func playSound(s int) {
	if sounds[s].player.IsPlaying() {
		return
	}
	sounds[s].player.SetVolume(sounds[s].vol)
	sounds[s].player.Rewind()
	sounds[s].player.Play()
}

func playVariated(s int) {

	//Get a random variation of the sound
	vSounds := varSounds[s]
	sound := vSounds.variants[rand.Intn(vSounds.numSounds)].player

	//Sound channel busy, try to find another
	if sound.IsPlaying() {
		for x := 0; x < vSounds.numSounds; x++ {
			sound = vSounds.variants[x].player
			if sound.IsPlaying() {
				continue
			} else {
				break
			}
		}
	}

	sound.Rewind()
	sound.SetVolume(sounds[s].vol)
	sound.Play()
}

type soundData struct {
	file     string
	player   *audio.Player
	variated bool
	vol      float64
}

type variSoundData struct {
	numSounds int
	variants  []soundData
}

// Sounds that have variations for variety
var varSounds [SND_MAX]variSoundData

var sounds = [SND_MAX]soundData{
	{
		variated: true,
		file:     "arrow-shoot",
		vol:      0.5,
	},
	{
		variated: true,
		file:     "goblin-die",
		vol:      0.25,
	},
	{
		file: "grass-walk.wav",
		vol:  0.1,
	},
	{
		file: "wind.wav",
		vol:  0.10,
	},
	{
		variated: true,
		file:     "axe",
		vol:      0.4,
	},
	{
		file: "tower-die.wav",
		vol:  0.3,
	},
	{
		file: "gameover.wav",
		vol:  1.0,
	},
	{
		file: "tension.wav",
		vol:  0.7,
	},
	{
		file: "gamewon.wav",
		vol:  0.7,
	},
}
