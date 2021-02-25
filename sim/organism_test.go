package sim

import (
	"context"
	"testing"

	"github.com/golang/geo/r2"
)

func addSpecies(s Species) *Species {
	return &s
}

func TestOrganismSplitting(t *testing.T) {
	// Given
	o := Organism{
		action: idle,
		cells: CellList{{
			id:       0,
			position: r2.Point{0, 0},
		}, {
			id:       1,
			position: r2.Point{1, 0},
		}, {
			id:       2,
			position: r2.Point{1, 1},
		}, {
			id:       3,
			position: r2.Point{2, 2},
		}},
	}

	// When
	os := o.split(context.TODO(), true, 1)

	// Then
	if len(os) != 1 {
		t.Errorf("Expected 1, got %d", len(os))
	}

	smaller := os[0]
	if len(os[0].cells) > len(o.cells) {
		smaller = o
	}

	if len(smaller.cells) != 1 {
		t.Errorf("Expected 1, got %d", len(smaller.cells))
	}

	if (smaller.action) != idle {
		t.Errorf("Expected idle, got %s", smaller.action)
	}

	expectedPosition := r2.Point{X: 2, Y: 2}
	if smaller.position != expectedPosition {
		t.Errorf("Expected organism position at (2, 2), got (%.f, %.f)", smaller.position.X, smaller.position.Y)
	}

	cell := smaller.cells[0]
	expectedCellPosition := r2.Point{}
	if cell.position != expectedCellPosition {
		t.Errorf("Expected cell position at (0, 0), got (%.f, %.f)", cell.position.X, cell.position.Y)
	}

	if cell.id != 0 {
		t.Errorf("Expected cell ID to be 0, got %d", cell.id)
	}
}

func TestOrganismEating(t *testing.T) {
	// Given
	env := Environment{0, 10, 10}
	ct := CellType{
		consumption:  100,
		Herbivore:    100,
		diets:        []Diet{Herbivore},
		maxCapacity:  200,
		maxSatiation: -50,
	}
	s := Species{
		produces: [][]int{{0}},
	}
	o := Organism{
		species: &s,
		cells: CellList{{
			id:       0,
			alive:    true,
			cellType: &ct,
		}, {
			id:       1,
			alive:    true,
			cellType: &ct,
		}, {
			id:       2,
			alive:    true,
			cellType: &ct,
		}, {
			id:       3,
			alive:    true,
			cellType: &ct,
		}},
	}

	// When
	left := o.eat(env, 0)

	// Then
	expected := 3736
	if left != expected {
		t.Errorf("Expected %d, got %d", expected, left)
	}
}

func TestOrganismDyingFromHunger(t *testing.T) {
	// Given
	env := Environment{0, 10, 10}
	ct := CellType{
		consumption:  100,
		maxSatiation: 150,

		timeToDie: 10,
	}
	s := Species{
		produces: [][]int{{0}},
	}
	o := Organism{
		species: &s,
		cells: CellList{{
			id:        0,
			alive:     true,
			cellType:  &ct,
			hp:        1,
			satiation: 0,
		}, {
			id:        1,
			alive:     true,
			cellType:  &ct,
			hp:        1,
			satiation: 200,
		}},
	}

	// When
	for cIndex := range o.cells {
		o.cells[cIndex].consume()
	}
	o.killCells(env, 1)

	// Then
	expected := 1
	result := o.cells.GetAliveCount()
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestOrganismDyingFromOutOfBounds(t *testing.T) {
	// Given
	env := Environment{0, 10, 10}
	ct := CellType{
		consumption:  100,
		maxSatiation: 150,

		timeToDie: 10,
	}
	s := Species{
		produces: [][]int{{0}},
	}
	o := Organism{
		species: &s,
		cells: CellList{{
			id:        0,
			alive:     true,
			cellType:  &ct,
			hp:        1,
			position:  r2.Point{1, 1},
			satiation: 1,
		}, {
			id:        1,
			alive:     true,
			cellType:  &ct,
			hp:        1,
			position:  r2.Point{-1, 1},
			satiation: 1,
		}},
	}

	// When
	o.killCells(env, 1)

	// Then
	expected := 1
	result := o.cells.GetAliveCount()
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestOrganismDyingFromToxicity(t *testing.T) {
	// Given
	env := Environment{1, 10, 10}
	ct := CellType{
		consumption:  100,
		maxSatiation: 150,

		timeToDie: 10,
	}
	s := Species{
		produces: [][]int{{0}},
	}
	o := Organism{
		species: &s,
		cells: CellList{{
			id:        0,
			alive:     true,
			cellType:  &ct,
			hp:        1,
			satiation: 1,
		}, {
			id:        1,
			alive:     true,
			cellType:  &ct,
			hp:        1,
			satiation: 1,
		}},
	}

	// When
	o.killCells(env, 1)

	// Then
	expected := 0
	result := o.cells.GetAliveCount()
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestOrganismDyingFromAge(t *testing.T) {
	// Given
	env := Environment{0, 10, 10}
	ct := CellType{
		consumption:  100,
		maxSatiation: 150,

		timeToDie: 10,
	}
	s := Species{
		produces: [][]int{{0}},
	}
	o := Organism{
		species: &s,
		cells: CellList{{
			id:        0,
			alive:     true,
			bornAt:    2,
			cellType:  &ct,
			hp:        1,
			satiation: 1,
		}, {
			id:        1,
			alive:     true,
			cellType:  &ct,
			hp:        1,
			satiation: 1,
		}},
	}

	// When
	o.killCells(env, 43)

	// Then
	expected := 1
	result := o.cells.GetAliveCount()
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestOrganismMutation(t *testing.T) {
	t.Run("creates copy of species", func(t *testing.T) {
		// Given
		ct := CellType{
			consumption:  100,
			diets:        []Diet{Funghi},
			maxSatiation: 150,

			timeToDie: 10,
		}
		s := Species{
			produces: [][]int{{0}},
			types:    []CellType{ct},
		}
		o1 := Organism{
			species: &s,
			cells: CellList{{
				id:        0,
				alive:     true,
				bornAt:    2,
				cellType:  &s.types[0],
				hp:        1,
				satiation: 1,
			}},
		}

		// When
		o1.procreate(true, 1, 25, true)
		os := o1.split(context.TODO(), true, 1)
		o2 := os[0]
		o2.mutate(addSpecies)
		o2.species.types[0].mutateDiet()

		// Then
		if o1.species == o2.species {
			t.Error("Does not create copy")
		}

		for ctIndex := range o1.species.types {
			if &o1.species.types[ctIndex] == &o2.species.types[ctIndex] {
				t.Error("Does not create copy")
			}
		}

		if o1.cells[0].cellType == o2.cells[0].cellType {
			t.Error("Does not create copy")
		}

		if len(o1.cells[0].cellType.diets) != 1 || o1.cells[0].cellType.diets[0] != Funghi {
			t.Error("Does not create copy")
		}

		if o2.cells[0].cellType.diets[0] == Funghi && len(o2.cells[0].cellType.diets) == 1 {
			t.Error("Does not create copy")
		}
	})
}

func TestRandomOrganism(t *testing.T) {
	t.Run("creates copy of species", func(t *testing.T) {
		// Given
		s := getRandomHerbivore()
		s.types[0].connects = 0
		o1 := Organism{
			species: &s,
			cells: CellList{{
				id:        0,
				alive:     true,
				bornAt:    2,
				cellType:  &s.types[0],
				hp:        1,
				satiation: 1,
			}},
		}

		// When
		o1.procreate(true, 1, 25, true)
		os := o1.split(context.TODO(), true, 1)
		o2 := os[0]
		o2.mutate(addSpecies)
		o2.species.types[0].mutateDiet()

		// Then
		if o1.species == o2.species {
			t.Error("Does not create copy")
		}

		for ctIndex := range o1.species.types {
			if &o1.species.types[ctIndex] == &o2.species.types[ctIndex] {
				t.Error("Does not create copy")
			}
		}

		if o1.cells[0].cellType == o2.cells[0].cellType {
			t.Error("Does not create copy")
		}

		if &o1.cells[0].cellType.diets == &o2.cells[0].cellType.diets {
			t.Error("Does not create copy")
		}

		if len(o1.cells[0].cellType.diets) != 1 || o1.cells[0].cellType.diets[0] != Herbivore {
			t.Error("Does not create copy")
		}

		if o2.cells[0].cellType.diets[0] == Herbivore && len(o2.cells[0].cellType.diets) == 1 {
			t.Errorf("Does not create copy, %s vs %s", o1.cells[0].cellType.diets, o2.cells[0].cellType.diets)
		}
	})
}

func TestOrganismProcreation(t *testing.T) {
	t.Run("can produce more than 1 cell type", func(t *testing.T) {
		// Given
		cts := []CellType{
			{
				ID:           0,
				consumption:  100,
				diets:        []Diet{Funghi},
				maxSatiation: 150,
				timeToDie:    10,
			},
			{
				ID:           1,
				consumption:  100,
				diets:        []Diet{Herbivore},
				maxSatiation: 150,
				timeToDie:    10,
			},
		}
		s := Species{
			produces: [][]int{{0, 1}, {}},
			types:    cts,
		}
		o := Organism{
			species: &s,
			cells: CellList{{
				id:        0,
				alive:     true,
				cellType:  &cts[0],
				hp:        1,
				satiation: 0,
			}, {
				id:        1,
				alive:     true,
				cellType:  &cts[0],
				hp:        1,
				satiation: 200,
			}},
		}

		// When
		for i := 0; i < 50; i++ {
			o.procreate(true, 0, 50, true)
		}

		// Then
		found := false
		for _, cell := range o.cells {
			if cell.cellType.ID == 1 {
				found = true
				break
			}
		}

		if !found {
			t.Error("Organism did not produce any type 1 cell")
		}
	})
}
