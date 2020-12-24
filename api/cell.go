package api

import (
	"github.com/dominik-zeglen/aquarium/sim"
	"github.com/golang/geo/r2"
)

type CellResolver struct {
	cell *sim.Cell
	s    *sim.Sim
}

func CreateCellResolver(cell *sim.Cell, sim *sim.Sim) CellResolver {
	return CellResolver{cell, sim}
}

func (res CellResolver) ID() int32 {
	return int32(res.cell.GetID())
}

func (res CellResolver) Alive() bool {
	return res.cell.IsAlive()
}

func (res CellResolver) BornAt() int32 {
	return int32(res.cell.GetBornAt())
}

func (res CellResolver) Capacity() int32 {
	return int32(res.cell.GetCapacity())
}

func (res CellResolver) Food() int32 {
	return int32(res.cell.GetFood(res.s.GetEnvironment(), res.s.GetIteration()))
}

func (res CellResolver) Hp() int32 {
	return int32(res.cell.GetHP())
}

func (res CellResolver) Position() r2.Point {
	return res.cell.GetPosition()
}

func (res CellResolver) Satiation() int32 {
	return int32(res.cell.GetSatiation())
}

func (res CellResolver) Species() SpeciesResolver {
	return CreateSpeciesResolver(res.cell.GetSpecies(), res.s)
}

type CellConnectionEdgeResolver struct {
	cell sim.Cell
	s    *sim.Sim
}

func CreateCellConnectionEdgeResolver(cell sim.Cell, sim *sim.Sim) CellConnectionEdgeResolver {
	return CellConnectionEdgeResolver{cell, sim}
}

func (res CellConnectionEdgeResolver) Node() CellResolver {
	return CreateCellResolver(&res.cell, res.s)
}

type CellConnectionResolver struct {
	cells []sim.Cell
	s     *sim.Sim
}

func CreateCellConnectionResolver(cells []sim.Cell, sim *sim.Sim) CellConnectionResolver {
	return CellConnectionResolver{cells, sim}
}

func (res CellConnectionResolver) Count() int32 {
	return int32(len(res.cells))
}

func (res CellConnectionResolver) Edges() []CellConnectionEdgeResolver {
	resolvers := make([]CellConnectionEdgeResolver, len(res.cells))

	for cellIndex := range res.cells {
		resolvers[cellIndex] = CreateCellConnectionEdgeResolver(res.cells[cellIndex], res.s)
	}

	return resolvers
}
