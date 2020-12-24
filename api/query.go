package api

import (
	"github.com/dominik-zeglen/aquarium/sim"
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
