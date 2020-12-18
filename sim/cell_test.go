package sim

import "testing"

func TestEating(t *testing.T) {
	e := Environment{9}

	t.Run("Simple eating by herbivore", func(t *testing.T) {
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

		c.eat(e)

		if c.satiation < satiation {
			t.Errorf(
				"Cell is not eating, before: %d, after: %d",
				satiation,
				c.satiation,
			)
		}
	})

	t.Run("Simple eating by funghi", func(t *testing.T) {
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

		c.eat(e)

		if c.satiation < satiation {
			t.Errorf(
				"Cell is not eating, before: %d, after: %d",
				satiation,
				c.satiation,
			)
		}
	})

	t.Run("Simple eating by carnivore", func(t *testing.T) {
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

		c.eat(e)

		if c.satiation > satiation {
			t.Errorf(
				"Cell is eating, before: %d, after: %d",
				satiation,
				c.satiation,
			)
		}
	})
}
