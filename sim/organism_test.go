package sim

import (
	"testing"

	"github.com/golang/geo/r2"
)

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
	os := o.split()

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
		maxSatiation: 400,
		size:         10,
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
	expected := 13336
	if left != expected {
		t.Errorf("Expected %d, got %d", expected, left)
	}
}

func TestOrganismDyingFromHunger(t *testing.T) {
	// Given
	env := Environment{0, 10, 10}
	ct := CellType{
		consumption:  100,
		maxSatiation: 200,
		size:         10,
		TimeToDie:    2,
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
		maxSatiation: 200,
		size:         10,
		TimeToDie:    2,
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
		maxSatiation: 200,
		size:         10,
		TimeToDie:    2,
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
		maxSatiation: 200,
		size:         10,
		TimeToDie:    2,
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
	o.killCells(env, 3)

	// Then
	expected := 1
	result := o.cells.GetAliveCount()
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
