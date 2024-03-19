package main

import (
	"image/color"
	"math"
	"strings"
)

/* give distance between two coordinates */
func Distance(a, b xyi) float64 {
	x := math.Abs(float64(a.X - b.X))
	y := math.Abs(float64(a.Y - b.Y))
	return float64(math.Sqrt(x*x + y*y))
}

func healthColor(input float32) color.NRGBA {
	var healthColor color.NRGBA = color.NRGBA{R: 255, G: 255, B: 255, A: 0}
	health := input * 100

	if health < 100 && health > 0 {
		healthColor.A = 128
		healthColor.B = 0

		r := int(float32(100-(health)) * 5)
		if r > 255 {
			r = 255
		}
		healthColor.R = uint8(r)

		g := int(float32(health) * 4)
		if g > 255 {
			g = 255
		}
		healthColor.G = uint8(g)

	}

	return healthColor
}

func makeEllipsis() string {
	return strings.Repeat(".", (int(frameCount%120) / 30))
}
