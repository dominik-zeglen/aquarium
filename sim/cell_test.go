package sim

import (
	"testing"
)

func TestEating(t *testing.T) {
	e := Environment{9, 100, 100}

	t.Run("Simple eating by herbivore", func(t *testing.T) {
		// Given
		satiation := 5
		s := CellType{
			consumption:  10,
			maxSatiation: 30,
			Herbivore:    30,
		}
		c := Cell{
			cellType:  &s,
			satiation: satiation,
		}

		// When
		food := c.GetFood(e, 0, 0)
		c.eat(food)

		// Then
		if c.satiation < satiation {
			t.Errorf(
				"Cell is not eating, before: %d, after: %d",
				satiation,
				c.satiation,
			)
		}
	})

	t.Run("Simple eating by funghi", func(t *testing.T) {
		// Given
		satiation := 5
		s := CellType{
			consumption:  10,
			maxSatiation: 30,
			Funghi:       30,
		}
		c := Cell{
			cellType:  &s,
			satiation: satiation,
		}

		// When
		food := c.GetFood(e, 0, 0)
		c.eat(food)

		//Then
		if c.satiation < satiation {
			t.Errorf(
				"Cell is not eating, before: %d, after: %d",
				satiation,
				c.satiation,
			)
		}
	})

	t.Run("Simple eating by carnivore", func(t *testing.T) {
		// Given
		satiation := 5
		s := CellType{
			consumption:  10,
			maxSatiation: 30,
			Carnivore:    30,
		}
		c := Cell{
			cellType:  &s,
			satiation: satiation,
		}

		// When
		food := c.GetFood(e, 0, 0)
		c.eat(food)

		//Then
		if c.satiation > satiation {
			t.Errorf(
				"Cell is eating, before: %d, after: %d",
				satiation,
				c.satiation,
			)
		}
	})
}

func TestProcreation(t *testing.T) {
	// Given
	ct := CellType{
		consumption:  100,
		maxSatiation: 150,
		size:         10,
		timeToDie:    10,
	}
	c := Cell{
		id:        0,
		alive:     true,
		bornAt:    2,
		cellType:  &ct,
		hp:        1,
		satiation: 200,
	}

	// When
	child := c.procreate(10, []*CellType{&ct})

	// Then
	if !child.alive {
		t.Error("Child is not alive")
	}
	if child.satiation <= 0 {
		t.Error("Child is not fed")
	}
	if child.position.X == c.position.X {
		t.Error("Child is spawned in the same position")
	}
}
