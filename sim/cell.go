package sim

import (
	"math/rand"

	gauss "github.com/chobie/go-gaussian"
	"github.com/golang/geo/r2"
)

const procreationCd = 10

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

	alive        bool
	timeToDie    int
	hp           int
	bornAt       int
	diedAt       int
	procreatedAt int

	satiation    int
	maxSatiation int
	consumption  int
	extrements   int
	transport    int
	capacity     int
	maxCapacity  int

	division     int8
	connectivity int8

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

func (c *Cell) eat() {
	c.satiation -= c.consumption
	food := int(c.herbivore * 2)

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

func (c Cell) canProcreate(iteration int) bool {
	if c.procreatedAt == 0 {
		return !c.shouldEat()
	}
	return iteration-c.procreatedAt > procreationCd && !c.shouldEat()
}

func (c *Cell) procreate(iteration int, lastID int) []Cell {
	descendants := []Cell{}
	if c.canProcreate(iteration) {
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

func (c Cell) shouldDie(iteration int) bool {
	p := gauss.NewGaussian(float64(c.timeToDie), 4)
	prob := p.Cdf(float64(iteration - c.bornAt))
	roll := rand.Float64()
	shouldDie := prob > roll

	return shouldDie
}

func (c *Cell) sim(iteration int, lastID int) []Cell {
	descendants := []Cell{}

	if c.alive {
		c.eat()
		c.move()
		descendants = c.procreate(iteration, lastID)
		if c.shouldDie(iteration) {
			c.alive = false
			c.diedAt = iteration
		}
	}

	return descendants
}

func getRandomCell(id int) Cell {
	points := int((17 * 100) / (3 + rand.Float64()))

	c := Cell{}

	c.id = id
	c.action = idle
	c.alive = true
	c.bornAt = 0

	c.capacity = rand.Intn(30)
	points -= c.capacity

	c.size = rand.Intn(30) + 1
	points -= c.size

	c.herbivore = int8(rand.Intn(100))
	points -= int(c.herbivore)

	c.division = 1
	points -= int(c.division * 10)

	c.timeToDie = int(rand.NormFloat64()*20) + 50
	points -= c.timeToDie * 10

	c.maxSatiation = int(rand.Intn(100)) + 300
	c.satiation = c.maxSatiation / 2

	c.consumption = 10

	return c
}
