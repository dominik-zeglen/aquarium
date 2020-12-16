package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dominik-zeglen/aquarium/sim"
)

func save(data []sim.SimData) {
	file, _ := json.Marshal(data)

	_ = ioutil.WriteFile("out/data.json", file, 0644)
}

func main() {
	s := sim.CreateSim()
	data := []sim.SimData{}

	for {
		if s.GetCellCount() == 0 {
			break
		}
		data = append(data, s.RunStep())
		if data[len(data)-1].Iteration > 1000 {
			break
		}
	}

	defer save(data)
}
