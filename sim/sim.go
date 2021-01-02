package sim

import (
	"fmt"
	"sync"
)

type Sim struct {
	cells         CellList
	species       SpeciesList
	env           Environment
	speciesLastID int
	iteration     int
	maxCells      int
	verbose       bool
	lock          sync.Mutex
}

func (d *IterationData) from(from IterationData) {
	d.AliveCellCount = from.AliveCellCount
	d.CellCount = from.CellCount
	d.Iteration = from.Iteration
	d.Procreation = from.Procreation
	d.Waste = from.Waste
}

func (s Sim) GetEnvironment() Environment {
	return s.env
}

func (s Sim) GetIteration() int {
	return s.iteration
}

func (s Sim) GetAliveCellCount() int {
	return s.cells.GetAliveCount()
}

func (s Sim) GetCellCount() int {
	return len(s.cells)
}

func (s *Sim) GetCells() CellList {
	return s.cells
}

func (s Sim) GetSpecies() SpeciesList {
	return s.species
}

func (s *Sim) Lock() {
	s.lock.Lock()
}

func (s *Sim) Unlock() {
	s.lock.Unlock()
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

func (s *Sim) KillOldestCells() {
	for i := 1; s.GetAliveCellCount() >= s.maxCells; i++ {
		for cellIndex := range s.cells {
			if s.cells[cellIndex].species.TimeToDie+s.cells[cellIndex].bornAt-s.iteration < i {
				s.cells[cellIndex].alive = false
			}
		}
	}
}

func (s *Sim) Create(verbose bool) {
	s.iteration = 0
	s.env = Environment{4, 10000, 10000}

	startCells := make(CellList, 100)

	for i := 0; i < 100; i++ {
		startCells[i] = getRandomCell(i, s.env, s.addSpecies)
	}

	s.cells = startCells
	s.maxCells = 2e4
	s.verbose = verbose
}

func (s *Sim) RunStep() IterationData {
	s.iteration++

	nextGenCells := make(CellList, 2*s.maxCells)
	waste := float64(0)

	data := IterationData{
		CellCount:      len(s.cells),
		AliveCellCount: s.GetAliveCellCount(),
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

	data.Procreation.CanProcreate = data.AliveCellCount < s.maxCells
	index := 0

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

		descendant := s.cells[cellIndex].sim(
			s.env,
			s.iteration,
			nextGenCells[data.CellCount-1].id,
			s.addSpecies,
			data.Procreation.CanProcreate,
		)

		if descendant != nil {
			nextGenCells[index] = *descendant
			index++
		}

		if !s.cells[cellIndex].alive && s.iteration-s.cells[cellIndex].diedAt > 10 {
			waste += s.cells[cellIndex].species.getWasteAfterDeath()
		} else {
			if s.cells[cellIndex].alive {
				waste += s.cells[cellIndex].species.getWaste(s.env.getToxicityOnHeight(s.cells[cellIndex].position.Y))
			}
			nextGenCells[index] = s.cells[cellIndex]
			index++
		}

	}

	s.cells = nextGenCells[:index]
	s.env.changeToxicity(waste)
	s.cleanupSpecies()
	data.Procreation.Species = s.species

	if s.verbose || s.GetAliveCellCount() == 0 {
		fmt.Printf(
			"Iteration %6d, cell count: %5d, alive cells: %5d, waste: %.4f, tolerance: %.2f-%.2f, %3d species",
			s.iteration,
			len(s.cells),
			s.GetAliveCellCount(),
			s.env.toxicity,
			data.Waste.MinTolerance,
			data.Waste.MaxTolerance,
			len(s.GetSpecies().GetAlive()),
		)
		if data.Procreation.CanProcreate {
			fmt.Printf("\n")
		} else {
			fmt.Print(" x\n")
		}
	}

	return data
}

func (s *Sim) RunLoop(data *IterationData) {
	consecutiveNoProcreateIterations := 0

	for {
		s.lock.Lock()
		iterationData := s.RunStep()
		data.from(iterationData)

		if !iterationData.Procreation.CanProcreate {
			consecutiveNoProcreateIterations++
		} else {
			consecutiveNoProcreateIterations = 0
		}

		if consecutiveNoProcreateIterations > 2 {
			s.KillOldestCells()
		}
		s.lock.Unlock()

		if s.GetCellCount() == 0 {
			break
		}

		// if iterationData.Iteration == 1 {
		// 	break
		// }

		// if true {
		// 	time.Sleep(time.Second * 2)
		// }

	}
}
