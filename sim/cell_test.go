package sim

import (
	"testing"
)

var env Environment

func addSpecies(species Species) *Species {
	return &species
}

func init() {
	env = Environment{9, 100, 100}
}

func TestEating(t *testing.T) {
	e := Environment{9, 100, 100}

	t.Run("Simple eating by herbivore", func(t *testing.T) {
		// Given
		satiation := 5
		s := Species{
			consumption:  10,
			maxSatiation: 30,
			Herbivore:    30,
		}
		c := Cell{
			species:   &s,
			satiation: satiation,
		}

		// When
		c.eat(e, 0)

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
		s := Species{
			consumption:  10,
			maxSatiation: 30,
			Funghi:       30,
		}
		c := Cell{
			species:   &s,
			satiation: satiation,
		}

		// When
		c.eat(e, 0)

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
		s := Species{
			consumption:  10,
			maxSatiation: 30,
			Carnivore:    30,
		}
		c := Cell{
			species:   &s,
			satiation: satiation,
		}

		// When
		c.eat(e, 0)

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

func BenchmarkEatHerbivore(b *testing.B) {
	c := getRandomCell(0, env, addSpecies)

	for i := 0; i < b.N; i++ {
		c.eat(env, 1)
	}
}

func BenchmarkEatFunghi(b *testing.B) {
	c := getRandomCell(0, env, addSpecies)
	c.species.Funghi = c.species.Herbivore
	c.species.Herbivore = 0

	for i := 0; i < b.N; i++ {
		c.eat(env, 1)
	}
}

func BenchmarkEatOmnivore(b *testing.B) {
	c := getRandomCell(0, env, addSpecies)
	c.species.Funghi = c.species.Herbivore

	for i := 0; i < b.N; i++ {
		c.eat(env, 1)
	}
}

func BenchmarkProcreate(b *testing.B) {
	c := getRandomCell(0, env, addSpecies)
	c.satiation = c.species.maxSatiation

	for i := 0; i < b.N; i++ {
		c.procreate(true, 1, 1, env, addSpecies)
	}
}
