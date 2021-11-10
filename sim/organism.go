package sim

import (
	"context"
	"math"
	"math/rand"

	"github.com/golang/geo/r2"
	"github.com/opentracing/opentracing-go"
)

const gridSize = 11

type Organism struct {
	id       int
	angle    float64
	position r2.Point
	action   Action
	target   r2.Point

	cells      CellList
	lastCellId int

	speciesID int
	species   *Species

	bornAt int
	diedAt int
}

func (o *Organism) eat(e Environment, iteration int) int {
	food := 0

	// Get produced and stored food
	for cellIndex := range o.cells {
		if o.cells[cellIndex].alive {
			food += o.cells[cellIndex].GetFood(e, iteration, o.position.Y)
			food += o.cells[cellIndex].capacity
			o.cells[cellIndex].capacity = 0
		}
	}

	// Feed cells to compensate their consumption
	for cellIndex := range o.cells {
		if o.cells[cellIndex].alive && food > 0 {
			consumed := o.cells[cellIndex].consume()
			food -= o.cells[cellIndex].eat(consumed)
		}
	}

	// Feed cells that can reproduce
	for cellIndex := range o.cells {
		canReproduce := len(o.species.produces[o.cells[cellIndex].cellType.ID])
		if o.cells[cellIndex].alive && canReproduce > 0 && food > 0 {
			food -= o.cells[cellIndex].eat(food)
		}
	}

	// Store food that had not been eaten
	for cellIndex := range o.cells {
		if o.cells[cellIndex].alive && food > 0 {
			food -= o.cells[cellIndex].storeFood(food)
		}
	}

	return food
}

func (o *Organism) procreate(
	canProcreate bool,
	iteration int,
	maxCells int,
	force bool,
) {
	for cellIndex, cell := range o.cells.GetAlive() {
		if len(o.cells) >= maxCells ||
			!(cell.shouldProcreate(iteration) || force) {
			return
		}

		produced := o.species.produces[cell.cellType.ID]
		producedCt := make([]*CellType, len(produced))

		for ctIndex := range produced {
			producedCt[ctIndex] = &o.species.types[ctIndex]
		}

		if len(producedCt) > 0 {
			freeSpot := getFreeSpot(o.cells, cell, cell.cellType.CanConnect())
			if freeSpot != nil {
				child := o.cells[cellIndex].procreate(iteration, producedCt)
				o.lastCellId++
				child.id = o.lastCellId
				child.position = *freeSpot

				o.cells = append(o.cells, child)
			}
		}

	}
}

func (o Organism) shouldMutate() bool {
	return rand.Float32() > .999
}

func (o *Organism) mutate(addSpecies AddSpecies) {
	newSpecies := o.species.mutate()
	o.species = addSpecies(newSpecies)
	o.speciesID = o.species.id

	for cellIndex := range o.cells {
		ctID := o.cells[cellIndex].cellType.ID

		for cellTypeIndex := range o.species.types {
			if o.species.types[cellTypeIndex].ID == ctID {
				o.cells[cellIndex].cellType = &o.species.types[cellTypeIndex]
			}
		}
	}
}

func (o Organism) IsAlive() bool {
	return o.cells.GetAliveCount() > 0
}

func (o *Organism) die(iteration int) {
	for cellIndex := range o.cells {
		o.cells[cellIndex].die(iteration)
	}
}

func (o *Organism) move() r2.Point {
	var moveVec r2.Point

	if o.action == idle {
		if rand.Float32() > .9 {
			o.angle = getRandomAngle()
		}
		moveVec = getVecFromAngle(o.angle)
	} else {
		moveVec = o.
			target.
			Sub(o.position).
			Normalize()
	}

	scaledMoveVec := moveVec.Mul(
		float64(o.GetMobility()*10) / float64(o.GetMass()),
	)
	o.position = o.position.Add(scaledMoveVec)

	return scaledMoveVec
}

func (o *Organism) killCells(env Environment, iteration int) {
	for cellIndex := range o.cells {
		cell := &o.cells[cellIndex]
		if cell.shouldDie(env, iteration, o.position) {
			cell.die(iteration)
		}
	}
}

func (o *Organism) split(
	ctx context.Context,
	canProcreate bool,
	iteration int,
) []Organism {
	if len(o.cells) == 1 {
		return OrganismList{}
	}

	splitSpan, _ := opentracing.StartSpanFromContext(
		ctx,
		"split-grids",
	)
	defer splitSpan.Finish()

	grid := make([][]*int, gridSize)
	for i := 0; i < gridSize; i++ {
		grid[i] = make([]*int, gridSize)
	}

	for cellIndex := range o.cells {
		x := int(o.cells[cellIndex].position.X) + gridSize/2
		y := int(o.cells[cellIndex].position.Y) + gridSize/2

		grid[y][x] = &cellIndex
	}

	graph := CellGraph{}
	graph.Init(gridSize, gridSize, grid)
	ptrGrids := graph.GetIslands()
	grids := make([]CellList, len(ptrGrids))

	for i := 0; i < len(grids); i++ {
		grids[i] = make(CellList, len(ptrGrids[i]))
		for j := 0; j < len(ptrGrids[i]); j++ {
			grids[i][j] = o.cells[*ptrGrids[i][j]]
		}
	}

	if len(grids) > 1 {
		organisms := make(OrganismList, len(grids)-1)

		maxCellsIndex := 0
		for gridIndex := range grids {
			if len(grids[maxCellsIndex]) < len(grids[gridIndex]) {
				maxCellsIndex = gridIndex
			}
		}

		o.cells = grids[maxCellsIndex]
		gridsSliced := append(grids[:maxCellsIndex], grids[maxCellsIndex+1:]...)

		if !canProcreate {
			return OrganismList{}
		}

		for gridIndex, grid := range gridsSliced {
			center := grid.GetCenter()
			grid = grid.Center()

			for cIndex := range grid {
				if !canProcreate {
					grid[cIndex].alive = false
					grid[cIndex].diedAt = iteration
				}
			}

			organisms[gridIndex] = *o
			organisms[gridIndex].cells = grid
			organisms[gridIndex].position = o.position.Add(center)
			organisms[gridIndex].angle = rand.Float64() * 2 * math.Pi
			organisms[gridIndex].lastCellId = len(grid) - 1
		}

		return organisms
	}

	return OrganismList{}
}

func (o *Organism) sim(
	ctx context.Context,
	env Environment,
	iteration int,
	maxCells int,
	addSpecies AddSpecies,
	canProcreate bool,
) OrganismList {
	age := iteration - o.bornAt
	if o.IsAlive() {
		o.eat(env, iteration)
		o.move()

		if o.shouldMutate() {
			o.mutate(addSpecies)
		}

		if age < 3 || age > 200+iteration/3200 {
			if rand.Float64() > .66 {
				o.die(iteration)
			}
		} else {
			o.procreate(canProcreate, iteration, maxCells, false)
			o.killCells(env, iteration)
		}

		if o.cells.GetAliveCount() == 0 {
			o.diedAt = iteration
			return []Organism{}
		}

		return o.split(ctx, canProcreate, iteration)
	}

	return OrganismList{}
}

// Getters

func (o Organism) GetMass() int {
	mass := 0
	for cellIndex := range o.cells {
		mass += o.cells[cellIndex].cellType.getMass()
	}

	return mass
}
func (o Organism) GetMobility() int {
	mobility := 0
	for cellIndex := range o.cells {
		mobility += o.cells[cellIndex].cellType.GetMobility()
	}

	return mobility
}
func (o Organism) GetPosition() r2.Point {
	return o.position
}
func (o Organism) GetBornAt() int {
	return o.bornAt
}
func (o Organism) GetCells() CellList {
	return o.cells
}
func (o Organism) GetID() int {
	return o.id
}
func (o Organism) GetSpecies() Species {
	return *o.species
}

func getRandomOrganism(id int, e Environment, addSpecies AddSpecies) Organism {
	s := addSpecies(getRandomHerbivore())
	ct := &s.types[0]

	c := Cell{
		alive:     true,
		cellType:  ct,
		satiation: 20,
		hp:        ct.getMaxHP(),
	}

	return Organism{
		id:     id,
		angle:  rand.Float64() * 2 * math.Pi,
		cells:  CellList{c},
		action: idle,
		position: r2.Point{
			X: float64(e.width)*rand.Float64()*.8 + float64(e.width)/10,
			Y: float64(e.height)*rand.Float64()*.8 + float64(e.height)/10,
		},
		species:   s,
		speciesID: s.id,
	}

}

type OrganismList []Organism

func (ol OrganismList) GetArea(start r2.Point, end r2.Point) OrganismList {
	organisms := make(OrganismList, len(ol))

	index := 0
	for organismIndex, organism := range ol {
		position := organism.GetPosition()
		if position.X > start.X && position.X < end.X &&
			position.Y > start.Y && position.Y < end.Y {
			organisms[index] = ol[organismIndex]
			index++
		}
	}

	return organisms[:index]
}

// GetAreaCount serves as an optimisation
func (ol OrganismList) GetAreaCount(start r2.Point, end r2.Point) int {
	counter := 0

	for _, organism := range ol {
		position := organism.position
		if position.X > start.X && position.X < end.X &&
			position.Y > start.Y && position.Y < end.Y {
			counter++
		}
	}

	return counter
}

func (ol OrganismList) GetAlive() OrganismList {
	organisms := make(OrganismList, len(ol))

	index := 0
	for organismIndex, organism := range ol {
		if organism.IsAlive() {
			organisms[index] = ol[organismIndex]
			index++
		}
	}

	return organisms[:index]
}

// GetAliveCount serves as an optimisation
func (ol OrganismList) GetAliveCount() int {
	counter := 0
	for _, organism := range ol {
		if organism.IsAlive() {
			counter++
		}
	}

	return counter
}

func (ol OrganismList) GetSpecies(id int) OrganismList {
	organisms := make(OrganismList, len(ol))

	index := 0
	for organismIndex, organism := range ol {
		if organism.species.id == id {
			organisms[index] = ol[organismIndex]
			index++
		}
	}

	return organisms[:index]
}
