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
	points    int

	types    []CellType
	produces [][]int
}

const startingPoints = 30

var startingCellType CellType

func init() {
	diets := []Diet{Herbivore}
	startingCellType = CellType{
		ID:             0,
		diets:          diets,
		Herbivore:      10,
		wasteTolerance: 16,
		mobility:       10,
		points:         startingPoints,
	}
}

func (s Species) copy() Species {
	n := s

	n.produces = make([][]int, len(s.produces))
	n.types = make([]CellType, len(s.types))

	copy(n.produces, s.produces)

	for typeIndex := range s.types {
		n.types[typeIndex] = s.types[typeIndex].copy()
	}

	return n
}

func (s Species) getMaxTypes() int {
	if !s.types[len(s.types)-1].CanConnect() {
		return len(s.types)
	}

	return (s.points-startingPoints)/30 + 1
}

func (s Species) mutate() Species {
	n := s.copy()
	n.points++

	typeCount := len(n.types)
	typeIndex := rand.Intn(typeCount)
	mutatedType := n.types[typeIndex].copy().mutate()
	n.types[typeIndex] = mutatedType

	if s.getMaxTypes() > len(s.types) {
		ct := startingCellType.copy()
		ct.ID = s.types[len(s.types)-1].ID + 1
		for ct.points > ct.getInvestedPoints() {
			ct = ct.mutateOnce()
		}

		n.types = append(n.types, ct)
		prod := n.produces[typeCount-1]
		n.produces[typeCount-1] = append(prod, typeCount)
		n.produces = append(n.produces, []int{})
	}

	return n
}

func getRandomHerbivore() Species {
	ct := startingCellType.copy()

	for ct.points > ct.getInvestedPoints() {
		ct = ct.mutateOnce()
	}

	types := []CellType{ct}

	return Species{
		types:    types,
		produces: [][]int{{0}},
		points:   startingPoints,
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
			if !HasDiet(diet, diets) {
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
