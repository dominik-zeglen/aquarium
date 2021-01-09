package api

import (
	"context"

	"github.com/dominik-zeglen/aquarium/middleware"
	"github.com/dominik-zeglen/aquarium/sim"
	"github.com/golang/geo/r2"
)

type SpeciesResolver struct {
	species sim.Species
}

func createSpeciesResolverList(species sim.SpeciesList) []SpeciesResolver {
	resolvers := make([]SpeciesResolver, len(species))

	for speciesIndex := range species {
		resolvers[speciesIndex] = SpeciesResolver{species[speciesIndex]}
	}

	return resolvers
}

func (res SpeciesResolver) ID() int32 {
	return int32(res.species.GetID())
}

func (res SpeciesResolver) Name() string {
	return res.species.GetName()
}

func (res SpeciesResolver) EmergedAt() int32 {
	return int32(0)
}
func (res SpeciesResolver) Diet() []string {
	diets := res.species.GetDiets()
	dietNames := make([]string, len(diets))

	for dietIndex, diet := range diets {
		dietNames[dietIndex] = diet.String()
	}

	return dietNames
}
func (res SpeciesResolver) Organisms(ctx context.Context) []OrganismResolver {
	s := ctx.Value(middleware.SimContextKey).(*sim.Sim)
	organisms := s.GetOrganisms().GetSpecies(res.species.GetID())

	return createOrganismResolverList(organisms)
}
func (res SpeciesResolver) CellTypes() []CellTypeResolver {
	return createCellTypeResolverList(res.species.GetTypes())
}

type SpeciesGridElementResolver struct {
	Position r2.Point
	species  []sim.Species
	s        *sim.Sim
}

func CreateSpeciesGridElementResolver(
	position r2.Point,
	species []sim.Species,
	sim *sim.Sim,
) SpeciesGridElementResolver {
	return SpeciesGridElementResolver{position, species, sim}
}

func (res SpeciesGridElementResolver) Species() []SpeciesResolver {
	resolvers := make([]SpeciesResolver, len(res.species))

	for speciesIndex := range res.species {
		resolvers[speciesIndex] = SpeciesResolver{res.species[speciesIndex]}
	}

	return resolvers
}
