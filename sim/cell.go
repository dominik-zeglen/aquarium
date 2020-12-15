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

	Shape string
	size  int

	membrane int
	enzymes  int

	herbivore int8
	carnivore int8
	funghi    int8

	alive          bool
	timeToDie      int
	hp             int
	bornAt         int
	diedAt         int
	procreatedAt   int
	wasteTolerance float64

	satiation    int
	maxSatiation int
	consumption  int
	transport    int
	capacity     int
	maxCapacity  int

	division      int8
	connectivity  int8
	procreationCd int8

	mobility int

	// True if cell reproduces by budding
	reproducingMethod bool
}

func (c Cell) getMaxHP() int {
	return c.size * 23
}

func (c Cell) getMass() int {
	return c.size * 10
}

func (c Cell) getFoodValue() int {
	return c.size * 2
}

func (c Cell) getDefence() int {
	return c.membrane / 10
}

func (c Cell) getAttack() int {
	return (int(c.carnivore)*2 + c.size + c.enzymes) / 3
}

func (c Cell) getLeftToFull() int {
	return c.maxSatiation - c.satiation
}

func (c Cell) shouldEat() bool {
	return c.maxSatiation > c.satiation
}

func (c Cell) getWaste(e Environment) float64 {
	waste := float64(c.size)
	if (e.toxicity) > 0 {
		waste -= float64(c.funghi) * 0.4
	}

	return waste / 6e8
}

func (c Cell) getWasteAfterDeath() float64 {
	return (float64(c.size)) / 8e8
}

func (c *Cell) eat(e Environment) {
	c.satiation -= c.consumption
	food := int(float32(c.herbivore) * 1.2)
	if e.toxicity > 0 {
		food += int(float32(c.funghi) * 0.4)
	}

	if food > c.getLeftToFull() {
		c.satiation = c.maxSatiation
		leftover := food - c.getLeftToFull()

		if leftover > c.maxCapacity {
			c.capacity = c.maxCapacity
		} else {
			c.capacity += leftover
		}

	} else {
		c.satiation += food
	}

	if c.shouldEat() {
		if c.capacity > c.getLeftToFull() {
			c.capacity -= c.getLeftToFull()
			c.satiation = c.maxSatiation
		} else {
			c.satiation += c.capacity
			c.capacity = 0
		}
	}
}

func (c *Cell) validate() {
	if c.carnivore < 0 {
		c.carnivore = 0
	}
	if c.herbivore < 0 {
		c.herbivore = 0
	}
	if c.funghi < 0 {
		c.funghi = 0
	}
	if c.consumption < 1 {
		c.consumption = 1
	}
	if c.division < 0 {
		c.division = 0
	}
}

func (c *Cell) mutate() {
	attr := rand.Intn(7)

	switch attr {
	case 0:
		c.carnivore += int8(rand.Intn(3) - 1)
		break
	case 1:
		c.herbivore += int8(rand.Intn(3) - 1)
		break
	case 2:
		c.funghi += int8(rand.Intn(3) - 1)
		break
	case 3:
		c.capacity += rand.Intn(3) - 1
		break
	case 4:
		c.consumption += rand.Intn(3) - 1
		break
	case 5:
		c.division += int8(rand.Intn(3) - 1)
		break
	case 6:
		c.wasteTolerance += float64(rand.Intn(3)-1) / 4
		break
	case 7:
		c.timeToDie += rand.Intn(3) - 1
		break
	}

	c.validate()
	if rand.Float32() > .999 {
		c.mutate()
	}
}

func (c Cell) canProcreate(iteration int) bool {
	if c.procreatedAt == 0 {
		return !c.shouldEat()
	}
	return iteration-c.procreatedAt > int(c.procreationCd) && !c.shouldEat()
}

func (c *Cell) procreate(canProcreate bool, iteration int, lastID int) []Cell {
	descendants := []Cell{}
	if canProcreate && c.canProcreate(iteration) && rand.Float32() > .8 {
		food := c.maxCapacity / int(c.division+1)

		for i := 0; i < int(c.division); i++ {
			descendant := *c
			descendant.id = lastID + i + 1
			descendant.satiation = food
			descendant.action = idle
			descendant.bornAt = iteration
			descendant.hp = c.getMaxHP()
			descendant.capacity = 0
			descendant.procreatedAt = iteration

			c.satiation = food

			if rand.Float32() > .99 {
				c.mutate()
			}

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
		Mul(float64(c.mobility / (*c).getMass()))
	c.position = c.position.Add(moveVec)
}

func (c Cell) shouldDie(env Environment, iteration int) bool {
	var prob float64
	if env.toxicity > c.wasteTolerance {
		prob = rand.Float64() + (env.toxicity - c.wasteTolerance)
	}

	p := gauss.NewGaussian(float64(c.timeToDie), 10)
	prob += p.Cdf(float64(iteration - c.bornAt))
	roll := rand.Float64()
	shouldDie := prob > roll

	return shouldDie
}

func (c *Cell) sim(env Environment, iteration int, lastID int, canProcreate bool) []Cell {
	descendants := []Cell{}

	if c.alive {
		c.eat(env)
		c.move()
		descendants = c.procreate(canProcreate, iteration, lastID)
		if c.shouldDie(env, iteration) {
			c.alive = false
			c.diedAt = iteration
		}
	}

	return descendants
}

func getRandomFunghiCell() Cell {
	points := int((17 * 100) / (3 + rand.Float64()))

	c := Cell{}

	c.capacity = rand.Intn(30)
	points -= c.capacity

	c.size = rand.Intn(30) + 1
	points -= c.size

	c.funghi = int8(rand.Intn(100))
	points -= int(c.funghi)

	c.division = 1
	points -= int(c.division * 10)

	c.timeToDie = int(rand.NormFloat64()*20) + 30
	points -= c.timeToDie * 10

	c.maxSatiation = int(rand.Intn(100)) + 300
	c.satiation = c.maxSatiation / 2

	c.consumption = 10
	c.procreationCd = int8(rand.Intn(4) + 8)

	c.wasteTolerance = rand.Float64()*3 + 10

	return c
}

func getRandomHerbivoreCell() Cell {
	points := int((17 * 100) / (3 + rand.Float64()))

	c := Cell{}

	c.capacity = rand.Intn(30)
	points -= c.capacity

	c.size = rand.Intn(30) + 1
	points -= c.size

	c.herbivore = int8(rand.Intn(100))
	points -= int(c.herbivore)

	c.division = 1
	points -= int(c.division * 10)

	c.timeToDie = int(rand.NormFloat64()*20) + 30
	points -= c.timeToDie * 10

	c.maxSatiation = int(rand.Intn(100)) + 300
	c.satiation = c.maxSatiation / 2

	c.consumption = 10
	c.procreationCd = int8(rand.Intn(4) + 8)

	c.wasteTolerance = rand.Float64()*6 + 4

	return c
}

func getRandomCell(id int) Cell {
	var c Cell

	if rand.Float32() > .5 {
		c = getRandomHerbivoreCell()
	} else {
		c = getRandomFunghiCell()
	}

	c.id = id
	c.action = idle
	c.alive = true
	c.bornAt = 0

	return c
}
