package sim

import (
	"math/rand"

	"github.com/golang/geo/r2"
)

type Cell struct {
	id       int
	position r2.Point

	cellType *CellType

	alive        bool
	hp           int
	bornAt       int
	diedAt       int
	procreatedAt int

	satiation int
	capacity  int
}

func (c Cell) GetFood(e Environment, iteration int, organismHeight float64) int {
	food := 0
	if c.cellType.Herbivore > 0 {
		food += int(float64(c.cellType.Herbivore) * e.getLightOnHeight(c.position.Y+organismHeight, iteration))
	}
	if c.cellType.Funghi > 0 {
		food += int(c.cellType.getProcessedWaste(e.getToxicityOnHeight(c.position.Y + organismHeight)))
	}

	return food
}

func (c Cell) getLeftToFull() int {
	return c.cellType.maxSatiation - c.satiation
}

func (c Cell) shouldEat() bool {
	return c.cellType.maxSatiation > c.satiation
}

func (c *Cell) consume() int {
	c.satiation -= c.cellType.GetConsumption()
	return c.cellType.GetConsumption()
}

func (c *Cell) eat(food int) int {
	if food > c.getLeftToFull() {
		c.satiation = c.cellType.maxSatiation
		leftover := food - c.getLeftToFull()

		return food - leftover

	}

	c.satiation += food
	return food
}

func (c *Cell) storeFood(food int) int {
	toStore := c.cellType.maxCapacity - c.capacity

	if food > toStore {
		c.capacity = c.cellType.maxCapacity
		return toStore
	} else {
		c.capacity += food
		return food
	}
}

func (c Cell) canProcreate(iteration int) bool {
	if c.procreatedAt == 0 {
		return !c.shouldEat()
	}
	return iteration-c.procreatedAt > int(c.cellType.procreationCd) && !c.shouldEat()
}

func (c *Cell) shouldProcreate(iteration int, produces []*CellType) bool {
	return c.canProcreate(iteration) && rand.Float32() > .5 && len(produces) > 0
}

func (c *Cell) procreate(iteration int, produces []*CellType) Cell {
	food := c.cellType.maxSatiation / 2
	vec := getRandomVec().Mul(float64(c.cellType.size))
	ct := produces[rand.Intn(len(produces))]

	descendant := Cell{
		satiation: food,
		bornAt:    iteration,
		alive:     true,
		position:  c.position.Add(vec),
		cellType:  ct,
		hp:        ct.getMaxHP(),
	}

	c.satiation = food
	c.bornAt = iteration
	c.procreatedAt = iteration

	return descendant
}

func (c Cell) shouldDie(
	env Environment,
	iteration int,
	organismPosition r2.Point,
) bool {
	age := int8(iteration - c.bornAt)
	isStarving := c.satiation <= 0
	isPastLifetime := c.cellType.TimeToDie < age
	isEnvironmentTooToxic := env.getToxicityOnHeight(c.position.Y+organismPosition.Y) > c.cellType.WasteTolerance

	mustDie := isPastLifetime ||
		isEnvironmentTooToxic ||
		c.hp == 0 ||
		isOutOfBounds(c.position.Add(organismPosition), env)

	if mustDie {
		return true
	}

	if isStarving {
		if age > 0 {
			return true
		}
	}

	return false
}

func (c *Cell) die(iteration int) {
	c.alive = false
	c.diedAt = iteration
}

// Getters
func (c Cell) GetPosition() r2.Point {
	return c.position
}
func (c Cell) GetType() *CellType {
	return c.cellType
}
func (c Cell) IsAlive() bool {
	return c.alive
}
func (c Cell) GetHP() int {
	return c.hp
}
func (c Cell) GetBornAt() int {
	return c.bornAt
}
func (c Cell) GetSatiation() int {
	return c.satiation
}
func (c Cell) GetCapacity() int {
	return c.capacity
}

type CellList []Cell

func (cl CellList) GetArea(start r2.Point, end r2.Point) CellList {
	cells := make(CellList, len(cl))

	index := 0
	for cellIndex, cell := range cl {
		position := cell.GetPosition()
		if position.X > start.X && position.X < end.X &&
			position.Y > start.Y && position.Y < end.Y {
			cells[index] = cl[cellIndex]
			index++
		}
	}

	return cells[:index]
}

func (cl CellList) GetAlive() CellList {
	cells := make(CellList, len(cl))

	index := 0
	for cellIndex, cell := range cl {
		if cell.alive {
			cells[index] = cl[cellIndex]
			index++
		}
	}

	return cells[:index]
}

// GetAliveCount serves as an optimisation
func (cl CellList) GetAliveCount() int {
	counter := 0
	for _, cell := range cl {
		if cell.alive {
			counter++
		}
	}

	return counter
}
