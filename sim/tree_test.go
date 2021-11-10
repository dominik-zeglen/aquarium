package sim

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIslandCounting(t *testing.T) {
	// Given
	cell := Cell{}
	m := [][]*Cell{
		{&cell, &cell, nil, nil, nil},
		{&cell, nil, nil, nil, nil},
		{nil, nil, &cell, nil, nil},
		{&cell, &cell, nil, nil, nil},
		{nil, nil, nil, nil, &cell},
	}

	// When
	graph := CellGraph{}
	graph.Init(5, 5, m)
	islands := graph.GetIslands()

	fmt.Println(islands)

	// Then
	assert.Equal(t, 4, len(islands))
	assert.Equal(t, 3, len(islands[0]))
	assert.Equal(t, 1, len(islands[1]))
	assert.Equal(t, 2, len(islands[2]))
	assert.Equal(t, 1, len(islands[3]))
}

func TestNotRectangular(t *testing.T) {
	// Given
	cell := Cell{}
	m := [][]*Cell{
		{&cell, &cell, nil, nil, nil},
		{&cell, nil, nil, nil, nil},
		{nil, nil, nil, nil, &cell},
	}

	// When
	graph := CellGraph{}
	graph.Init(5, 3, m)
	islands := graph.GetIslands()

	// Then
	assert.Equal(t, 2, len(islands))
}
