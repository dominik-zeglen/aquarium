package sim

import "math"

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
	return e.toxicity / 2 * (height/float64(e.height) + 1)
}

func (e Environment) getLightOnHeight(height float64, iteration int) float64 {
	hour := (iteration / 10) % 24
	light := math.Abs(float64((hour - 12)))
	return light*(1-height/float64(e.height))*.8 + .2
}
