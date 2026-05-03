package game

import "testing"

func TestGenerateAndSolve(t *testing.T) {
	for _, size := range []int{5, 6, 7} {
		regions := GeneratePuzzle(size)
		if regions == nil {
			t.Fatalf("GeneratePuzzle(%d) returned nil", size)
		}
		solutions := Solve(size, regions, 2)
		if len(solutions) != 1 {
			t.Fatalf("Expected 1 solution for size %d, got %d", size, len(solutions))
		}
	}
}

func TestBoardConflicts(t *testing.T) {
	regions := [][]int{
		{0, 1, 2},
		{0, 1, 2},
		{0, 1, 2},
	}
	b := NewBoard(3, regions)
	// Place two queens in same row
	b.Marks[0][0] = MarkQueen
	b.Marks[0][2] = MarkQueen
	conflicts := b.Conflicts()
	if len(conflicts) == 0 {
		t.Fatal("Expected conflicts for queens in same row")
	}
}

func TestBoardWin(t *testing.T) {
	regions := GeneratePuzzle(5)
	solutions := Solve(5, regions, 1)
	if len(solutions) == 0 {
		t.Fatal("No solution found")
	}
	b := NewBoard(5, regions)
	for row, col := range solutions[0] {
		b.Marks[row][col] = MarkQueen
	}
	if !b.IsWon() {
		t.Fatal("Board should be won with correct solution")
	}
}
