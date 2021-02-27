package sim

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/geo/r2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SimConfig struct {
	EnvDivisions       int
	MaxCellsInOrganism int
	MaxOrganisms       int
	StartCells         int
	Verbose            bool
	WarmupIterations   int
}

type Sim struct {
	areaCount          int
	env                Environment
	iteration          int
	lock               sync.Mutex
	maxCells           int
	maxCellsInOrganism int
	organismLastID     int
	organisms          OrganismList
	species            SpeciesList
	speciesLastID      int
	speciesLock        sync.Mutex
	verbose            bool
	warmupIterations   int
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
	s.speciesLock.Lock()

	sp := species.copy()
	sp.id = s.speciesLastID
	sp.emergedAt = s.iteration
	sp.count = 1
	s.speciesLastID++
	s.species = append(s.species, sp)
	s.speciesLock.Unlock()
	return &s.species[len(s.species)-1]
}

func (s *Sim) removeSpecies(id int) {
	i := 0
	for ; i < len(s.species); i++ {
		if s.species[i].id == id {
			break
		}
	}

	copy(s.species[i:], s.species[i+1:])

	s.species = s.species[:len(s.species)-1]
}

func (s *Sim) cleanupSpecies(ctx context.Context) {
	span, spanCtx := opentracing.StartSpanFromContext(
		ctx,
		"cleanup-species",
	)
	defer span.Finish()

	specimenCount := make(map[int]int, len(s.species))
	idsToDelete := []int{}

	for oIndex := range s.organisms {
		specimenCount[s.organisms[oIndex].species.id]++
	}

	for speciesIndex, species := range s.species {
		count, found := specimenCount[species.id]

		if !found || count == 0 {
			idsToDelete = append(idsToDelete, s.species[speciesIndex].id)
		}
		s.species[speciesIndex].extinct = false
		s.species[speciesIndex].count = count
	}

	for _, id := range idsToDelete {
		s.removeSpecies(id)
	}

	reindexSpan, _ := opentracing.StartSpanFromContext(
		spanCtx,
		"reindex",
	)
	speciesMap := make(map[int]*Species, len(s.species))
	for speciesIndex := range s.species {
		speciesMap[s.species[speciesIndex].id] = &s.species[speciesIndex]
	}
	for organismIndex, organism := range s.organisms {
		s.organisms[organismIndex].species = speciesMap[organism.speciesID]
	}
	reindexSpan.Finish()
}

func (s Sim) getAreas(ctx context.Context) []bool {
	span, spanCtx := opentracing.StartSpanFromContext(
		ctx,
		"get-areas",
	)
	defer span.Finish()

	areas := make([]bool, s.areaCount*s.areaCount+(s.areaCount-1)*(s.areaCount-1))
	aliveOrganisms := s.organisms.GetAlive()

	for areaIndex := range areas {
		auxArea := areaIndex >= s.areaCount*s.areaCount
		max := s.areaCount

		if auxArea {
			max--
		}

		start := r2.Point{
			X: float64(s.env.width / s.areaCount * (areaIndex % s.areaCount)),
			Y: float64(s.env.height / s.areaCount * ((areaIndex / s.areaCount) % s.areaCount)),
		}
		end := start.Add(r2.Point{
			X: float64(s.env.width / s.areaCount),
			Y: float64(s.env.height / s.areaCount),
		})

		if auxArea {
			offset := r2.Point{
				X: float64(s.env.width / s.areaCount / 2),
				Y: float64(s.env.height / s.areaCount / 2),
			}

			start = start.Add(offset)
			end = end.Add(offset)
		}

		countSpan, _ := opentracing.StartSpanFromContext(
			spanCtx,
			"count-area",
		)
		organisms := aliveOrganisms.GetAreaCount(start, end)
		areas[areaIndex] = organisms < (s.maxCells / s.areaCount / s.areaCount)
		countSpan.Finish()
	}

	return areas
}

func (s *Sim) Create(config SimConfig) {
	s.iteration = 0
	s.env = Environment{4, 10000, 10000}

	startCells := make(OrganismList, config.StartCells)

	for i := 0; i < config.StartCells; i++ {
		startCells[i] = getRandomOrganism(i, s.env, s.addSpecies)
	}

	s.organisms = startCells
	s.maxCells = config.MaxOrganisms
	s.areaCount = config.EnvDivisions
	s.verbose = config.Verbose
	s.warmupIterations = config.WarmupIterations
	s.maxCellsInOrganism = config.MaxCellsInOrganism
}

func (s *Sim) RunStep(ctx context.Context) IterationData {
	span, stepSpanCtx := opentracing.StartSpanFromContext(
		ctx,
		"step",
	)
	span.LogFields(
		log.Int("Iteration", s.iteration),
		log.Int("Organisms", len(s.organisms)),
		log.Int("Species", len(s.species)),
	)
	defer span.Finish()

	s.iteration++

	nextGenOrganisms := make(OrganismList, s.maxCells*5)
	waste := float64(0)

	dataSpan, _ := opentracing.StartSpanFromContext(stepSpanCtx, "get-data")
	data := IterationData{
		CellCount: len(s.organisms),
		Iteration: s.iteration,
		Waste: WasteData{
			MinTolerance: s.species[0].types[0].GetWasteTolerance(),
			Waste:        s.env.toxicity,
		},
		Procreation: ProcreationData{
			MinCd:     s.species[0].types[0].GetProcreationCd(),
			MinHeight: float64(s.env.height),
		},
	}
	dataSpan.Finish()

	data.Procreation.CanProcreate = data.AliveCellCount < s.maxCells
	index := 0
	highestPoints := 0
	allPoints := 0
	highestConnectivity := 1
	allConnectivity := 0
	allConnecting := 0

	for _, species := range s.species {
		if !species.extinct {
			for tIndex := range species.types {
				allPoints += species.points
				allConnectivity += int(species.types[0].connects)
				if data.Waste.MaxTolerance < species.types[tIndex].GetWasteTolerance() {
					data.Waste.MaxTolerance = species.types[tIndex].GetWasteTolerance()
				}
				if data.Waste.MinTolerance > species.types[tIndex].GetWasteTolerance() {
					data.Waste.MinTolerance = species.types[tIndex].GetWasteTolerance()
				}
				if highestPoints < species.points {
					highestPoints = species.points
				}
				if highestConnectivity < int(species.types[0].connects) {
					highestConnectivity = int(species.types[0].connects)
				}
				if species.types[0].CanConnect() {
					allConnecting++
				}
			}
		}
	}

	avgPoints := allPoints / len(s.species)
	avgConnectivity := allConnectivity / len(s.species)

	areas := s.getAreas(stepSpanCtx)

	simSpan, simSpanCtx := opentracing.StartSpanFromContext(stepSpanCtx, "sim")
	removedCellCounter := 0
	for organismIndex, organism := range s.organisms {
		if organism.IsAlive() {
			if data.Procreation.MaxHeight < organism.position.Y {
				data.Procreation.MaxHeight = organism.position.Y
			}
			if data.Procreation.MinHeight > organism.position.Y {
				data.Procreation.MinHeight = organism.position.Y
			}
		}

		pos := fitToBoundary(s.organisms[organismIndex].position, s.env)
		mainArea := int(pos.Y)*s.areaCount/s.env.height*s.areaCount + int(pos.X)*s.areaCount/s.env.width
		canProcreate := areas[mainArea]
		xAux := (int(pos.X) % (s.env.width / s.areaCount)) > s.env.width/s.areaCount/2
		yAux := (int(pos.Y) % (s.env.height / s.areaCount)) > s.env.height/s.areaCount/2
		if xAux && yAux {
			xo := int(pos.X) - s.env.width/s.areaCount/2
			yo := int(pos.Y) - s.env.height/s.areaCount/2

			auxArea := yo*s.areaCount/s.env.height*(s.areaCount-1) + xo*s.areaCount/s.env.width
			canProcreate = canProcreate && areas[auxArea]
		}

		descendants := s.organisms[organismIndex].sim(
			simSpanCtx,
			s.env,
			s.iteration,
			s.maxCellsInOrganism,
			s.addSpecies,
			canProcreate,
		)

		for dIndex := range descendants {
			descendants[dIndex].id = s.GetNewOrganismID()
			descendants[dIndex].bornAt = s.iteration
			nextGenOrganisms[index] = descendants[dIndex]
			index++
		}

		alive := false
		removeMap := map[int]bool{}
		for cellIndex, cell := range organism.cells {
			if !cell.alive && s.iteration-cell.diedAt > 5 {
				waste += cell.cellType.getWasteAfterDeath()
				removeMap[cellIndex] = true
			} else {
				if cell.alive {
					alive = true
					waste += cell.cellType.getWaste(s.env.getToxicityOnHeight(cell.position.Y))
					data.AliveCellCount++
				}
			}
		}

		if len(removeMap) > 0 {
			removedCellCounter += len(removeMap)
			s.organisms[organismIndex].cells = organism.cells.Remove(removeMap)
		}

		if alive {
			nextGenOrganisms[index] = s.organisms[organismIndex]
			index++
		}
	}
	simSpan.LogFields(
		log.Int("removed cells", removedCellCounter),
	)
	simSpan.Finish()

	s.organisms = nextGenOrganisms[:index]
	s.env.changeToxicity(waste)

	s.cleanupSpecies(stepSpanCtx)

	data.Procreation.Species = s.species

	if s.verbose {
		fmt.Printf(
			"It: %6d, o: %5d, w: %.4f sp: %4d, HLvl: %3d, AvgLvl: %3d, HCn: %2d, AvgCn: %2d, CCn: %2d\n",
			s.iteration,
			len(s.organisms),
			s.env.toxicity,
			len(s.GetSpecies().GetAlive()),
			highestPoints-startingPoints+1,
			avgPoints-startingPoints+1,
			highestConnectivity,
			avgConnectivity,
			allConnecting,
		)
	}

	return data
}

func (s *Sim) RunLoop(data *IterationData) {
	for {
		spanName := fmt.Sprintf("loop %d", data.Iteration+1)
		span := opentracing.GlobalTracer().StartSpan(spanName)
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		s.lock.Lock()
		iterationData := s.RunStep(ctx)
		data.from(iterationData)

		s.lock.Unlock()

		if s.GetCellCount() == 0 {
			break
		}

		span.Finish()

		if iterationData.Iteration > s.warmupIterations {
			time.Sleep(time.Second)
		}
	}
}
