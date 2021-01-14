package sim

import (
	"testing"
)

func TestDietMutation(t *testing.T) {
	t.Run("returns immutable", func(t *testing.T) {
		// Given
		ct := CellType{
			diets: []Diet{Herbivore},
		}

		// When
		newType := ct.copy()
		newType.mutateDiet()

		// Then
		if &newType == &ct {
			t.Errorf("New cell type should be a copy, not reference")
		}
		if &newType.diets == &ct.diets {
			t.Errorf("New cell type should be a copy, not reference")
		}
	})

	t.Run("removes one diet if has two", func(t *testing.T) {
		// Given
		ct := CellType{
			diets: []Diet{Herbivore, Funghi},
		}

		// When
		newType := ct.copy()
		newType.mutateDiet()

		// Then
		if len(newType.diets) != 1 {
			t.Errorf("New cell type should have only one diet")
		}
	})
}
