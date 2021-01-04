package sim

import (
	"math/rand"
)

type Species struct {
	ID        int  `json:"id"`
	EmergedAt int  `json:"emergedAt"`
	Extinct   bool `json:"extinct"`
	Count     int  `json:"count"`

	Types    []CellType
	Produces [][]int
}

func (s Species) mutate() Species {
	n := s

	typeIndex := rand.Intn(len(n.Types))
	mutatedType := n.Types[typeIndex].mutate()
	n.Types[typeIndex] = mutatedType

	return n
}

func getRandomHerbivore() Species {
	return Species{
		Types: []CellType{{
			ID:             0,
			diets:          []Diet{Herbivore},
			maxCapacity:    rand.Intn(30),
			size:           rand.Intn(20) + 10,
			Herbivore:      int8(rand.Intn(20)) + 5,
			TimeToDie:      40,
			maxSatiation:   int(rand.Intn(100)) + 300,
			consumption:    10,
			procreationCd:  int8(rand.Intn(4) + 8),
			WasteTolerance: float64(rand.Intn(16))/4 + 4,
			mobility:       20,
		}},
		Produces: [][]int{{0}},
	}
}

type SpeciesList []Species

func (sl SpeciesList) GetAlive() SpeciesList {
	species := make(SpeciesList, len(sl))

	index := 0
	for speciesIndex := range sl {
		if !sl[speciesIndex].Extinct {
			species[index] = sl[speciesIndex]
			index++
		}
	}

	return species[:index]

}

type SpeciesGridRow = map[int]SpeciesList
type SpeciesGrid = map[int]SpeciesGridRow

func (sl SpeciesList) GetArea(organisms OrganismList, scale int) SpeciesGrid {
	grid := SpeciesGrid{}

	for cellIndex := range organisms {
		x := int(organisms[cellIndex].position.X) / scale
		y := int(organisms[cellIndex].position.Y) / scale
		cellSpecies := *organisms[cellIndex].species

		_, ok := grid[y]
		if !ok {
			grid[y] = SpeciesGridRow{}
			grid[y][x] = []Species{cellSpecies}
		} else {
			species, ok := grid[y][x]
			if !ok {
				grid[y][x] = []Species{cellSpecies}
			} else {
				found := false
				for speciesIndex := range species {
					if species[speciesIndex].ID == cellSpecies.ID {
						found = true
						break
					}
				}

				if !found {
					grid[y][x] = append(species, cellSpecies)
				}
			}
		}

	}

	return grid
}
