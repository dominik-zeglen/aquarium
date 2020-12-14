package sim

import "fmt"

type Sim struct {
	cells     []Cell
	iteration int
}

func (s *Sim) RunStep() {
	nextGenCells := []Cell{}

	for cellIndex := range s.cells {
		descendants := s.cells[cellIndex].sim(s.iteration, s.cells[len(s.cells)-1].id)
		if len(descendants) > 0 {
			nextGenCells = append(nextGenCells, descendants...)
		}

		if s.cells[cellIndex].diedAt-s.iteration < 10 || s.cells[cellIndex].alive {
			nextGenCells = append(nextGenCells, s.cells[cellIndex])
		}
	}

	s.iteration++
	s.cells = nextGenCells

	fmt.Printf(
		"Iteration %d, cell count: %d, alive cells: %d\n",
		s.iteration,
		len(s.cells),
		s.getAliveCells(),
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

func CreateSim() Sim {
	startCells := []Cell{}

	for i := 0; i < 10; i++ {
		startCells = append(startCells, getRandomCell(i))
	}

	return Sim{
		cells:     startCells,
		iteration: 0,
	}
}
