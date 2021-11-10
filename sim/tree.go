package sim

type CellGraph struct {
	m [][]*Cell
	x int
	y int
}

func (cg *CellGraph) Init(x int, y int, m [][]*Cell) {
	cg.x = x
	cg.y = y
	cg.m = m
}

func (cg *CellGraph) isSafe(x int, y int, visited [][]bool) bool {
	isInBounds := x >= 0 && x < cg.x && y >= 0 && y < cg.y

	if !isInBounds {
		return false
	}

	return !visited[y][x] && cg.m[y][x] != nil
}

func (cg *CellGraph) dfs(x int, y int, visited [][]bool) []*Cell {
	rowNbr := []int{-1, 0, 0, 1}
	colNbr := []int{0, -1, 1, 0}

	visited[y][x] = true
	subGraph := []*Cell{cg.m[y][x]}

	for i := 0; i < 4; i++ {
		if cg.isSafe(x+rowNbr[i], y+colNbr[i], visited) {
			subGraph = append(
				subGraph,
				cg.dfs(x+rowNbr[i], y+colNbr[i], visited)...,
			)
		}
	}

	return subGraph
}

func (cg *CellGraph) GetIslands() [][]*Cell {
	visited := make([][]bool, cg.y)
	for i := 0; i < cg.y; i++ {
		visited[i] = make([]bool, cg.x)
	}

	islands := [][]*Cell{}

	for y := 0; y < cg.y; y++ {
		for x := 0; x < cg.x; x++ {
			if !visited[y][x] && cg.m[y][x] != nil {
				islands = append(islands, cg.dfs(x, y, visited))
			}
		}
	}

	return islands
}
