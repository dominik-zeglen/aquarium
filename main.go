package main

import "github.com/dominik-zeglen/aquarium/sim"

func main() {
	s := sim.CreateSim()

	for {
		if s.GetCellCount() == 0 {
			break
		}
		s.RunStep()
	}
}
