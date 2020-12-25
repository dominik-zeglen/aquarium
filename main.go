package main

import (
	"encoding/json"
	"net/http"

	"github.com/dominik-zeglen/aquarium/api"
	"github.com/dominik-zeglen/aquarium/sim"
)

func main() {
	s := sim.Sim{}
	s.Create()

	var data sim.IterationData

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("content-type", "application/json")
		res.Header().Add("Access-Control-Allow-Origin", "*")
		d, _ := json.Marshal(data)
		res.Write(d)
	})
	http.Handle("/api", api.InitAPI(&s, &data))
	go http.ListenAndServe(":8000", nil)

	consecutiveNoProcreateIterations := 0

	for {
		if s.GetCellCount() == 0 {
			break
		}
		iterationData := s.RunStep()
		data = iterationData

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

		// if true {
		// 	time.Sleep(time.Second / 8)
		// }
	}
}
