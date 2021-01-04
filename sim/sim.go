package sim

import (
	"fmt"
	"sync"
)

type Sim struct {
	organisms      OrganismList
	species        SpeciesList
	env            Environment
	speciesLastID  int
	organismLastID int
	iteration      int
	maxCells       int
	verbose        bool
	lock           sync.Mutex
}

func (d *IterationData) from(from IterationData) {
	d.AliveCellCount = from.AliveCellCount
	d.CellCount = from.CellCount
	d.Iteration = from.Iteration
	d.Procreation = from.Procreation
	d.Waste = from.Waste
}

func (s *Sim) GetNewOrganismID() int {
	s.organismLastID++
	return s.organismLastID
}

func (s Sim) GetEnvironment() Environment {
	return s.env
}

func (s Sim) GetIteration() int {
	return s.iteration
}

func (s Sim) GetAliveCount() int {
	return s.organisms.GetAliveCount()
}

func (s Sim) GetAliveCellCount() int {
	counter := 0
	alive := s.organisms.GetAlive()

	for oIndex := range alive {
		counter += alive[oIndex].cells.GetAliveCount()
	}
	return counter
}

func (s Sim) GetCellCount() int {
	return len(s.organisms)
}

func (s *Sim) GetOrganisms() OrganismList {
	return s.organisms
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
	species.EmergedAt = s.iteration
	species.Count = 1
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

	for oIndex := range s.organisms {
		specimenCount[s.organisms[oIndex].species.ID]++
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
	for i := int8(1); s.GetAliveCellCount() >= s.maxCells; i++ {
		for organismIndex := range s.organisms {
			for cellIndex := range s.organisms[organismIndex].cells {
				age := int8(s.iteration - s.organisms[organismIndex].cells[cellIndex].bornAt)
				if s.organisms[organismIndex].cells[cellIndex].cellType.TimeToDie-age < i {
					s.organisms[organismIndex].cells[cellIndex].alive = false
				}
			}
		}
	}
}

func (s *Sim) Create(verbose bool) {
	s.iteration = 0
	s.env = Environment{4, 10000, 10000}

	startCells := make(OrganismList, 100)

	for i := 0; i < 100; i++ {
		startCells[i] = getRandomOrganism(i, s.env, s.addSpecies)
	}

	s.organisms = startCells
	s.maxCells = 2e4
	s.verbose = verbose
}

func (s *Sim) RunStep() IterationData {
	s.iteration++

	nextGenOrganisms := make(OrganismList, s.maxCells*5)
	waste := float64(0)

	data := IterationData{
		CellCount:      len(s.organisms),
		AliveCellCount: s.GetAliveCellCount(),
		Iteration:      s.iteration,
		Waste: WasteData{
			MinTolerance: s.species[0].Types[0].WasteTolerance,
			Waste:        s.env.toxicity,
		},
		Procreation: ProcreationData{
			MinCd:     127,
			MinHeight: float64(s.env.height),
		},
	}

	data.Procreation.CanProcreate = data.AliveCellCount < s.maxCells
	index := 0

	for _, species := range s.species {
		if !species.Extinct {
			for tIndex := range species.Types {
				if data.Waste.MaxTolerance < species.Types[tIndex].WasteTolerance {
					data.Waste.MaxTolerance = species.Types[tIndex].WasteTolerance
				}
				if data.Waste.MinTolerance > species.Types[tIndex].WasteTolerance {
					data.Waste.MinTolerance = species.Types[tIndex].WasteTolerance
				}
			}
		}
	}

	for organismIndex, organism := range s.organisms {
		if organism.IsAlive() {
			if data.Procreation.MaxHeight < organism.position.Y {
				data.Procreation.MaxHeight = organism.position.Y
			}
			if data.Procreation.MinHeight > organism.position.Y {
				data.Procreation.MinHeight = organism.position.Y
			}
		}

		descendants := s.organisms[organismIndex].sim(
			s.env,
			s.iteration,
			s.addSpecies,
			data.Procreation.CanProcreate,
		)

		for dIndex := range descendants {
			descendants[dIndex].id = s.GetNewOrganismID()
			descendants[dIndex].bornAt = s.iteration
			nextGenOrganisms[index] = descendants[dIndex]
			index++
		}

		for _, cell := range s.organisms[organismIndex].cells {
			if !cell.alive && s.iteration-cell.diedAt > 10 {
				waste += cell.cellType.getWasteAfterDeath()
			} else {
				if cell.alive {
					waste += cell.cellType.getWaste(s.env.getToxicityOnHeight(cell.position.Y))
				}
				nextGenOrganisms[index] = s.organisms[organismIndex]
				index++
			}
		}
	}

	s.organisms = nextGenOrganisms[:index]
	s.env.changeToxicity(waste)
	s.cleanupSpecies()
	data.Procreation.Species = s.species

	if s.iteration == 1000 {
		fmt.Println("oh")
	}

	if s.verbose || s.GetAliveCellCount() == 0 {
		fmt.Printf(
			"Iteration %6d, organisms: %5d, cells: %5d, waste: %.4f %d species",
			s.iteration,
			s.GetAliveCount(),
			s.GetAliveCellCount(),
			s.env.toxicity,
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
