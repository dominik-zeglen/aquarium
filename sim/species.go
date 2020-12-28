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
	diets []Diet

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

func (s Species) GetName() string {
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

func (c Species) GetConsumption() int {
	return int(
		float32(c.maxSatiation) / 20 *
			float32(c.size) / 30 *
			float32(c.consumption) / 10,
	)
}

func (s Species) GetDiet() []Diet {
	return s.diets
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

	if rand.Float32() > .9 {
		if len(st.diets) == 1 {
			if rand.Float32() > .9 {
				if st.hasDiet(Herbivore) {
					st.diets = append(st.diets, Funghi)
					st.Herbivore /= 2
					st.Funghi = st.Herbivore
				}
				if st.hasDiet(Funghi) {
					st.diets = append(st.diets, Herbivore)
					st.Funghi /= 2
					st.Herbivore = st.Funghi
				}
			} else {
				if st.hasDiet(Herbivore) {
					st.diets = []Diet{Funghi}
					st.Funghi = st.Herbivore
					st.Herbivore = 0
				}
				if st.hasDiet(Funghi) {
					st.diets = []Diet{Herbivore}
					st.Herbivore = st.Funghi
					st.Funghi = 0
				}
			}
		} else {
			diets := []Diet{Herbivore, Funghi}
			diet := diets[rand.Intn(len(diets))]
			if diet == Herbivore {
				st.Herbivore += st.Funghi
				st.Funghi = 0
			} else if diet == Funghi {
				st.Funghi += st.Herbivore
				st.Herbivore = 0
			}
		}

		mutationCount := (rand.Intn(10) + 10)
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

		if attr < .21 {
			if len(s.diets) > 1 {
				if rand.Float32() > .5 {
					n.Herbivore += int8(value * 4)
				} else {
					n.Funghi += int8(value * 4)
				}
			} else {
				if s.hasDiet(Herbivore) {
					n.Herbivore += int8(value * 4)
				}
				if s.hasDiet(Funghi) {
					n.Funghi += int8(value * 4)
				}
			}
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

		if attr > .73 && attr < .9 {
			n.WasteTolerance += float64(value) / 4
		}

		if attr > .9 && attr < .95 {
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
	s.diets = []Diet{Herbivore}

	s.maxCapacity = rand.Intn(30)

	s.size = rand.Intn(20) + 10

	s.Herbivore = int8(rand.Intn(20)) + 5

	s.division = 1

	s.TimeToDie = 30

	s.maxSatiation = int(rand.Intn(100)) + 300

	s.consumption = 10
	s.procreationCd = int8(rand.Intn(4) + 8)

	s.WasteTolerance = float64(rand.Intn(16))/4 + 4
	s.mobility = 20

	return s
}
