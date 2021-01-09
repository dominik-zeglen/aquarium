package sim

import (
	"math/rand"
)

type CellType struct {
	ID int

	shape string
	diets []Diet

	size int

	membrane int
	enzymes  int

	Herbivore int8
	Carnivore int8
	Funghi    int8

	TimeToDie      int8
	WasteTolerance float64

	maxSatiation int
	consumption  int
	transport    int
	maxCapacity  int

	connectivity  int8
	procreationCd int8

	mobility int
}

func (t CellType) getMaxHP() int {
	return t.size * 23
}

func (t CellType) getMass() int {
	return t.size * 10
}

func (t CellType) getFoodValue() int {
	return t.size * 2
}

func (t CellType) getDefence() int {
	return t.membrane / 10
}

func (t CellType) getAttack() int {
	return (int(t.Carnivore)*2 + t.size + t.enzymes) / 3
}

func (t CellType) getProcessedWaste(toxicity float64) float64 {
	if toxicity > 0 && t.Funghi > 0 {
		waste := (float64(t.Funghi)) * .75 * (toxicity - .5)
		if waste > 0 {
			return waste
		}
	}

	return 0
}

func (t CellType) getWaste(toxicity float64) float64 {
	waste := float64(t.size)
	if (toxicity) > 0 {
		waste -= t.getProcessedWaste(toxicity)
	}

	return waste / 6e8
}

func (c CellType) getWasteAfterDeath() float64 {
	return (float64(c.size)) / 6e8
}

func (c CellType) GetConsumption() int {
	return int(
		float32(c.maxSatiation) / 20 *
			float32(c.size) / 30 *
			float32(c.consumption) / 10,
	)
}

func (t CellType) GetDiet() []Diet {
	return t.diets
}

func (t *CellType) validate() bool {
	if t.Carnivore < 0 {
		t.Carnivore = 0
		return false
	}
	if t.Herbivore < 0 {
		t.Herbivore = 0
		return false
	}
	if t.Funghi < 0 {
		t.Funghi = 0
		return false
	}
	if t.consumption < 3 {
		t.consumption = 3
		return false
	}
	if t.TimeToDie > 90 {
		t.TimeToDie = 90
		return false
	}
	if t.procreationCd < 5 {
		t.procreationCd = 5
		return false
	}
	if t.procreationCd >= t.TimeToDie {
		t.procreationCd = int8(t.TimeToDie) - 1
		return false
	}
	if t.maxSatiation < 50 {
		t.maxSatiation = 50
		return false
	}
	if t.size < 10 {
		t.size = 10
		return false
	}
	if t.WasteTolerance < 0 {
		t.WasteTolerance = 0
		return false
	}
	if t.mobility < 0 {
		t.mobility = 0
		return false
	}

	return true
}

func (t CellType) getDietPoints() int {
	diets := len(t.diets)
	dietPoints := 0

	if t.Carnivore > 0 {
		diets++
		dietPoints += int(t.Carnivore)
	}
	if t.Herbivore > 0 {
		diets++
		dietPoints += int(t.Herbivore)
	}
	if t.Funghi > 0 {
		diets++
		dietPoints += int(t.Funghi)
	}

	if diets > 2 {
		return dietPoints * 2
	}
	if diets > 1 {
		return dietPoints * 3 / 2
	}

	return dietPoints
}

func (t *CellType) mutateDiet() {
	if len(t.diets) == 1 {
		if rand.Float32() > .9 {
			if t.hasDiet(Herbivore) {
				t.diets = append(t.diets, Funghi)
				t.Herbivore /= 2
				t.Funghi = t.Herbivore
			}
			if t.hasDiet(Funghi) {
				t.diets = append(t.diets, Herbivore)
				t.Funghi /= 2
				t.Herbivore = t.Funghi
			}
		} else {
			if t.hasDiet(Herbivore) {
				t.diets = []Diet{Funghi}
				t.Funghi = t.Herbivore
				t.Herbivore = 0
			}
			if t.hasDiet(Funghi) {
				t.diets = []Diet{Herbivore}
				t.Herbivore = t.Funghi
				t.Funghi = 0
			}
		}
	} else {
		diets := []Diet{Herbivore, Funghi}
		diet := diets[rand.Intn(len(diets))]
		if diet == Herbivore {
			t.Herbivore += t.Funghi
			t.Funghi = 0
		} else if diet == Funghi {
			t.Funghi += t.Herbivore
			t.Herbivore = 0
		}
	}
}

func (t CellType) mutate() CellType {
	ct := t
	ct.diets = make([]Diet, len(t.diets))
	copy(ct.diets, t.diets)

	if rand.Float32() > .9 {
		ct.mutateDiet()

		mutationCount := (rand.Intn(10) + 10)
		for i := 0; i < mutationCount; i++ {
			ct = ct.mutateOnce()
		}
	} else {
		do := true
		for do || rand.Float32() > .5 {
			ct = ct.mutateOnce()
			do = false
		}
	}

	return ct
}

func (t CellType) mutateOnce() CellType {
	n := t
	do := true

	for do || !n.validate() {
		attr := rand.Float64()
		value := rand.Intn(2)*2 - 1

		for attr < .21 && n.getDietPoints() >= 100 && value > 0 {
			attr = rand.Float64()
		}

		if attr < .21 {
			if len(t.diets) > 1 {
				if rand.Float32() > .5 {
					n.Herbivore += int8(value * 4)
				} else {
					n.Funghi += int8(value * 4)
				}
			} else {
				if t.hasDiet(Herbivore) {
					n.Herbivore += int8(value * 4)
				}
				if t.hasDiet(Funghi) {
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
			n.TimeToDie += int8(value)
		}

		if attr > .63 && attr < .73 {
			n.procreationCd += int8(value)
		}

		if attr > .73 && attr < .85 {
			n.WasteTolerance += float64(value) / 4
		}

		if attr > .85 && attr < .9 {
			n.maxSatiation += value
		}

		if attr > .9 && attr < .95 {
			n.mobility += value
		}

		if attr > .95 {
			n.size += value
		}

		do = false
	}

	return n
}
