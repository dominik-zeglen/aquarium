package sim

type Diet int

const (
	Herbivore Diet = iota
	Funghi
)

func (d Diet) String() string {
	return [...]string{"herbivore", "funghi"}[d]
}

func (s Species) hasDiet(diet Diet) bool {
	for _, sd := range s.diets {
		if sd == diet {
			return true
		}
	}

	return false
}
