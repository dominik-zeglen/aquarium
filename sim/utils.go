package sim

import "github.com/golang/geo/r2"

func fitToBoundary(p r2.Point, env Environment) r2.Point {
	x := p.X
	if x > float64(env.width) {
		x = float64(env.width) - 1
	}
	if x < 0 {
		x = 0
	}

	y := p.Y
	if y > float64(env.height) {
		y = float64(env.height - 1)
	}
	if y < 0 {
		y = 0
	}

	return r2.Point{x, y}
}
