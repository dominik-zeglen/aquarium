package sim

import (
	"math/rand"
	"sort"

	"github.com/golang/geo/r2"
)

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

type ByLength []r2.Point

func (a ByLength) Len() int           { return len(a) }
func (a ByLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLength) Less(i, j int) bool { return a[i].Norm() < a[j].Norm() }

func getFreeSpot(
	cells CellList,
	cell Cell,
	canConnect bool,
) *r2.Point {
	dist := float64(1)
	if !canConnect || rand.Float64() > .8 {
		dist = 2
	}

	candidates := []r2.Point{
		{X: 0, Y: dist},
		{X: dist, Y: 0},
		{X: 0, Y: -dist},
		{X: -dist, Y: 0},
	}

	newPositions := make([]r2.Point, 4)
	for posIndex, pos := range candidates {
		newPositions[posIndex] = cell.position.Add(pos)
	}
	sort.Sort(ByLength(newPositions))

	for _, newPos := range newPositions {
		isInBounds := newPos.X >= 0 && newPos.X < gridSize/2 && newPos.Y >= 0 &&
			newPos.Y < gridSize/2

		if !isInBounds {
			continue
		}

		available := true

		for _, cellToCheck := range cells {
			if cellToCheck.position.Sub(newPos).Norm() < dist {
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
