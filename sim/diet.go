package sim

type Diet int

const (
	Herbivore Diet = iota
	Funghi
)

func (d Diet) String() string {
	return [...]string{"herbivore", "funghi"}[d]
}

func HasDiet(diet Diet, diets []Diet) bool {
	for _, sd := range diets {
		if sd == diet {
			return true
		}
	}

	return false
}

func (t CellType) hasDiet(diet Diet) bool {
	return HasDiet(diet, t.diets)
}

func (s Species) hasDiet(diet Diet) bool {
	for _, t := range s.types {
		if t.hasDiet(diet) {
			return true
		}
	}

	return false
}
