package main

import (
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	monoFontBig   font.Face
	monoFont      font.Face
	monoFontSmall font.Face
)

func init() {

	/* Mono font */
	fontData := getFont("Hack-Regular.ttf")
	collection, err := opentype.ParseCollection(fontData)
	if err != nil {
		log.Fatal(err)
	}

	mono, err := collection.Font(0)
	if err != nil {
		log.Fatal(err)
	}

	/* Mono font big */
	monoFont, err = opentype.NewFace(mono, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	/* Mono font */
	monoFont, err = opentype.NewFace(mono, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	/* Mono font small */
	monoFontSmall, err = opentype.NewFace(mono, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}
