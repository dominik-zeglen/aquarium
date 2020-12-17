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

func (c *Cell) eat(e Environment) {
	c.satiation -= c.species.consumption
	food := int(float32(c.species.Herbivore) * 1.2)
	if e.toxicity > 0 {
		food += int(c.species.getProcessedWaste(e))
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
	addSpecies AddSpecies,
) []Cell {
	descendants := []Cell{}

	if canProcreate && c.canProcreate(iteration) && rand.Float32() > .8 {
		food := c.species.maxCapacity / int(c.species.division+1)

		for i := 0; i < int(c.species.division); i++ {
			descendant := *c
			descendant.id = lastID + i + 1
			descendant.satiation = food
			descendant.action = idle
			descendant.bornAt = iteration
			descendant.capacity = 0
			descendant.procreatedAt = iteration

			c.satiation = food

			if rand.Float32() > .99 {
				species := c.species.mutate()
				species.emergedAt = iteration
				c.species = addSpecies(species)
			} else {
				descendant.species = c.species
			}

			descendant.hp = c.species.getMaxHP()
			descendants = append(descendants, descendant)
		}
	}

	return descendants
}

func (c *Cell) move() {
	moveVec := c.
		target.
		Sub(c.position).
		Normalize().
		Mul(float64(c.species.mobility / c.species.getMass()))
	c.position = c.position.Add(moveVec)
}

func (c Cell) shouldDie(env Environment, iteration int) bool {
	var prob float64
	if env.toxicity > c.species.wasteTolerance {
		prob = rand.Float64() + (env.toxicity - c.species.wasteTolerance)
	}

	p := gauss.NewGaussian(float64(c.species.timeToDie), 10)
	prob += p.Cdf(float64(iteration - c.bornAt))
	roll := rand.Float64()
	shouldDie := prob > roll

	return shouldDie
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
		c.eat(env)
		c.move()
		descendants = c.procreate(canProcreate, iteration, lastID, addSpecies)
		if c.shouldDie(env, iteration) {
			c.alive = false
			c.diedAt = iteration
		}
	}

	return descendants
}

func getRandomCell(id int, addSpecies AddSpecies) Cell {
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
	c.satiation = c.species.maxSatiation / 2

	return c
}
