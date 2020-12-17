package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
		d, _ := json.Marshal(data[len(data)-1])
		res.Write(d)
	})
	go http.ListenAndServe(":8000", nil)

	for {
		if s.GetCellCount() == 0 {
			break
		}
		data = append(data, s.RunStep())
		// if data[len(data)-1].Iteration > 1000 {
		// 	break
		// }
	}

	defer save(data)
}
