package sim

import (
	"fmt"
	"math/rand"
)

type Species struct {
	ID        int  `json:"id"`
	EmergedAt int  `json:"emergedAt"`
	Extinct   bool `json:"extinct"`
	Count     int  `json:"count"`

	shape string

	size int

	membrane int
	enzymes  int

	Herbivore int8 `json:"herbivore"`
	Carnivore int8 `json:"carnivore"`
	Funghi    int8 `json:"funghi"`

	TimeToDie      int     `json:"timeToDie"`
	WasteTolerance float64 `json:"wasteTolerance"`

	maxSatiation int
	consumption  int
	transport    int
	maxCapacity  int

	division      int8
	connectivity  int8
	procreationCd int8

	mobility int

	// True if cell reproduces by budding
	reproducingMethod bool
}

func (s Species) getName() string {
	diet := "F"
	if s.Herbivore > s.Funghi {
		diet = "H"
	}
	return fmt.Sprintf("%s-%d-%d", diet, s.EmergedAt, s.ID)
}

func (s Species) getMaxHP() int {
	return s.size * 23
}

func (s Species) getMass() int {
	return s.size * 10
}

func (s Species) getFoodValue() int {
	return s.size * 2
}

func (s Species) getDefence() int {
	return s.membrane / 10
}

func (s Species) getAttack() int {
	return (int(s.Carnivore)*2 + s.size + s.enzymes) / 3
}

func (s Species) getProcessedWaste(toxicity float64) float64 {
	if toxicity > 0 && s.Funghi > 0 {
		waste := (float64(s.Funghi)) * .75 * (toxicity - .5)
		if waste > 0 {
			return waste
		}
	}

	return 0
}

func (s Species) getWaste(toxicity float64) float64 {
	waste := float64(s.size)
	if (toxicity) > 0 {
		waste -= s.getProcessedWaste(toxicity)
	}

	return waste / 6e8
}

func (c Species) getWasteAfterDeath() float64 {
	return (float64(c.size)) / 6e8
}

func (c Species) getConsumption() int {
	return int(float32(c.maxSatiation) / 20 * float32(c.size) / 30)
}

func (s *Species) validate() bool {
	if s.Carnivore < 0 {
		s.Carnivore = 0
		return false
	}
	if s.Herbivore < 0 {
		s.Herbivore = 0
		return false
	}
	if s.Funghi < 0 {
		s.Funghi = 0
		return false
	}
	if s.consumption < 3 {
		s.consumption = 3
		return false
	}
	if s.division < 0 {
		s.division = 0
		return false
	}
	if s.TimeToDie > 60 {
		s.TimeToDie = 60
		return false
	}
	if s.procreationCd < 5 {
		s.procreationCd = 5
		return false
	}
	if s.maxSatiation < 50 {
		s.maxSatiation = 50
		return false
	}
	if s.size < 10 {
		s.size = 10
		return false
	}

	return true
}

func (s Species) getDietPoints() int {
	diets := 0
	dietPoints := 0

	if s.Carnivore > 0 {
		diets++
		dietPoints += int(s.Carnivore)
	}
	if s.Herbivore > 0 {
		diets++
		dietPoints += int(s.Herbivore)
	}
	if s.Funghi > 0 {
		diets++
		dietPoints += int(s.Funghi)
	}

	if diets > 2 {
		return dietPoints * 2
	}
	if diets > 1 {
		return dietPoints * 3 / 2
	}

	return dietPoints
}

func (s Species) mutate() Species {
	st := s

	if rand.Float32() > .99 {
		mutationCount := (rand.Intn(40) + 1)
		for i := 0; i < mutationCount; i++ {
			st = st.mutateOnce()
		}
	} else {
		do := true
		for do || rand.Float32() > .5 {
			st = st.mutateOnce()
			do = false
		}
	}

	return st
}

func (s Species) mutateOnce() Species {
	n := s
	do := true

	for do || !n.validate() {
		attr := rand.Float64()
		value := rand.Intn(2)*2 - 1

		for attr < .21 && n.getDietPoints() > 100 && value > 0 {
			attr = rand.Float64()
		}

		if attr < .07 {
			n.Carnivore += int8(value * 4)
		} else if attr < .14 {
			n.Herbivore += int8(value * 4)
		} else if attr < .21 {
			n.Funghi += int8(value * 4)
		}

		if attr > .21 && attr > .41 {
			n.maxCapacity += value
		}

		if attr > .41 && attr < .61 {
			n.consumption += value
		}

		if attr > .61 && attr < .63 {
			n.TimeToDie += value
		}

		if attr > .63 && attr < .73 {
			n.procreationCd += int8(value)
		}

		if attr > .73 && attr < .88 {
			n.WasteTolerance += float64(value) / 4
		}

		if attr > .88 && attr < .95 {
			n.maxSatiation += value
		}

		if attr > .95 {
			n.size += value
		}

		do = false
	}

	return n
}

func getRandomHerbivore() Species {
	s := Species{}

	s.maxCapacity = rand.Intn(30)

	s.size = rand.Intn(20) + 10

	s.Herbivore = int8(rand.Intn(20)) + 5

	s.division = 1

	s.TimeToDie = 30

	s.maxSatiation = int(rand.Intn(100)) + 300

	s.consumption = 10
	s.procreationCd = int8(rand.Intn(4) + 8)

	s.WasteTolerance = float64(rand.Intn(16))/4 + 1
	s.mobility = 20

	return s
}

func getRandomFunghi() Species {
	s := getRandomHerbivore()
	s.Funghi = s.Herbivore
	s.Herbivore = 0
	s.maxSatiation -= 50
	s.TimeToDie += 5

	return s
}
