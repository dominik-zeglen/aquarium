package api

import (
	"github.com/dominik-zeglen/aquarium/sim"
)

type IterationWasteResolver struct {
	w *sim.WasteData
	s *sim.Sim
}

func CreateIterationWasteResolver(data *sim.WasteData, sim *sim.Sim) IterationWasteResolver {
	return IterationWasteResolver{data, sim}
}

func (res IterationWasteResolver) MaxTolerance() float64 {
	return res.w.MaxTolerance
}
func (res IterationWasteResolver) MinTolerance() float64 {
	return res.w.MinTolerance
}
func (res IterationWasteResolver) Toxicity() float64 {
	return res.w.Waste
}

type IterationProcreationResolver struct {
	p *sim.ProcreationData
	s *sim.Sim
}

func CreateIterationProcreationResolver(data *sim.ProcreationData, sim *sim.Sim) IterationProcreationResolver {
	return IterationProcreationResolver{data, sim}
}

func (res IterationProcreationResolver) CanProcreate() bool {
	return res.p.CanProcreate
}
func (res IterationProcreationResolver) MaxCd() int32 {
	return int32(res.p.MaxCd)
}
func (res IterationProcreationResolver) MaxHeight() float64 {
	return res.p.MaxHeight
}
func (res IterationProcreationResolver) MinCd() int32 {
	return int32(res.p.MinCd)
}
func (res IterationProcreationResolver) MinHeight() float64 {
	return res.p.MinHeight
}
func (res IterationProcreationResolver) Species() []SpeciesResolver {
	return createSpeciesResolverList(res.p.Species)
}

type IterationResolver struct {
	d *sim.IterationData
	s *sim.Sim
}

func CreateIterationResolver(data *sim.IterationData, sim *sim.Sim) IterationResolver {
	return IterationResolver{data, sim}
}

func (res IterationResolver) AliveCellCount() int32 {
	return int32(res.d.AliveCellCount)
}

func (res IterationResolver) CellCount() int32 {
	return int32(res.d.CellCount)
}

func (res IterationResolver) Number() int32 {
	return int32(res.d.Iteration)
}

func (res IterationResolver) Procreation() IterationProcreationResolver {
	return CreateIterationProcreationResolver(&res.d.Procreation, res.s)
}

func (res IterationResolver) Waste() IterationWasteResolver {
	return CreateIterationWasteResolver(&res.d.Waste, res.s)
}
