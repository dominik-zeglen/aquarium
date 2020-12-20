package sim

import (
	"fmt"
)

type AddSpecies func(species Species) *Species

type Sim struct {
	cells         []Cell
	species       []Species
	env           Environment
	speciesLastID int
	iteration     int
}

type WasteData struct {
	Waste        float64 `json:"waste"`
	MinTolerance float64 `json:"minTolerance"`
	MaxTolerance float64 `json:"maxTolerance"`
}

type ProcreationData struct {
	CanProcreate bool      `json:"canProcreate"`
	MinCd        int8      `json:"minCd"`
	MaxCd        int8      `json:"maxCd"`
	MinHeight    float64   `json:"minHeight"`
	MaxHeight    float64   `json:"maxHeight"`
	Species      []Species `json:"species"`
}

type SimData struct {
	CellCount      int             `json:"cellCount"`
	AliveCellCount int             `json:"aliveCellCount"`
	Waste          WasteData       `json:"waste"`
	Iteration      int             `json:"iteration"`
	Procreation    ProcreationData `json:"procreation"`
}

const maxCells = 3e4

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
			MinCd:     127,
			MinHeight: float64(s.env.height),
		},
	}

	data.Procreation.CanProcreate = data.AliveCellCount < maxCells

	for cellIndex, cell := range s.cells {
		if cell.alive {
			if data.Waste.MaxTolerance < cell.species.WasteTolerance {
				data.Waste.MaxTolerance = cell.species.WasteTolerance
			}
			if data.Waste.MinTolerance > cell.species.WasteTolerance {
				data.Waste.MinTolerance = cell.species.WasteTolerance
			}

			if data.Procreation.MaxCd < cell.species.procreationCd {
				data.Procreation.MaxCd = cell.species.procreationCd
			}
			if data.Procreation.MinCd > cell.species.procreationCd {
				data.Procreation.MinCd = cell.species.procreationCd
			}

			if data.Procreation.MaxHeight < cell.position.Y {
				data.Procreation.MaxHeight = cell.position.Y
			}
			if data.Procreation.MinHeight > cell.position.Y {
				data.Procreation.MinHeight = cell.position.Y
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
				waste += s.cells[cellIndex].species.getWaste(s.env.getToxicityOnHeight(s.cells[cellIndex].position.Y))
			}
			nextGenCells = append(nextGenCells, s.cells[cellIndex])
		}

	}

	s.cells = nextGenCells
	s.env.changeToxicity(waste)
	s.cleanupSpecies()
	data.Procreation.Species = s.species

	fmt.Printf(
		"Iteration %6d, cell count: %5d, alive cells: %5d, waste: %.4f, tolerance: %.2f-%.2f, %3d species",
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
	species.ID = s.speciesLastID
	s.speciesLastID++
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
	specimenCount := make(map[int]int, len(s.species))
	idsToDelete := []int{}

	for cellIndex := range s.cells {
		specimenCount[s.cells[cellIndex].species.ID]++
	}

	for speciesIndex, species := range s.species {
		count, found := specimenCount[species.ID]

		if !found {
			idsToDelete = append(idsToDelete, s.species[speciesIndex].ID)
		}
		s.species[speciesIndex].Extinct = count == 0
		s.species[speciesIndex].Count = count
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

func (s *Sim) KillOldestCells() {
	for i := 2; s.getAliveCells() >= maxCells; i++ {
		for cellIndex := range s.cells {
			if s.cells[cellIndex].species.TimeToDie+s.cells[cellIndex].bornAt-s.iteration < i {
				s.cells[cellIndex].alive = false
			}
		}
	}
}

func (s *Sim) Create() {
	s.iteration = 0
	s.env = Environment{4, 10000, 10000}

	startCells := []Cell{}

	for i := 0; i < 100; i++ {
		startCells = append(startCells, getRandomCell(i, s.env, s.addSpecies))
	}

	s.cells = startCells
}
