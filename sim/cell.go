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
	return c.cellType.GetMaxSatiation() - c.satiation
}

func (c Cell) shouldEat() bool {
	return c.getLeftToFull() > 0
}

func (c *Cell) consume() int {
	c.satiation -= c.cellType.GetConsumption()
	return c.cellType.GetConsumption()
}

func (c *Cell) eat(food int) int {
	if food > c.getLeftToFull() {
		c.satiation = c.cellType.GetMaxSatiation()
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
	return iteration-c.procreatedAt > int(c.cellType.GetProcreationCd()) && !c.shouldEat()
}

func (c *Cell) shouldProcreate(iteration int, produces []*CellType) bool {
	return c.canProcreate(iteration) && rand.Float32() > .5 && len(produces) > 0
}

func (c *Cell) procreate(iteration int, produces []*CellType) Cell {
	food := c.cellType.GetMaxSatiation() / 2

	ct := produces[rand.Intn(len(produces))]

	descendant := Cell{
		satiation:    food,
		bornAt:       iteration,
		alive:        true,
		cellType:     ct,
		hp:           ct.getMaxHP(),
		procreatedAt: iteration,
	}

	c.satiation = food
	c.bornAt = iteration
	c.procreatedAt = iteration

	return descendant
}

func (c Cell) getAge(iteration int) int8 {
	return int8(iteration - c.bornAt)
}

func (c Cell) shouldDie(
	env Environment,
	iteration int,
	organismPosition r2.Point,
) bool {
	age := c.getAge(iteration)
	isStarving := c.satiation <= 0
	isPastLifetime := c.cellType.GetTimeToDie() < age
	isEnvironmentTooToxic := env.getToxicityOnHeight(c.position.Y+organismPosition.Y) > c.cellType.GetWasteTolerance()

	mustDie := isPastLifetime ||
		isEnvironmentTooToxic ||
		c.hp == 0 ||
		isOutOfBounds(c.position.Add(organismPosition), env)

	dies := mustDie || (isStarving && age > 0)

	return dies
}

func (c *Cell) die(iteration int) {
	c.alive = false
	c.diedAt = iteration
}

// Getters
func (c Cell) GetID() int {
	return c.id
}
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

func (cl CellList) GetCenter() r2.Point {
	boxStart := cl[0].position
	boxEnd := cl[0].position

	for cIndex := range cl {
		if boxStart.X > cl[cIndex].position.X {
			boxStart.X = cl[cIndex].position.X
		}
		if boxStart.Y > cl[cIndex].position.Y {
			boxStart.Y = cl[cIndex].position.Y
		}

		if boxEnd.X < cl[cIndex].position.X {
			boxEnd.X = cl[cIndex].position.X
		}
		if boxEnd.Y < cl[cIndex].position.Y {
			boxEnd.Y = cl[cIndex].position.Y
		}
	}

	center := r2.Point{
		X: (boxEnd.X + boxStart.X) / 2,
		Y: (boxEnd.Y + boxStart.Y) / 2,
	}

	return center
}

func (cl CellList) Remove(removeMap map[int]bool) CellList {
	cells := make(CellList, len(cl))
	index := 0

	for cellIndex := range cl {
		if !removeMap[cellIndex] {
			cells[index] = cl[cellIndex]
			index++
		}
	}

	return cells[:index]
}
