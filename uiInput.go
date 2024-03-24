package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {

	if !ebiten.IsFocused() {
		return nil
	}

	mx, my := ebiten.CursorPosition()
	lmb := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	//mmb := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle)
	//rmb := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

	if gameMode == MODE_SPLASH {
		if lmb {
			for _, button := range splashButtons {
				if mx >= button.topLeft.X &&
					mx <= button.bottomRight.X &&
					my >= button.topLeft.Y &&
					my <= button.bottomRight.Y {

					button.action()
					log.Println("meep meep click go")
					break
				}
			}
		}
	}

	return nil
}
