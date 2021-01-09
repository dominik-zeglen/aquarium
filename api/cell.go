package api

import (
	"github.com/dominik-zeglen/aquarium/sim"
	"github.com/golang/geo/r2"
)

type CellResolver struct {
	cell sim.Cell
}

func createCellResolverList(cells sim.CellList) []CellResolver {
	resolvers := make([]CellResolver, len(cells))

	for speciesIndex := range cells {
		resolvers[speciesIndex] = CellResolver{cells[speciesIndex]}
	}

	return resolvers
}

func (res CellResolver) ID() int32 {
	return int32(res.cell.GetID())
}

func (res CellResolver) Alive() bool {
	return res.cell.IsAlive()
}
func (res CellResolver) Hp() int32 {
	return int32(res.cell.GetHP())
}

func (res CellResolver) Position() r2.Point {
	return res.cell.GetPosition()
}

func (res CellResolver) Type() CellTypeResolver {
	return CellTypeResolver{*res.cell.GetType()}
}
