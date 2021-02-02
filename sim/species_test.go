package sim

import (
	"testing"
)

func TestMutation(t *testing.T) {
	t.Run("returns immutable", func(t *testing.T) {
		// Given
		ct := CellType{
			consumption:  100,
			maxSatiation: 150,
			timeToDie:    10,
		}
		s := Species{
			produces: [][]int{{0}},
			types:    []CellType{ct},
		}

		// When
		newSpecies := s.mutate()

		// Then
		if &newSpecies == &s {
			t.Errorf("New species should be a copy, not reference")
		}
		if &newSpecies.produces == &s.produces {
			t.Errorf("New species should be a copy, not reference")
		}
		if &newSpecies.types == &s.types {
			t.Errorf("New species should be a copy, not reference")
		}
		for cellTypeIndex := range newSpecies.types {
			if &newSpecies.types[cellTypeIndex] == &s.types[cellTypeIndex] {
				t.Errorf("New species should be a copy, not reference")
			}
		}
	})

	t.Run("is able to create multiple cell types", func(t *testing.T) {
		// Given
		ct := CellType{
			consumption:  100,
			maxSatiation: 150,
			timeToDie:    10,
			connects:     15,
		}
		s := Species{
			produces: [][]int{{0}},
			types:    []CellType{ct},
			points:   100,
		}

		// When
		newSpecies := s.mutate()

		// Then
		if newSpecies.getMaxTypes() != 2 {
			t.Errorf("New species should be able to hold 2 types, got: %d", newSpecies.getMaxTypes())
		}
		if len(newSpecies.types) != 2 {
			t.Errorf("New species should have 2 types, got: %d", len(newSpecies.types))
		}
		if len(newSpecies.produces[0]) != 2 {
			t.Errorf("First type should be able to produce 2 types, got: %d", len(newSpecies.produces[0]))
		}
	})
}
