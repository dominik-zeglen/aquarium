package sim

import "github.com/golang/geo/r2"

type Action string

var attack = Action("attack")
var idle = Action("idle")

func isOutOfBounds(p r2.Point, e Environment) bool {
	return p.X < 0 || p.Y < 0 || p.X > float64(e.width) || p.Y > float64(e.height)
}
