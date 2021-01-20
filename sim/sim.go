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

type Sim struct {
	organisms      OrganismList
	species        SpeciesList
	env            Environment
	speciesLastID  int
	organismLastID int
	iteration      int
	maxCells       int
	areaCount      int
	debug          bool
	lock           sync.Mutex
	speciesLock    sync.Mutex
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
		organisms := 0
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
		organisms = aliveOrganisms.GetAreaCount(start, end)
		countSpan.Finish()
		areas[areaIndex] = organisms < (s.maxCells / s.areaCount / s.areaCount)
	}

	return areas
}

func (s *Sim) Create(verbose bool) {
	s.iteration = 0
	s.env = Environment{4, 10000, 10000}
	startCellCount := 5

	startCells := make(OrganismList, startCellCount)

	for i := 0; i < startCellCount; i++ {
		startCells[i] = getRandomOrganism(i, s.env, s.addSpecies)
	}

	s.organisms = startCells
	s.maxCells = 2e4
	s.areaCount = 5
	s.debug = verbose
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
		CellCount:      len(s.organisms),
		AliveCellCount: s.GetAliveCount(),
		Iteration:      s.iteration,
		Waste: WasteData{
			MinTolerance: s.species[0].types[0].WasteTolerance,
			Waste:        s.env.toxicity,
		},
		Procreation: ProcreationData{
			MinCd:     s.species[0].types[0].procreationCd,
			MinHeight: float64(s.env.height),
		},
	}
	dataSpan.Finish()

	data.Procreation.CanProcreate = data.AliveCellCount < s.maxCells
	index := 0

	for _, species := range s.species {
		if !species.extinct {
			for tIndex := range species.types {
				if data.Waste.MaxTolerance < species.types[tIndex].WasteTolerance {
					data.Waste.MaxTolerance = species.types[tIndex].WasteTolerance
				}
				if data.Waste.MinTolerance > species.types[tIndex].WasteTolerance {
					data.Waste.MinTolerance = species.types[tIndex].WasteTolerance
				}
			}
		}
	}

	areas := s.getAreas(stepSpanCtx)

	simSpan, _ := opentracing.StartSpanFromContext(stepSpanCtx, "sim")
	for organismIndex, organism := range s.organisms {
		if organism.IsAlive() {
			if data.Procreation.MaxHeight < organism.position.Y {
				data.Procreation.MaxHeight = organism.position.Y
			}
			if data.Procreation.MinHeight > organism.position.Y {
				data.Procreation.MinHeight = organism.position.Y
			}
		}

		pos := s.organisms[organismIndex].position
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
			s.env,
			s.iteration,
			s.addSpecies,
			canProcreate,
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
	simSpan.Finish()

	s.organisms = nextGenOrganisms[:index]
	s.env.changeToxicity(waste)

	s.cleanupSpecies(stepSpanCtx)

	data.Procreation.Species = s.species

	if s.debug {
		fmt.Printf(
			"Iteration %6d, organisms: %5d, alive: %5d, waste: %.4f %d species\n",
			s.iteration,
			len(s.organisms),
			s.GetAliveCount(),
			s.env.toxicity,
			len(s.GetSpecies().GetAlive()),
		)
	}

	return data
}

func (s *Sim) RunLoop(data *IterationData) {

	for {
		span := opentracing.GlobalTracer().StartSpan("loop")
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		s.lock.Lock()
		iterationData := s.RunStep(ctx)
		data.from(iterationData)

		s.lock.Unlock()

		if s.GetCellCount() == 0 {
			break
		}

		span.Finish()

		if iterationData.Iteration > 380 && !s.debug {
			time.Sleep(time.Second)
		}
	}
}
