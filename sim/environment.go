package sim

type Environment struct {
	toxicity float64
	width    int
	height   int
}

func (e *Environment) changeToxicity(value float64) {
	e.toxicity += value
	if e.toxicity < 0 {
		e.toxicity = 0
	}
}

func (e Environment) getToxicityOnHeight(height float64) float64 {
	return height / float64(e.height) * e.toxicity
}
