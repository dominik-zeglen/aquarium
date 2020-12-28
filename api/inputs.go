package api

import "github.com/golang/geo/r2"

type AreaInput struct {
	Start r2.Point
	End   r2.Point
	Scale *int32
}

type CellFilter struct {
	Area *AreaInput
}

type SpeciesFilter struct {
	Area *AreaInput
}
