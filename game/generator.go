package game

import (
	"math/rand"
)

type cell struct{ r, c int }

// GeneratePuzzle creates a new puzzle of the given size.
// It returns the region map (regions[row][col] = regionID 0..size-1).
func GeneratePuzzle(size int) [][]int {
	for {
		queens := generateQueenPlacement(size)
		if queens == nil {
			continue
		}
		regions := buildRegions(size, queens)
		if regions == nil {
			continue
		}
		if HasUniqueSolution(size, regions) {
			return regions
		}
	}
}

// generateQueenPlacement places N queens (one per row, one per col, no diagonal adjacency).
func generateQueenPlacement(size int) []int {
	var colOrders [][]int
	for i := 0; i < size; i++ {
		colOrders = append(colOrders, rand.Perm(size))
	}

	result := make([]int, size)
	colUsed := make([]bool, size)

	var solve func(row int) bool
	solve = func(row int) bool {
		if row == size {
			return true
		}
		for _, col := range colOrders[row] {
			if colUsed[col] {
				continue
			}
			if row > 0 {
				prevCol := result[row-1]
				if col >= prevCol-1 && col <= prevCol+1 {
					continue
				}
			}
			result[row] = col
			colUsed[col] = true
			if solve(row + 1) {
				return true
			}
			colUsed[col] = false
		}
		return false
	}

	if solve(0) {
		return result
	}
	return nil
}

// buildRegions creates contiguous regions around each queen position using random growth.
func buildRegions(size int, queens []int) [][]int {
	regions := make([][]int, size)
	for r := range regions {
		regions[r] = make([]int, size)
		for c := range regions[r] {
			regions[r][c] = -1
		}
	}

	// Seed each region with its queen's cell
	frontiers := make([][]cell, size)
	for r, c := range queens {
		regions[r][c] = r
		for _, nb := range neighbors(r, c, size) {
			frontiers[r] = append(frontiers[r], nb)
		}
	}

	// Round-robin growth: each region claims one cell per round
	unassigned := size*size - size
	staleRounds := 0

	for unassigned > 0 {
		grewThisRound := false
		order := rand.Perm(size)
		for _, reg := range order {
			// Clean stale frontier entries and try to claim one cell
			claimed := false
			rand.Shuffle(len(frontiers[reg]), func(i, j int) {
				frontiers[reg][i], frontiers[reg][j] = frontiers[reg][j], frontiers[reg][i]
			})
			newFrontier := frontiers[reg][:0]
			for _, fc := range frontiers[reg] {
				if regions[fc.r][fc.c] != -1 {
					continue // already taken
				}
				if !claimed {
					regions[fc.r][fc.c] = reg
					unassigned--
					claimed = true
					grewThisRound = true
					// Add neighbors of claimed cell
					for _, nb := range neighbors(fc.r, fc.c, size) {
						if regions[nb.r][nb.c] == -1 {
							newFrontier = append(newFrontier, nb)
						}
					}
				} else {
					newFrontier = append(newFrontier, fc)
				}
			}
			frontiers[reg] = newFrontier
		}
		if !grewThisRound {
			staleRounds++
			if staleRounds > 5 {
				// Some cells are unreachable by any frontier — assign them to nearest region
				assignOrphans(size, regions)
				return regions
			}
		} else {
			staleRounds = 0
		}
	}

	return regions
}

// assignOrphans assigns any remaining -1 cells to an adjacent region.
func assignOrphans(size int, regions [][]int) {
	changed := true
	for changed {
		changed = false
		for r := 0; r < size; r++ {
			for c := 0; c < size; c++ {
				if regions[r][c] != -1 {
					continue
				}
				// Find any adjacent assigned region
				for _, nb := range neighbors(r, c, size) {
					if regions[nb.r][nb.c] != -1 {
						regions[r][c] = regions[nb.r][nb.c]
						changed = true
						break
					}
				}
			}
		}
	}
}

func neighbors(r, c, size int) []cell {
	var result []cell
	for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		nr, nc := r+d[0], c+d[1]
		if nr >= 0 && nr < size && nc >= 0 && nc < size {
			result = append(result, cell{nr, nc})
		}
	}
	return result
}
