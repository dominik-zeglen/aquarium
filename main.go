package main

import "github.com/dominik-zeglen/aquarium/sim"

func main() {
	s := sim.CreateSim()

	for i := 0; i < 1000; i++ {
		s.RunStep()
	}
}
