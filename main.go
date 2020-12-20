package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dominik-zeglen/aquarium/sim"
)

func save(data []sim.SimData) {
	file, _ := json.Marshal(data)

	_ = ioutil.WriteFile("out/data.json", file, 0644)
}

func main() {
	s := sim.Sim{}
	s.Create()

	data := []sim.SimData{}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("content-type", "application/json")
		res.Header().Add("Access-Control-Allow-Origin", "*")
		d, _ := json.Marshal(data[len(data)-1])
		res.Write(d)
	})
	go http.ListenAndServe(":8000", nil)

	consecutiveNoProcreateIterations := 0

	for {
		if s.GetCellCount() == 0 {
			break
		}
		iterationData := s.RunStep()
		data = append(data, iterationData)

		if !iterationData.Procreation.CanProcreate {
			consecutiveNoProcreateIterations++
		} else {
			consecutiveNoProcreateIterations = 0
		}

		if consecutiveNoProcreateIterations > 2 {
			s.KillOldestCells()
		}

		// if iterationData.Iteration == 1 {
		// 	break
		// }

		if false {
			time.Sleep(time.Second / 8)
		}
	}

	defer save(data)
}
