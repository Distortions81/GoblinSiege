package main

type buttonData struct {
	name        string
	topLeft     xyi
	bottomRight xyi
	action      func()
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

func actPlayGame() {

}

func actSettings() {

}

func actQuit() {

}
