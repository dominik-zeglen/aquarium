package sim

import "fmt"

type AddSpecies func(species Species) *Species

type Sim struct {
	cells     []Cell
	species   []Species
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
	Species      []Species
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
	s.iteration++

	nextGenCells := []Cell{}
	waste := float64(0)

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
			if data.Waste.MaxTolerance < cell.species.wasteTolerance {
				data.Waste.MaxTolerance = cell.species.wasteTolerance
			}
			if data.Waste.MinTolerance > cell.species.wasteTolerance {
				data.Waste.MinTolerance = cell.species.wasteTolerance
			}

			if data.Procreation.MaxCd < cell.species.procreationCd {
				data.Procreation.MaxCd = cell.species.procreationCd
			}
			if data.Procreation.MinCd > cell.species.procreationCd {
				data.Procreation.MinCd = cell.species.procreationCd
			}
		}

		descendants := s.cells[cellIndex].sim(
			s.env,
			s.iteration,
			s.cells[data.CellCount-1].id,
			s.addSpecies,
			data.Procreation.CanProcreate,
		)

		if len(descendants) > 0 {
			nextGenCells = append(nextGenCells, descendants...)
		}

		if !s.cells[cellIndex].alive && s.iteration-s.cells[cellIndex].diedAt > 10 {
			waste += s.cells[cellIndex].species.getWasteAfterDeath()
		} else {
			if s.cells[cellIndex].alive {
				waste += s.cells[cellIndex].species.getWaste(s.env)
			}
			nextGenCells = append(nextGenCells, s.cells[cellIndex])
		}

	}

	s.cells = nextGenCells
	s.env.changeToxicity(waste)
	s.cleanupSpecies()
	data.Procreation.Species = s.species

	fmt.Printf(
		"Iteration %d, cell count: %d, alive cells: %5d, waste: %.4f, tolerance: %.2f-%.2f, %d species",
		s.iteration,
		len(s.cells),
		s.getAliveCells(),
		s.env.toxicity,
		data.Waste.MinTolerance,
		data.Waste.MaxTolerance,
		s.GetAliveSpecies(),
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

func (s *Sim) addSpecies(species Species) *Species {
	species.ID = len(s.species) + 1
	s.species = append(s.species, species)
	return &s.species[len(s.species)-1]
}

func (s *Sim) removeSpecies(id int) {
	i := 0
	for ; i < len(s.species); i++ {
		if s.species[i].ID == id {
			break
		}
	}
	copy(s.species[i:], s.species[i+1:])
	s.species = s.species[:len(s.species)-1]
}

func (s *Sim) cleanupSpecies() {
	idsToDelete := []int{}
	for speciesIndex, species := range s.species {
		extinct := true
		found := false
		for cellIndex := range s.cells {
			if s.cells[cellIndex].species.ID == species.ID {
				found = true
				if s.cells[cellIndex].alive {
					extinct = false
					break
				}
			}
		}

		if !found {
			idsToDelete = append(idsToDelete, s.species[speciesIndex].ID)
		}
		s.species[speciesIndex].Extinct = extinct
	}

	for _, id := range idsToDelete {
		s.removeSpecies(id)
	}
}

func (s Sim) GetAliveSpecies() int {
	counter := 0
	for speciesIndex := range s.species {
		if !s.species[speciesIndex].Extinct {
			counter++
		}
	}

	return counter
}

func (s Sim) GetCellCount() int {
	return len(s.cells)
}

func (s *Sim) Create() {
	s.iteration = 0
	s.env = Environment{1}

	startCells := []Cell{}

	for i := 0; i < 10; i++ {
		startCells = append(startCells, getRandomCell(i, s.addSpecies))
	}

	s.cells = startCells
}
