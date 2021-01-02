package sim

import "fmt"

func ExampleEating() {
	envs := []Environment{
		{0, 10, 10},
		{1, 10, 10},
		{2, 10, 10},
		{5, 10, 10},
		{10, 10, 10},
	}

	species := []Species{
		{
			diets:     []Diet{Herbivore},
			Herbivore: 25,
		},
		{
			diets:     []Diet{Herbivore},
			Herbivore: 50,
		},
		{
			diets:     []Diet{Herbivore},
			Herbivore: 75,
		},
		{
			diets:     []Diet{Herbivore},
			Herbivore: 100,
		},
		{
			diets:  []Diet{Funghi},
			Funghi: 25,
		},
		{
			diets:  []Diet{Funghi},
			Funghi: 50,
		},
		{
			diets:  []Diet{Funghi},
			Funghi: 75,
		},
		{
			diets:  []Diet{Funghi},
			Funghi: 100,
		},
	}

	cells := make(CellList, len(species)*3)
	for cellIndex := range cells {
		cells[cellIndex].position.Y = float64(cellIndex%3)*3 + 1.5
		cells[cellIndex].species = &species[cellIndex/3]
	}

	for _, env := range envs {
		for _, cell := range cells {
			f := 0
			d := cell.species.Herbivore
			if cell.species.Funghi > 0 {
				d = cell.species.Funghi
			}

			for it := 0; it < 24; it++ {
				f += cell.GetFood(env, it)
			}

			fmt.Printf(
				"%s: %d ate %d at %.2f toxicity at -%.1f depth\n",
				cell.species.diets[0],
				d,
				f/24,
				env.toxicity,
				cell.position.Y,
			)
		}
	}
}
