package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type spriteSheetData struct {
	name      string
	health    int
	file      string
	frameSize xyi
	frames    int
	anims     [ANI_MAX]animationData
	img       *ebiten.Image
}

type animationData struct {
	name   string
	row    int
	mirror bool
	img    []*ebiten.Image
}

const (
	ANI_IDLE = iota
	ANI_RUN
	ANI_FADE
	ANI_DIE
	ANI_ATTACK
	ANI_SWIM
	ANI_MAX
)

// heh
var sheetPile = []*spriteSheetData{
	&obj_goblinBarb,
	&obj_tower1,
	&obj_arrow,
}

var obj_goblinBarb = spriteSheetData{
	name:      "goblin barbarian",
	file:      "Goblin_Barbarian",
	health:    100,
	frameSize: xyi{X: 32, Y: 32},
	frames:    4,

	anims: [ANI_MAX]animationData{
		{
			name:   "idle",
			row:    1,
			mirror: true,
		},
		{
			name:   "run",
			row:    4,
			mirror: true,
		},
		{
			name:   "fade",
			row:    7,
			mirror: true,
		},
		{
			name:   "die",
			row:    10,
			mirror: true,
		},
		{
			name:   "attack",
			row:    13,
			mirror: true,
		},
		{
			name:   "swim",
			row:    15,
			mirror: true,
		},
	},
}

var obj_tower1 = spriteSheetData{
	name:      "tower",
	file:      "tower",
	health:    100,
	frameSize: xyi{X: 32, Y: 64},
	frames:    3,

	anims: [ANI_MAX]animationData{
		{
			name: "idle",
			row:  0,
		},
		{
			name: "build",
			row:  1,
		},
		{
			name: "broken",
			row:  2,
		},
	},
}

var obj_arrow = spriteSheetData{
	name: "arrow",
	file: "arrow",
}

func getAni(sheet *spriteSheetData, auiNum int, frame int) *ebiten.Image {
	aniData := sheet.anims[auiNum]
	var r image.Rectangle
	if aniData.mirror {
		r = image.Rect(
			(sheet.frameSize.X*frame)+sheet.frameSize.X,
			aniData.row*sheet.frameSize.Y,
			sheet.frameSize.X*frame,
			(aniData.row*sheet.frameSize.Y)+sheet.frameSize.Y)
	} else {
		r = image.Rect(
			sheet.frameSize.X*frame,
			aniData.row*sheet.frameSize.Y,
			(sheet.frameSize.X*frame)+sheet.frameSize.X,
			(aniData.row*sheet.frameSize.Y)+sheet.frameSize.Y)
	}
	return sheet.img.SubImage(r).(*ebiten.Image)
}
