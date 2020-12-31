package api

import (
	"github.com/dominik-zeglen/aquarium/sim"
	"github.com/golang/geo/r2"
)

type Query struct {
	s         *sim.Sim
	iteration *sim.IterationData
}

type CellArgs struct {
	ID int32
}

func (q *Query) Cell(args CellArgs) *CellResolver {
	cells := q.s.GetCells()
	id := int(args.ID)

	for _, cell := range cells {
		if cell.GetID() == id {
			resolver := CreateCellResolver(&cell, q.s)
			return &resolver
		}
	}

	return nil

}

type CellListArgs struct {
	Filter *CellFilter
}

func (q *Query) CellList(args CellListArgs) CellConnectionResolver {
	var cells []sim.Cell

	if args.Filter != nil && args.Filter.Area != nil {
		cells = q.s.GetCells().GetArea(args.Filter.Area.Start, args.Filter.Area.End)
	} else {

		cells = q.s.GetCells()
	}

	return CreateCellConnectionResolver(cells, q.s)
}

type SpeciesArgs struct {
	ID int32
}

func (q *Query) Species(args SpeciesArgs) *SpeciesResolver {
	species := q.s.GetAliveSpecies()
	id := int(args.ID)

	for _, species := range species {
		if species.ID == id {
			resolver := CreateSpeciesResolver(&species, q.s)
			return &resolver
		}
	}

	return nil
}

func (q *Query) SpeciesList() SpeciesConnectionResolver {
	species := q.s.GetAliveSpecies()

	return CreateSpeciesConnectionResolver(species, q.s)
}

type SpeciesGridArgs struct {
	Area AreaInput
}

func (q *Query) SpeciesGrid(args SpeciesGridArgs) []SpeciesGridElementResolver {
	scale := int32(1)
	if args.Area.Scale != nil {
		scale = *args.Area.Scale
	}
	cells := q.s.GetCells().GetArea(args.Area.Start, args.Area.End)
	grid := q.s.GetAliveSpecies().GetArea(cells, int(scale))

	resolvers := []SpeciesGridElementResolver{}

	for y := range grid {
		for x := range grid[y] {
			resolvers = append(resolvers, CreateSpeciesGridElementResolver(
				r2.Point{X: float64(x), Y: float64(y)},
				grid[y][x],
				q.s,
			))
		}
	}

	return resolvers
}

func (q *Query) Iteration() IterationResolver {
	return CreateIterationResolver(q.iteration, q.s)
}
