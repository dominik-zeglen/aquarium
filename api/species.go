package api

import (
	"github.com/dominik-zeglen/aquarium/sim"
)

type SpeciesResolver struct {
	species *sim.Species
	s       *sim.Sim
}

func CreateSpeciesResolver(species *sim.Species, sim *sim.Sim) SpeciesResolver {
	return SpeciesResolver{species, sim}
}

func (res SpeciesResolver) ID() int32 {
	return int32(res.species.ID)
}

func (res SpeciesResolver) Consumption() int32 {
	return int32(res.species.GetConsumption())
}

func (res SpeciesResolver) Name() string {
	return res.species.GetName()
}

func (res SpeciesResolver) EmergedAt() int32 {
	return int32(res.species.EmergedAt)
}
func (res SpeciesResolver) Diet() []string {
	var diets []string
	speciesDiets := res.species.GetDiet()

	for _, diet := range speciesDiets {
		diets = append(diets, diet.String())
	}

	return diets
}
func (res SpeciesResolver) Carnivore() int32 {
	return int32(res.species.Carnivore)
}
func (res SpeciesResolver) Herbivore() int32 {
	return int32(res.species.Herbivore)
}
func (res SpeciesResolver) Funghi() int32 {
	return int32(res.species.Funghi)
}
func (res SpeciesResolver) Cells() []CellResolver {
	var resolvers []CellResolver
	simCells := res.s.GetCells()

	for _, cell := range simCells {
		if cell.GetSpecies().ID == res.species.ID {
			resolvers = append(resolvers, CreateCellResolver(&cell, res.s))
		}
	}

	return resolvers
}
