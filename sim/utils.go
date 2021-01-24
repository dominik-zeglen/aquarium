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

func getFreeSpot(
	cells CellList,
	cell Cell,
	canConnect bool,
) *r2.Point {
	dist := float64(1)
	if !canConnect {
		dist = 5
	}

	candidates := []r2.Point{
		{X: 0, Y: dist},
		{X: dist, Y: 0},
		{X: 0, Y: -dist},
		{X: -dist, Y: 0},
	}

	for _, candidate := range candidates {
		newPos := cell.position.Add(candidate)
		available := true

		for _, cellToCheck := range cells {
			if cellToCheck.position.Sub(newPos).Norm() == 0 {
				available = false
				break
			}
		}

		if available {
			return &newPos
		}

	}

	return nil
}
