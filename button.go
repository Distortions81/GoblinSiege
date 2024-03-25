package main

import (
	"fmt"
	"log"
	"os"
)

type buttonData struct {
	name        string
	topLeft     xyi
	bottomRight xyi

	size   xyf32
	pos    xyf32
	action func()
}

var splashButtons = []buttonData{
	{
		name:        "Play Game",
		topLeft:     xyi{X: 100, Y: 215},
		bottomRight: xyi{X: 319, Y: 289},
		action:      actPlayGame,
	},
	{
		name:        "Settings",
		topLeft:     xyi{X: 100, Y: 324},
		bottomRight: xyi{X: 319, Y: 397},
		action:      actSettings,
	},
	{
		name:        "Quit",
		topLeft:     xyi{X: 100, Y: 431},
		bottomRight: xyi{X: 319, Y: 505},
		action:      actQuit,
	},
}

func init() {
	for b, button := range splashButtons {
		splashButtons[b].size = xyf32{
			X: float32(button.bottomRight.X - button.topLeft.X),
			Y: float32(button.topLeft.Y - button.bottomRight.Y),
		}
		splashButtons[b].pos = xyf32{
			X: float32(button.topLeft.X) + float32(button.size.X/2),
			Y: float32(button.bottomRight.Y) - float32(button.size.Y/2),
		}
		buf := fmt.Sprintf("%v: size: %v,%v pos: %v,%v",
			splashButtons[b].name,
			splashButtons[b].size.X, splashButtons[b].size.Y,
			splashButtons[b].pos.X, splashButtons[b].pos.Y)
		log.Println(buf)
	}
}

func actPlayGame() {
	gameMode = MODE_PLAY_TWITCH
}

func actSettings() {
	gameMode = MODE_SPLASH
}

func actQuit() {
	os.Exit(0)
}
