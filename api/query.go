package api

import (
	"github.com/dominik-zeglen/aquarium/sim"
	"github.com/golang/geo/r2"
)

type Query struct {
	s *sim.Sim
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

func (q *Query) CellList() CellConnectionResolver {
	cells := q.s.GetCells()

	return CreateCellConnectionResolver(cells, q.s)
}

type AreaArgs struct {
	Start r2.Point
	End   r2.Point
}

func (q *Query) Area(args AreaArgs) []CellResolver {
	cells := q.s.GetCells()
	resolvers := []CellResolver{}

	for cellIndex, cell := range cells {
		position := cell.GetPosition()
		if position.X > args.Start.X && position.X < args.End.X &&
			position.Y > args.Start.Y && position.Y < args.End.Y {
			resolvers = append(resolvers, CreateCellResolver(&cells[cellIndex], q.s))
		}
	}

	return resolvers
}
