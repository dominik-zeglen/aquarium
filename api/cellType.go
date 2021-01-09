package api

import (
	"github.com/dominik-zeglen/aquarium/sim"
)

type CellTypeResolver struct {
	cellType sim.CellType
}

func createCellTypeResolverList(species []sim.CellType) []CellTypeResolver {
	resolvers := make([]CellTypeResolver, len(species))

	for speciesIndex := range species {
		resolvers[speciesIndex] = CellTypeResolver{species[speciesIndex]}
	}

	return resolvers
}

func (res CellTypeResolver) ID() int32 {
	return int32(res.cellType.ID)
}

func (res CellTypeResolver) Diet() []string {
	diets := res.cellType.GetDiet()
	dietNames := make([]string, len(diets))

	for dietIndex, diet := range diets {
		dietNames[dietIndex] = diet.String()
	}

	return dietNames
}
func (res CellTypeResolver) Funghi() int32 {
	return int32(res.cellType.Funghi)
}
func (res CellTypeResolver) Herbivore() int32 {
	return int32(res.cellType.Herbivore)
}
