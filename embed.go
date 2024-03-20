package main

import (
	"embed"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	//go:embed data
	f     embed.FS
	bgimg *ebiten.Image
)

func init() {
	var err error
	bgimg, _, err = ebitenutil.NewImageFromFile("data/maps/main.png")
	if err != nil {
		log.Fatal(err)
	}

	for s, sheet := range sheetPile {
		tmp, _, err := ebitenutil.NewImageFromFile("data/sprites/" + sheet.file + ".png")
		sheetPile[s].img = tmp

		log.Printf("Loaded spritesheet: %v", sheet.name)

		if err != nil {
			log.Fatal(err)
		}

		if sheet.frames <= 0 {
			sheetPile[s].frameSize.X = tmp.Bounds().Dx()
			sheetPile[s].frameSize.Y = tmp.Bounds().Dy()
			continue
		}
		for a, ani := range sheet.anims {
			for x := 0; x < sheet.frames; x++ {
				log.Printf("Sliced: %v frame %v\n", ani.name, x)
				sheet.anims[a].img = append(sheet.anims[a].img, getAni(sheet, a, x))
			}
		}
	}

	for s, sound := range sounds {
		sRead, err := os.Open("data/sounds/" + sound.file)
		if err != nil {
			log.Fatal(err)
		}
		sounds[s].wav, err = wav.DecodeWithoutResampling(sRead)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Loaded %v.", sound.file)
	}
}

func getFont(name string) []byte {
	data, err := f.ReadFile("data/fonts/" + name)
	if err != nil {
		log.Fatal(err)
	}
	return data

}
