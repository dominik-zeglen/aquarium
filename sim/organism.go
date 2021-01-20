package sim

import (
	"math/rand"

	"github.com/golang/geo/r2"
)

type Organism struct {
	id       int
	angle    float32
	position r2.Point
	action   Action
	target   r2.Point

	cells CellList

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
	force bool,
) {
	if canProcreate {
		for cellIndex := range o.cells {
			produced := o.species.produces[o.cells[cellIndex].cellType.ID]
			producedCt := make([]*CellType, len(produced))

			for ctIndex := range produced {
				producedCt[ctIndex] = &o.species.types[ctIndex]
			}

			if o.cells[cellIndex].shouldProcreate(iteration, producedCt) || force {
				child := o.cells[cellIndex].procreate(iteration, producedCt)
				child.id = o.cells[len(o.cells)-1].id + 1
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

func (o *Organism) move() r2.Point {
	var moveVec r2.Point

	if o.action == idle {
		moveVec = getRandomVec()
	} else {
		moveVec = o.
			target.
			Sub(o.position).
			Normalize()
	}

	scaledMoveVec := moveVec.Mul(float64(o.GetMobility()*10) / float64(o.GetMass()))
	o.position = o.position.Add(scaledMoveVec)

	return scaledMoveVec
}

func (o *Organism) killCells(env Environment, iteration int) {
	for cellIndex := range o.cells {
		if o.cells[cellIndex].shouldDie(env, iteration, o.position) {
			o.cells[cellIndex].die(iteration)
		}
	}
}

func (o *Organism) split() []Organism {
	grids := []CellList{}
	for _, cell := range o.cells {
		found := false
		for gridIndex, grid := range grids {
			for _, cellInGrid := range grid {
				if cell.position.Sub(cellInGrid.position).Norm() == 1 {
					grids[gridIndex] = append(grids[gridIndex], cell)
					found = true
					break
				}
			}
		}

		if !found {
			grids = append(grids, CellList{cell})
		}
	}

	lastLen := len(grids)
	do := true

	for (lastLen < len(grids) || do) && len(grids) > 1 {
		lastLen = len(grids)
		do = false

		for gridAIndex, gridA := range grids {

			for gridBIndex, gridB := range grids {
				found := false
				if gridAIndex == gridBIndex {
					break
				}

				for _, cellA := range gridA {
					if found {
						break
					}
					for _, cellB := range gridB {
						if cellA.id == cellB.id {
							found = true
							break
						}
					}
				}

				if found {
					grids[gridAIndex] = append(grids[gridAIndex], grids[gridBIndex]...)
					grids = append(grids[:gridBIndex], grids[:gridBIndex+1]...)
				}
			}
		}
	}

	if len(grids) > 1 {
		organisms := make(OrganismList, len(grids)-1)
		o.cells = grids[0]

		for gridIndex, grid := range grids[1:] {
			center := grid.GetCenter()

			for cIndex := range grid {
				grid[cIndex].id = cIndex
				grid[cIndex].position = grid[cIndex].position.Sub(center)
			}

			organisms[gridIndex] = *o
			organisms[gridIndex].cells = grid
			organisms[gridIndex].position = o.position.Add(center)
		}

		return organisms
	}

	return OrganismList{}
}

func (o *Organism) sim(
	env Environment,
	iteration int,
	addSpecies AddSpecies,
	canProcreate bool,
) OrganismList {
	if o.IsAlive() {
		o.eat(env, iteration)
		o.move()

		if o.shouldMutate() {
			o.mutate(addSpecies)
		}

		o.procreate(canProcreate, iteration, false)
		o.killCells(env, iteration)

		if o.cells.GetAliveCount() == 0 {
			o.diedAt = iteration
		}

		return o.split()
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
		mobility += o.cells[cellIndex].cellType.mobility
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
		cells:  CellList{c},
		action: idle,
		position: r2.Point{
			X: float64(e.width)*rand.Float64()*.8 + float64(e.width)/10,
			Y: float64(e.height)*rand.Float64()*.8 + float64(e.height)/10,
		},
		species: s,
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
