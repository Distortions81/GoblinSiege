package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type objectData struct {
	Pos  xyi
	Type int
}

var (
	gameMap map[xyi]*objectData
)

func drawGameBoard(screen *ebiten.Image) {
	for _, item := range gameMap {
		vector.DrawFilledRect(screen, float32(item.Pos.X), float32(item.Pos.Y), 1, 1, color.White, false)
	}
}
