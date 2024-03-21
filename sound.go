package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	SND_ARROW_MISS = iota
	SND_ARROW_SHOOT
	SND_GOBLIN_DIE
	SND_GRASS_WALK
	SND_WIND
	SND_AXE
	SND_TOWER_DIE
	SND_MAX
)

func playSound(s int) {
	if sounds[s].player.IsPlaying() {
		return
	}
	sounds[s].player.SetVolume(0)
	sounds[s].player.Pause()
	sounds[s].player.Rewind()
	sounds[s].player.SetVolume(defaultVolume)
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

	sound.SetVolume(0)
	sound.Pause()
	sound.Rewind()
	sound.SetVolume(defaultVolume)
	sound.Play()
}

type soundData struct {
	file     string
	player   *audio.Player
	variated bool
}

type variSoundData struct {
	numSounds int
	variants  []soundData
}

var varSounds [SND_MAX]variSoundData

var sounds = [SND_MAX]soundData{
	{
		file: "arrow-miss.wav",
	},
	{
		variated: true,
		file:     "arrow-shoot",
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
	{
		variated: true,
		file:     "axe",
	},
	{
		file: "tower-die.wav",
	},
}
