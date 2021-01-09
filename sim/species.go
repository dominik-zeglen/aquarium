package sim

import (
	"fmt"
	"math/rand"
)

type Species struct {
	id        int
	emergedAt int
	extinct   bool
	count     int

	types    []CellType
	produces [][]int
}

func (s Species) mutate() Species {
	n := s

	n.produces = make([][]int, len(s.produces))
	n.types = make([]CellType, len(s.types))

	copy(n.produces, s.produces)
	copy(n.types, s.types)

	typeIndex := rand.Intn(len(n.types))
	mutatedType := n.types[typeIndex].mutate()
	n.types[typeIndex] = mutatedType

	return n
}

func getRandomHerbivore() Species {
	return Species{
		types: []CellType{{
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
			mobility:       100,
		}},
		produces: [][]int{{0}},
	}
}

// Getters

func (s Species) GetID() int {
	return s.id
}

func (s Species) GetEmergedAt() int {
	return s.emergedAt
}

func (s Species) GetDiets() []Diet {
	diets := []Diet{}

	for _, cellType := range s.types {
		for _, diet := range cellType.diets {
			if !hasDiet(diet, diets) {
				diets = append(diets, diet)
			}
		}
	}

	return diets
}

func (s Species) GetName() string {
	name := ""

	for _, diet := range s.GetDiets() {
		switch diet {
		case Funghi:
			name += "F"
			break
		case Herbivore:
			name += "H"
			break
		}
	}

	return fmt.Sprintf("%s-%d-%d", name, s.emergedAt, s.id)
}

func (s Species) GetTypes() []CellType {
	return s.types
}

type SpeciesList []Species

func (sl SpeciesList) GetAlive() SpeciesList {
	species := make(SpeciesList, len(sl))

	index := 0
	for speciesIndex := range sl {
		if !sl[speciesIndex].extinct {
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
					if species[speciesIndex].id == cellSpecies.id {
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
