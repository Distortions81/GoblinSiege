package main

import "math"

/* give distance between two coordinates */
func Distance(a, b xyi) float64 {
	x := math.Abs(float64(a.X - b.X))
	y := math.Abs(float64(a.Y - b.Y))
	return float64(math.Sqrt(x*x + y*y))
}
