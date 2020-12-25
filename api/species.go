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
func (res SpeciesResolver) Cells() CellConnectionResolver {
	cells := make([]sim.Cell, res.species.Count)
	simCells := res.s.GetCells()
	index := 0
	for cellIndex := range simCells {
		if simCells[cellIndex].GetSpecies().ID == res.species.ID {
			cells[index] = simCells[cellIndex]
		}
	}

	return CreateCellConnectionResolver(cells, res.s)
}

type SpeciesConnectionEdgeResolver struct {
	species sim.Species
	s       *sim.Sim
}

func CreateSpeciesConnectionEdgeResolver(species sim.Species, sim *sim.Sim) SpeciesConnectionEdgeResolver {
	return SpeciesConnectionEdgeResolver{species, sim}
}

func (res SpeciesConnectionEdgeResolver) Node() SpeciesResolver {
	return CreateSpeciesResolver(&res.species, res.s)
}

type SpeciesConnectionResolver struct {
	species []sim.Species
	s       *sim.Sim
}

func CreateSpeciesConnectionResolver(species []sim.Species, sim *sim.Sim) SpeciesConnectionResolver {
	return SpeciesConnectionResolver{species, sim}
}

func (res SpeciesConnectionResolver) Count() int32 {
	return int32(len(res.species))
}

func (res SpeciesConnectionResolver) Edges() []SpeciesConnectionEdgeResolver {
	resolvers := make([]SpeciesConnectionEdgeResolver, len(res.species))

	for speciesIndex := range res.species {
		resolvers[speciesIndex] = CreateSpeciesConnectionEdgeResolver(res.species[speciesIndex], res.s)
	}

	return resolvers
}
