package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {

	if !gameLoaded.Load() {
		return nil
	}

	if !ebiten.IsFocused() {
		return nil
	}

	mx, my := ebiten.CursorPosition()
	lmb := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	//mmb := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle)
	//rmb := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

	if gameMode.Load() == MODE_SPLASH {

		for b, button := range splashButtons {

			if mx >= button.topLeft.X &&
				mx <= button.bottomRight.X &&
				my >= button.topLeft.Y &&
				my <= button.bottomRight.Y {

				if lmb {
					log.Printf("Click: %v", button.name)
					splashButtons[b].clicked = time.Now()
					go func() {
						time.Sleep(flashSpeed)
						button.action()
					}()
					break
				} else {
					splashButtons[b].hover = true
					break
				}
			} else {
				splashButtons[b].hover = false
			}

		}
	}

	return nil
}
