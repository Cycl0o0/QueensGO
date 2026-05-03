package game

// CellMark represents the current mark on a cell.
type CellMark int

const (
	MarkEmpty CellMark = iota
	MarkX
	MarkQueen
)

// Board holds all the state for a Queens puzzle.
type Board struct {
	Size    int        // N (grid is NxN)
	Regions [][]int    // region ID for each cell [row][col]
	Marks   [][]CellMark // player marks [row][col]
}

// NewBoard creates a board with the given size and region map.
func NewBoard(size int, regions [][]int) *Board {
	marks := make([][]CellMark, size)
	for i := range marks {
		marks[i] = make([]CellMark, size)
	}
	return &Board{
		Size:    size,
		Regions: regions,
		Marks:   marks,
	}
}

// CycleMark advances the mark at (row, col): empty → X → queen → empty.
func (b *Board) CycleMark(row, col int) {
	b.Marks[row][col] = (b.Marks[row][col] + 1) % 3
}

// Queens returns a list of (row, col) positions where queens are placed.
func (b *Board) Queens() [][2]int {
	var qs [][2]int
	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			if b.Marks[r][c] == MarkQueen {
				qs = append(qs, [2]int{r, c})
			}
		}
	}
	return qs
}

// Conflicts returns the set of queen positions that violate constraints.
func (b *Board) Conflicts() map[[2]int]bool {
	conflicts := map[[2]int]bool{}
	queens := b.Queens()

	for i, q1 := range queens {
		for j, q2 := range queens {
			if i == j {
				continue
			}
			// Same row
			if q1[0] == q2[0] {
				conflicts[q1] = true
				conflicts[q2] = true
			}
			// Same column
			if q1[1] == q2[1] {
				conflicts[q1] = true
				conflicts[q2] = true
			}
			// Same region
			if b.Regions[q1[0]][q1[1]] == b.Regions[q2[0]][q2[1]] {
				conflicts[q1] = true
				conflicts[q2] = true
			}
			// Adjacent (including diagonal) — kings move distance
			dr := q1[0] - q2[0]
			dc := q1[1] - q2[1]
			if dr >= -1 && dr <= 1 && dc >= -1 && dc <= 1 {
				conflicts[q1] = true
				conflicts[q2] = true
			}
		}
	}
	return conflicts
}

// IsWon returns true if the puzzle is correctly solved:
// exactly N queens, one per row, col, and region, no adjacency conflicts.
func (b *Board) IsWon() bool {
	queens := b.Queens()
	if len(queens) != b.Size {
		return false
	}
	return len(b.Conflicts()) == 0
}

// Reset clears all marks.
func (b *Board) Reset() {
	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			b.Marks[r][c] = MarkEmpty
		}
	}
}
