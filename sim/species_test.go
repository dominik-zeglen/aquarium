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
			size:         10,
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
}
