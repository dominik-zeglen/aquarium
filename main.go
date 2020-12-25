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
		[]string{"http://localhost:3000", "http://aquarium.unicorn-lab.online"},
		api.InitAPI(&s, &data),
	))
	go http.ListenAndServe(":8000", nil)

	s.RunLoop(&data)
}
