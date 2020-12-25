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
	s.Create()
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
			cell(id: 1) {
				id
				alive
				bornAt
				capacity
				hp
				position {
					x
					y
				}
				satiation
				# species
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

func TestCellListResolver(t *testing.T) {
	// Given
	s := sim.Sim{}
	s.Create()
	d := sim.IterationData{}
	schema, err := GetSchema(&s, &d)
	if err != nil {
		t.Fatal(err)
	}

	// When
	variables := map[string]interface{}{}

	res := schema.Exec(
		context.TODO(),
		`query GetCells {
			cellList {
				count
				edges {
					node {
						id
					}
				}
			}
		}`,
		"GetCells",
		variables,
	)

	// Then
	if len(res.Errors) > 0 {
		log.Fatal(res.Errors)
	}
}

func TestAreaResolver(t *testing.T) {
	// Given
	s := sim.Sim{}
	s.Create()
	d := sim.IterationData{}
	schema, err := GetSchema(&s, &d)
	if err != nil {
		t.Fatal(err)
	}

	// When
	variables, err := unmarshallVariables(`{
		"start": {
			"x": 0,
			"y": 0
		},
		"end": {
			"x": 0,
			"y": 6000
		}
	}`)
	if err != nil {
		t.Fatal(err)
	}

	res := schema.Exec(
		context.TODO(),
		`query GetArea($start: PointInput!, $end: PointInput!) {
			area(start: $start, end: $end) {
				id
			}
		}`,
		"GetArea",
		variables,
	)

	// Then
	if len(res.Errors) > 0 {
		log.Fatal(res.Errors)
	}
}
