package game

// Solve finds solutions to the queens puzzle. It stops after finding `limit`
// solutions (0 = find all). Returns the list of solutions, where each solution
// is a slice of column indices (solution[row] = col).
func Solve(size int, regions [][]int, limit int) [][]int {
	var solutions [][]int
	colUsed := make([]bool, size)
	regionUsed := make([]bool, size)
	placement := make([]int, size)

	var backtrack func(row int)
	backtrack = func(row int) {
		if limit > 0 && len(solutions) >= limit {
			return
		}
		if row == size {
			sol := make([]int, size)
			copy(sol, placement)
			solutions = append(solutions, sol)
			return
		}
		for col := 0; col < size; col++ {
			if colUsed[col] {
				continue
			}
			reg := regions[row][col]
			if regionUsed[reg] {
				continue
			}
			// Check diagonal adjacency with the queen in the previous row
			if row > 0 {
				prevCol := placement[row-1]
				if col >= prevCol-1 && col <= prevCol+1 {
					continue
				}
			}
			placement[row] = col
			colUsed[col] = true
			regionUsed[reg] = true

			backtrack(row + 1)

			colUsed[col] = false
			regionUsed[reg] = false
		}
	}

	backtrack(0)
	return solutions
}

// HasUniqueSolution returns true if the puzzle has exactly one solution.
func HasUniqueSolution(size int, regions [][]int) bool {
	return len(Solve(size, regions, 2)) == 1
}
