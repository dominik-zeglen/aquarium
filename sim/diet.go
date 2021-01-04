package sim

type Diet int

const (
	Herbivore Diet = iota
	Funghi
)

func (d Diet) String() string {
	return [...]string{"herbivore", "funghi"}[d]
}

func (t CellType) hasDiet(diet Diet) bool {
	for _, sd := range t.diets {
		if sd == diet {
			return true
		}
	}

	return false
}

func (s Species) hasDiet(diet Diet) bool {
	for _, t := range s.Types {
		if t.hasDiet(diet) {
			return true
		}
	}

	return false
}
