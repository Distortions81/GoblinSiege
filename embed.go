package main

import (
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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

	for i, item := range oTypes {
		if item.spriteName == "" {
			continue
		}
		tmp, _, err := ebitenutil.NewImageFromFile("data/sprites/" + item.spriteName + ".png")
		oTypes[i].spriteImg = tmp

		log.Printf("Loaded sprite: %v", item.spriteName)

		if err != nil {
			log.Fatal(err)
		}
	}
	for i, item := range oTypes {
		if item.deadName == "" {
			continue
		}
		tmp, _, err := ebitenutil.NewImageFromFile("data/sprites/" + item.deadName + ".png")
		oTypes[i].deadImg = tmp

		log.Printf("Loaded sprite: %v", item.spriteName)

		if err != nil {
			log.Fatal(err)
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
