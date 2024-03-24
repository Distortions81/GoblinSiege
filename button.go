package main

import (
	"fmt"
	"log"
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
		topLeft:     xyi{X: 100, Y: 214},
		bottomRight: xyi{X: 320, Y: 289},
		action:      actPlayGame,
	},
	{
		name:        "Settings",
		topLeft:     xyi{X: 100, Y: 323},
		bottomRight: xyi{X: 320, Y: 397},
		action:      actPlayGame,
	},
	{
		name:        "Quit",
		topLeft:     xyi{X: 100, Y: 431},
		bottomRight: xyi{X: 320, Y: 504},
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
		buf := fmt.Sprintf("%v: size: %v,%v pos: %v,%v", button.name, button.size.X, button.size.Y, button.pos.X, button.pos.Y)
		log.Println(buf)
	}
}

func actPlayGame() {

}

func actSettings() {

}

func actQuit() {

}
