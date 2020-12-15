package sim

type Environment struct {
	toxicity float64
}

func (e *Environment) changeToxicity(value float64) {
	e.toxicity += value
	if e.toxicity < 0 {
		e.toxicity = 0
	}
}
