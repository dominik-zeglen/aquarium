package api

import (
	"context"
	"log"
	"testing"

	"github.com/dominik-zeglen/aquarium/sim"
)

func TestCellResolver(t *testing.T) {
	// Given
	s := sim.Sim{}
	s.Create(sim.SimConfig{})
	d := sim.IterationData{}
	schema, err := GetSchema(&s, &d)
	if err != nil {
		t.Fatal(err)
	}

	// When
	variables := map[string]interface{}{}

	res := schema.Exec(
		context.TODO(),
		`query GetCell {
			organism(id: 1) {
				id
				cells {
					id
					alive
					hp
					position {
						x
						y
					}
					type {
						id
						diet
						funghi
						herbivore
					}
				}
			}
		}`,
		"GetCell",
		variables,
	)

	// Then
	if len(res.Errors) > 0 {
		log.Fatal(res.Errors)
	}
}
