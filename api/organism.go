package api

import (
	"github.com/dominik-zeglen/aquarium/sim"
	"github.com/golang/geo/r2"
)

type OrganismResolver struct {
	organism sim.Organism
}

func createOrganismResolverList(organisms sim.OrganismList) []OrganismResolver {
	resolvers := make([]OrganismResolver, len(organisms))

	for speciesIndex := range organisms {
		resolvers[speciesIndex] = OrganismResolver{organisms[speciesIndex]}
	}

	return resolvers
}

func (res OrganismResolver) ID() int32 {
	return int32(res.organism.GetID())
}

func (res OrganismResolver) Species() SpeciesResolver {
	return SpeciesResolver{res.organism.GetSpecies()}
}

func (res OrganismResolver) BornAt() int32 {
	return int32(res.organism.GetBornAt())
}

func (res OrganismResolver) Cells() []CellResolver {
	return createCellResolverList(res.organism.GetCells())
}

func (res OrganismResolver) Position() r2.Point {
	return res.organism.GetPosition()
}
