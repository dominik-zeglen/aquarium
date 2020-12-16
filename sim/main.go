package sim

import "fmt"

type Sim struct {
	cells     []Cell
	env       Environment
	iteration int
}

type SpeciesMap map[string]int
type SpeciesData struct {
	Name      string
	Organisms int
}

type WasteData struct {
	Waste        float64
	MinTolerance float64
	MaxTolerance float64
}

type ProcreationData struct {
	CanProcreate bool
	MinCd        int8
	MaxCd        int8
	Species      []SpeciesData
}

type SimData struct {
	CellCount      int
	AliveCellCount int
	Waste          WasteData
	Iteration      int
	Procreation    ProcreationData
}

const maxCells = 1e4

func (s *Sim) RunStep() SimData {
	nextGenCells := []Cell{}
	waste := float64(0)
	species := SpeciesMap{}

	data := SimData{
		CellCount:      len(s.cells),
		AliveCellCount: s.getAliveCells(),
		Iteration:      s.iteration,
		Waste: WasteData{
			MinTolerance: 9999,
			Waste:        s.env.toxicity,
		},
		Procreation: ProcreationData{
			MinCd: 127,
		},
	}

	data.Procreation.CanProcreate = data.AliveCellCount < maxCells

	for cellIndex, cell := range s.cells {
		if cell.alive {
			species[cell.species]++
			if data.Waste.MaxTolerance < cell.wasteTolerance {
				data.Waste.MaxTolerance = cell.wasteTolerance
			}
			if data.Waste.MinTolerance > cell.wasteTolerance {
				data.Waste.MinTolerance = cell.wasteTolerance
			}

			if data.Procreation.MaxCd < cell.procreationCd {
				data.Procreation.MaxCd = cell.procreationCd
			}
			if data.Procreation.MinCd > cell.procreationCd {
				data.Procreation.MinCd = cell.procreationCd
			}
		}

		descendants := s.cells[cellIndex].sim(
			s.env,
			s.iteration,
			s.cells[data.CellCount-1].id,
			data.Procreation.CanProcreate,
		)

		if len(descendants) > 0 {
			nextGenCells = append(nextGenCells, descendants...)
		}

		if !s.cells[cellIndex].alive && s.iteration-s.cells[cellIndex].diedAt > 10 {
			waste += s.cells[cellIndex].getWasteAfterDeath()
		} else {
			if s.cells[cellIndex].alive {
				waste += s.cells[cellIndex].getWaste(s.env)
			}
			nextGenCells = append(nextGenCells, s.cells[cellIndex])
		}

	}

	for name, number := range species {
		data.Procreation.Species = append(data.Procreation.Species, SpeciesData{
			name,
			number,
		})
	}

	s.iteration++
	s.cells = nextGenCells
	s.env.changeToxicity(waste)

	fmt.Printf(
		"Iteration %d, cell count: %d, alive cells: %5d, waste: %.4f, tolerance: %.2f-%.2f, %d species",
		s.iteration,
		len(s.cells),
		s.getAliveCells(),
		s.env.toxicity,
		data.Waste.MinTolerance,
		data.Waste.MaxTolerance,
		len(data.Procreation.Species),
	)
	if data.Procreation.CanProcreate {
		fmt.Printf("\n")
	} else {
		fmt.Print(" x\n")
	}

	return data
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

	for i := 0; i < 10; i++ {
		startCells = append(startCells, getRandomCell(i))
	}

	return Sim{
		cells:     startCells,
		env:       Environment{1},
		iteration: 0,
	}
}
