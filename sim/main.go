package sim

import "fmt"

type Sim struct {
	cells     []Cell
	env       Environment
	iteration int
}

func (s *Sim) RunStep() {
	nextGenCells := []Cell{}
	waste := float64(0)
	minTolerance := float64(1000)
	maxTolerance := float64(0)

	aliveCells := s.getAliveCells()

	for cellIndex := range s.cells {
		if maxTolerance < s.cells[cellIndex].wasteTolerance {
			maxTolerance = s.cells[cellIndex].wasteTolerance
		}
		if minTolerance > s.cells[cellIndex].wasteTolerance {
			minTolerance = s.cells[cellIndex].wasteTolerance
		}
		descendants := s.cells[cellIndex].sim(s.env, s.iteration, s.cells[len(s.cells)-1].id, aliveCells < 2e5)
		if len(descendants) > 0 {
			nextGenCells = append(nextGenCells, descendants...)
		}

		if !s.cells[cellIndex].alive && s.iteration-s.cells[cellIndex].diedAt > 5 {
			waste += s.cells[cellIndex].getWasteAfterDeath()
		} else {
			if s.cells[cellIndex].alive {
				waste += s.cells[cellIndex].getWaste(s.env)
			}
			nextGenCells = append(nextGenCells, s.cells[cellIndex])
		}

	}

	s.iteration++
	s.cells = nextGenCells
	s.env.changeToxicity(waste)

	fmt.Printf(
		"Iteration %d, cell count: %d, alive cells: %d, waste: %.4f, tolerance: %.4f-%.4f\n",
		s.iteration,
		len(s.cells),
		s.getAliveCells(),
		s.env.toxicity,
		minTolerance,
		maxTolerance,
	)
}

func (s Sim) getAliveCells() int {
	counter := 0

	for cellIndex := range s.cells {
		if s.cells[cellIndex].alive {
			counter++
		}
	}

	return counter
}

func (s Sim) GetCellCount() int {
	return len(s.cells)
}

func CreateSim() Sim {
	startCells := []Cell{}

	for i := 0; i < 6; i++ {
		startCells = append(startCells, getRandomCell(i))
	}

	return Sim{
		cells:     startCells,
		env:       Environment{1},
		iteration: 0,
	}
}
