package sudokgo

import (
	"errors"
	"fmt"
)

const (
	GridSize = 3
	rowSize  = GridSize*GridSize

	RuleSimple     = 1
	RuleEasy       = 10
	RuleMedium     = 100
	RuleHard       = 1000
	RuleDiabolical = 10000

	ScoreSimple     = 25
	ScoreEasy       = 150
	ScoreMedium     = 900
	ScoreHard       = 5400
	ScoreImpossible = 99999
)

var (
	ErrLoadSize = errors.New(fmt.Sprintf("incorrect number of items to load, expcted %d",
		rowSize*rowSize))
	ErrEmptyGrid   = errors.New("grid must not be empty")
	ErrCannotSolve = errors.New("puzzle is too complex to solve")
	ErrImpossible  = errors.New("puzzle is impossible, broken")
)

type Possible struct {
	count         int
	possibilities [rowSize]bool
}

type Sudoku struct {
	Grid     [rowSize][rowSize]int
	possGrid [rowSize][rowSize]Possible
}

func (s *Sudoku) Reset() {
	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			s.Grid[x][y] = -1
		}
	}

	s.loadPossGrid()
}

func (s *Sudoku) Load(numbers string) error {
	if len(numbers) != rowSize*rowSize {
		return ErrLoadSize
	}

	offset := 0
	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			n := numbers[offset]

			if n == '-' {
				s.Grid[x][y] = -1
			} else {
				s.Grid[x][y] = int(n) - 48
			}

			offset++
		}
	}

	return nil
}

func (s *Sudoku) Solve() (int, error) {
	solved := true
	broken := false

	/* catch folk who might try and Solve an empty grid, it just wastes time */
	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			if s.Grid[x][y] != -1 {
				solved = false
				break
			}
		}

		if !solved {
			break
		}
	}
	if solved {
		return ScoreImpossible, ErrEmptyGrid
	}

	solved = true
	s.loadPossGrid()

	score := 0
	trying := true
	for trying {
		for _, rule := range rules {
			madeMove, points := rule(s)
			trying = madeMove

			if madeMove {
				score += points
				break
			}
		}
	}

	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			if s.Grid[x][y] == -1 {
				solved = false
				if s.possGrid[x][y].count == 0 {
					broken = true
				}
			}
		}
	}

	if !solved {
		s.printPossGrid()
		return ScoreImpossible, ErrCannotSolve
	} else if broken {
		return ScoreImpossible, ErrImpossible
	}

	return score, nil
}

func (s *Sudoku) Print() {
	for y := 0; y < rowSize; y++ {
		if y > 0 && (y%GridSize) == 0 {
			for x := 0; x < rowSize+GridSize-1; x++ {
				fmt.Print("-")
			}
			fmt.Println()
		}

		for x := 0; x < rowSize; x++ {
			if x > 0 && (x%GridSize) == 0 {
				fmt.Print("|")
			}

			if s.Grid[x][y] == -1 {
				fmt.Print(" ")
			} else {
				fmt.Print(s.Grid[x][y])
			}
		}

		fmt.Println()
	}
}

func NewSudoku() *Sudoku {
	return &Sudoku{}
}
