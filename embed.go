package main

import (
	"embed"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	//go:embed data
	f        embed.FS
	bgimg    *ebiten.Image
	splash   *ebiten.Image
	audioCon *audio.Context
)

func loadSheet(sheet *spriteSheetData) {
	tmp, _, err := ebitenutil.NewImageFromFile("data/sprites/" + sheet.file + ".png")
	if err != nil {
		log.Fatal(err)
	}
	sheet.img = tmp

	log.Printf("Loaded spritesheet: %v", sheet.name)

	if sheet.frames <= 0 {
		sheet.frameSize.X = tmp.Bounds().Dx()
		sheet.frameSize.Y = tmp.Bounds().Dy()
		return
	}
	for a, ani := range sheet.anims {
		if ani.name == "" {
			continue
		}
		for x := 0; x < sheet.frames; x++ {
			log.Printf("Sliced: %v frame %v\n", ani.name, x)
			sheet.anims[a].img = append(sheet.anims[a].img, getAni(sheet, a, x))
		}
	}
}

func getUpSheet(item *objectData) *spriteSheetData {
	//If item is not upgraded, or upgrade not found return
	//the item's standard sprite sheet
	if item.upgrade == 0 ||
		len(item.sheetP.upgrades) < (item.upgrade-1) {
		return item.sheetP
	} else {
		//Otherwise return upgraded copy
		return item.sheetP.upgrades[item.upgrade-1]
	}
}

func loadAssets() {
	var err error
	bgimg, _, err = ebitenutil.NewImageFromFile("data/maps/main.png")
	if err != nil {
		log.Fatal(err)
	}

	splash, _, err = ebitenutil.NewImageFromFile("data/sprites/splash.png")
	if err != nil {
		log.Fatal(err)
	}

	for _, sheet := range sheetPile {
		//Handle upgrades
		for _, upgrade := range sheet.upgrades {
			upgrade.name = upgrade.file
			upgrade.health = sheet.health
			upgrade.frameSize = sheet.frameSize
			upgrade.frames = sheet.frames
			upgrade.anims = sheet.anims
		}

		for _, gsheet := range sheet.upgrades {
			loadSheet(gsheet)
		}

		loadSheet(sheet)
	}

	audioCon = audio.NewContext(44100)

	//Standard sounds
	for s, snd := range sounds {
		if snd.variated {
			continue
		}
		audioPlayer := getSound("data/sounds/" + snd.file)
		sounds[s].player = audioPlayer
		log.Printf("Loaded %v.", snd.file)
	}
	//Variated
	for s, vsnd := range sounds {
		if !vsnd.variated {
			continue
		}
		variations, err := os.ReadDir("data/sounds/" + vsnd.file + "/")
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range variations {
			if item.IsDir() {
				continue
			}
			if strings.HasSuffix(item.Name(), ".wav") {
				audioPlayer := getSound("data/sounds/" + vsnd.file + "/" + item.Name())
				newSound := soundData{player: audioPlayer}
				varSounds[s].variants = append(varSounds[s].variants, newSound)
				varSounds[s].numSounds++
				log.Printf("Loaded %v.", item.Name())
			}
		}
	}
}

func getFont(name string) []byte {
	data, err := f.ReadFile("data/fonts/" + name)
	if err != nil {
		log.Fatal(err)
	}
	return data

}

func getSound(input string) *audio.Player {
	sndReader, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}

	wav, err := wav.DecodeWithoutResampling(sndReader)
	if err != nil {
		log.Fatal(err)
	}

	audioPlayer, err := audioCon.NewPlayer(wav)
	if err != nil {
		log.Fatal(err)
	}
	return audioPlayer
}
