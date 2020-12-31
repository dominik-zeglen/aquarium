package sim

type AddSpecies func(species Species) *Species

type WasteData struct {
	Waste        float64 `json:"waste"`
	MinTolerance float64 `json:"minTolerance"`
	MaxTolerance float64 `json:"maxTolerance"`
}

type ProcreationData struct {
	CanProcreate bool      `json:"canProcreate"`
	MinCd        int8      `json:"minCd"`
	MaxCd        int8      `json:"maxCd"`
	MinHeight    float64   `json:"minHeight"`
	MaxHeight    float64   `json:"maxHeight"`
	Species      []Species `json:"species"`
}

type IterationData struct {
	CellCount      int             `json:"cellCount"`
	AliveCellCount int             `json:"aliveCellCount"`
	Waste          WasteData       `json:"waste"`
	Iteration      int             `json:"iteration"`
	Procreation    ProcreationData `json:"procreation"`
}
