package sim

import (
	"math"
	"math/rand"

	"github.com/golang/geo/r2"
)

type Action string

var attack = Action("attack")
var idle = Action("idle")

func isOutOfBounds(p r2.Point, e Environment) bool {
	return p.X < 0 || p.Y < 0 || p.X > float64(e.width) || p.Y > float64(e.height)
}

func getVecFromAngle(angle float64) r2.Point {
	return r2.Point{
		X: math.Cos(angle),
		Y: math.Sin(angle),
	}
}

func getRandomAngle() float64 {
	return rand.Float64() * 2 * math.Pi
}

func getRandomVec() r2.Point {
	return getVecFromAngle(getRandomAngle())
}
