package main

import (
	"net/http"

	"github.com/dominik-zeglen/aquarium/api"
	"github.com/dominik-zeglen/aquarium/middleware"
	"github.com/dominik-zeglen/aquarium/sim"
)

func main() {
	s := sim.Sim{}
	s.Create()

	var data sim.IterationData
	http.Handle("/api", middleware.WithCors(
		[]string{"http://localhost:3000"},
		api.InitAPI(&s, &data),
	))
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
