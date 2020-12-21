package sim

import (
	"math/rand"

	gauss "github.com/chobie/go-gaussian"
	"github.com/golang/geo/r2"
)

type Cell struct {
	id       int
	position r2.Point
	action   Action
	target   r2.Point

	species *Species

	alive        bool
	hp           int
	bornAt       int
	diedAt       int
	procreatedAt int

	satiation int
	capacity  int
}

func (c Cell) getLeftToFull() int {
	return c.species.maxSatiation - c.satiation
}

func (c Cell) shouldEat() bool {
	return c.species.maxSatiation > c.satiation
}

func (c *Cell) eat(e Environment, iteration int) {
	c.satiation -= c.species.getConsumption()
	food := 0
	if c.species.Herbivore > 0 {
		food += int(float64(c.species.Herbivore) * e.getLightOnHeight(c.position.Y, iteration))
	}
	if c.species.Funghi > 0 {
		food += int(c.species.getProcessedWaste(e.getToxicityOnHeight(c.position.Y)))
	}

	if food > c.getLeftToFull() {
		c.satiation = c.species.maxSatiation
		leftover := food - c.getLeftToFull()

		if leftover > c.species.maxCapacity {
			c.capacity = c.species.maxCapacity
		} else {
			c.capacity += leftover
		}

	} else {
		c.satiation += food
	}

	if c.shouldEat() {
		if c.capacity > c.getLeftToFull() {
			c.capacity -= c.getLeftToFull()
			c.satiation = c.species.maxSatiation
		} else {
			c.satiation += c.capacity
			c.capacity = 0
		}
	}
}

func (c Cell) canProcreate(iteration int) bool {
	if c.procreatedAt == 0 {
		return !c.shouldEat()
	}
	return iteration-c.procreatedAt > int(c.species.procreationCd) && !c.shouldEat()
}

func (c *Cell) procreate(
	canProcreate bool,
	iteration int,
	lastID int,
	env Environment,
	addSpecies AddSpecies,
) []Cell {
	descendants := []Cell{}

	if canProcreate && c.canProcreate(iteration) && rand.Float32() > .7 {
		food := c.species.maxCapacity / int(c.species.division+1)

		for i := 0; i < int(c.species.division); i++ {
			descendant := Cell{}
			descendant.id = lastID + i + 1
			descendant.satiation = food
			descendant.action = idle
			descendant.bornAt = iteration
			descendant.capacity = 0
			descendant.procreatedAt = iteration
			descendant.alive = true

			vec := getRandomVec().Mul(float64(c.species.size))

			descendant.position = c.position.Add(vec)
			c.position = c.position.Sub(vec)

			c.satiation = food
			c.bornAt = iteration

			if rand.Float32() > .99 {
				species := c.species.mutate()
				species.EmergedAt = iteration
				c.species = addSpecies(species)
			}

			descendant.species = c.species
			descendant.hp = descendant.species.getMaxHP()

			descendants = append(descendants, descendant)
		}
	}

	return descendants
}

func (c *Cell) move() {
	var moveVec r2.Point

	if c.action == idle {
		moveVec = getRandomVec()
	} else {
		moveVec = c.
			target.
			Sub(c.position).
			Normalize()
	}

	moveVec = moveVec.Mul(float64(c.species.mobility / c.species.getMass()))
	c.position = c.position.Add(moveVec)
}

func (c Cell) shouldDie(env Environment, iteration int) bool {
	if c.satiation == 0 || isOutOfBounds(c.position, env) {
		return true
	}

	var prob float64
	if env.toxicity > c.species.WasteTolerance {
		prob = rand.Float64() + (env.getToxicityOnHeight(c.position.Y) - c.species.WasteTolerance)
	}
	if c.bornAt+c.species.TimeToDie-iteration < 0 {
		p := gauss.NewGaussian(float64(c.species.TimeToDie), 10)
		prob += p.Cdf(float64(iteration - c.bornAt))
	}

	roll := rand.Float64()
	result := prob > roll

	return result
}

func (c *Cell) sim(
	env Environment,
	iteration int,
	lastID int,
	addSpecies AddSpecies,
	canProcreate bool,
) []Cell {
	descendants := []Cell{}

	if c.alive {
		c.eat(env, iteration)
		c.move()
		descendants = c.procreate(canProcreate, iteration, lastID, env, addSpecies)
		if c.shouldDie(env, iteration) {
			c.alive = false
			c.diedAt = iteration
		}
	}

	return descendants
}

func getRandomCell(id int, e Environment, addSpecies AddSpecies) Cell {
	var c Cell
	var s Species

	if rand.Float32() > .5 {
		s = getRandomHerbivore()
	} else {
		s = getRandomFunghi()
	}

	c.species = addSpecies(s)

	c.id = id
	c.action = idle
	c.alive = true
	c.satiation = 20
	c.position = r2.Point{
		float64(e.width) * rand.Float64(),
		float64(e.height) * rand.Float64(),
	}

	return c
}
