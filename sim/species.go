package sim

import (
	"fmt"
	"math/rand"
)

type Species struct {
	ID        int `json:"id"`
	emergedAt int
	extinct   bool

	shape string

	size int

	membrane int
	enzymes  int

	Herbivore int8 `json:"herbivore"`
	Carnivore int8 `json:"carnivore"`
	Funghi    int8 `json:"funghi"`

	timeToDie      int
	wasteTolerance float64

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
	return fmt.Sprintf("%s-%d-%d", diet, s.emergedAt, s.ID)
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

func (s Species) getProcessedWaste(e Environment) float64 {
	if e.toxicity > 0 {
		waste := (float64(s.Funghi)) * .2 * (e.toxicity - .5)
		if waste > 0 {
			return waste
		}
	}

	return 0
}

func (c Species) getWaste(e Environment) float64 {
	waste := float64(c.size)
	if (e.toxicity) > 0 {
		waste -= c.getProcessedWaste(e)
	}

	return waste / 6e8
}

func (c Species) getWasteAfterDeath() float64 {
	return (float64(c.size)) / 8e8
}

func (s *Species) validate() {
	if s.Carnivore < 0 {
		s.Carnivore = 0
	}
	if s.Herbivore < 0 {
		s.Herbivore = 0
	}
	if s.Funghi < 0 {
		s.Funghi = 0
	}
	if s.consumption < 3 {
		s.consumption = 1
	}
	if s.division < 0 {
		s.division = 0
	}
}

func (s Species) mutate() Species {
	attr := rand.Intn(7)

	switch attr {
	case 0:
		s.Carnivore += int8(rand.Intn(3) - 1)
		break
	case 1:
		s.Herbivore += int8(rand.Intn(3) - 1)
		break
	case 2:
		s.Funghi += int8(rand.Intn(3) - 1)
		break
	case 3:
		s.maxCapacity += rand.Intn(3) - 1
		break
	case 4:
		s.consumption += rand.Intn(3) - 1
		break
	case 5:
		s.procreationCd += int8(rand.Intn(3) - 1)
		break
	case 6:
		s.wasteTolerance += float64(rand.Intn(3)-1) / 4
		break
	case 7:
		s.timeToDie += rand.Intn(3) - 1
		break
	}

	s.validate()
	if rand.Float32() > .5 {
		s.mutate()
	}

	return s
}

func getRandomFunghi() Species {
	points := int((17 * 100) / (3 + rand.Float64()))

	s := Species{}

	s.maxCapacity = rand.Intn(30)
	points -= s.maxCapacity

	s.size = rand.Intn(30) + 1
	points -= s.size

	s.Funghi = int8(rand.Intn(100))
	points -= int(s.Funghi)

	s.division = 1
	points -= int(s.division * 10)

	s.timeToDie = int(rand.NormFloat64()*20) + 30
	points -= s.timeToDie * 10

	s.maxSatiation = int(rand.Intn(100)) + 300

	s.consumption = 10
	s.procreationCd = int8(rand.Intn(4) + 8)

	s.wasteTolerance = float64(rand.Intn(3)-1)*2 + 10

	return s
}

func getRandomHerbivore() Species {
	points := int((17 * 100) / (3 + rand.Float64()))

	s := Species{}

	s.maxCapacity = rand.Intn(30)
	points -= s.maxCapacity

	s.size = rand.Intn(30) + 1
	points -= s.size

	s.Herbivore = int8(rand.Intn(100))
	points -= int(s.Herbivore)

	s.division = 1
	points -= int(s.division * 10)

	s.timeToDie = int(rand.NormFloat64()*20) + 30
	points -= s.timeToDie * 10

	s.maxSatiation = int(rand.Intn(100)) + 300

	s.consumption = 10
	s.procreationCd = int8(rand.Intn(4) + 8)

	s.wasteTolerance = float64(rand.Intn(3)-1)*6 + 4

	return s
}
