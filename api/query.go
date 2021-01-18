package api

import (
	"github.com/dominik-zeglen/aquarium/sim"
	"github.com/golang/geo/r2"
)

type Query struct {
	s         *sim.Sim
	iteration *sim.IterationData
}

type OrganismArgs struct {
	ID int32
}

func (q *Query) Organism(args OrganismArgs) *OrganismResolver {
	organisms := q.s.GetOrganisms()
	id := int(args.ID)

	for _, cell := range organisms {
		if cell.GetID() == id {
			resolver := OrganismResolver{cell}
			return &resolver
		}
	}

	return nil

}

type OrganismListArgs struct {
	Filter *OrganismFilter
}

func (q *Query) OrganismList(args OrganismListArgs) []OrganismResolver {
	var organisms sim.OrganismList

	if args.Filter != nil && args.Filter.Area != nil {
		organisms = q.s.GetOrganisms().GetAlive().GetArea(args.Filter.Area.Start, args.Filter.Area.End)
	} else {
		organisms = q.s.GetOrganisms()
	}

	return createOrganismResolverList(organisms)
}

type SpeciesArgs struct {
	ID int32
}

func (q *Query) Species(args SpeciesArgs) *SpeciesResolver {
	species := q.s.GetSpecies().GetAlive()
	id := int(args.ID)

	for _, species := range species {
		if species.GetID() == id {
			resolver := SpeciesResolver{species}
			return &resolver
		}
	}

	return nil
}

func (q *Query) SpeciesList() []SpeciesResolver {
	species := q.s.GetSpecies().GetAlive()

	return createSpeciesResolverList(species)
}

type SpeciesGridArgs struct {
	Area AreaInput
}

func (q *Query) SpeciesGrid(args SpeciesGridArgs) []SpeciesGridElementResolver {
	scale := int32(1)
	if args.Area.Scale != nil {
		scale = *args.Area.Scale
	}
	organisms := q.s.GetOrganisms().GetAlive().GetArea(args.Area.Start, args.Area.End)
	grid := q.s.GetSpecies().GetAlive().GetArea(organisms, int(scale))

	resolvers := []SpeciesGridElementResolver{}

	for y := range grid {
		for x := range grid[y] {
			resolvers = append(resolvers, CreateSpeciesGridElementResolver(
				r2.Point{X: float64(x), Y: float64(y)},
				grid[y][x],
			))
		}
	}

	return resolvers
}

func (q *Query) Iteration() IterationResolver {
	return CreateIterationResolver(q.iteration, q.s)
}
