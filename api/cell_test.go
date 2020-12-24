package api

import (
	"context"
	"testing"

	"github.com/dominik-zeglen/aquarium/sim"
)

func TestCellResolver(t *testing.T) {
	s := sim.Sim{}
	s.Create()
	schema, err := GetSchema(&s)
	if err != nil {
		t.Fatal(err)
	}

	variables := map[string]interface{}{}

	schema.Exec(
		context.TODO(),
		`query GetCell {
			cell(id: 1) {
				id
				alive
				bornAt
				capacity
				hp
				position
				satiation
				# species
			}
		}`,
		"GetCell",
		variables,
	)
}

func TestCellListResolver(t *testing.T) {
	s := sim.Sim{}
	s.Create()
	schema, err := GetSchema(&s)
	if err != nil {
		t.Fatal(err)
	}

	variables := map[string]interface{}{}

	schema.Exec(
		context.TODO(),
		`query GetCells {
			cellList {
				count
				edges {
					node {
						id
						alive
						bornAt
						capacity
						hp
						position
						satiation
						# species
					}
				}
			}
		}`,
		"GetCells",
		variables,
	)
}
