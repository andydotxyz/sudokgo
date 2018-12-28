package sudokgo

import "fmt"

func (s *Sudoku) loadPossGrid() {
	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			if s.Grid[x][y] == -1 {
				s.possGrid[x][y].count = rowSize
				for z := 0; z < rowSize; z++ {
					s.possGrid[x][y].possibilities[z] = true
				}
			} else {
				s.possGrid[x][y].count = 0
				for z := 0; z < rowSize; z++ {
					s.possGrid[x][y].possibilities[z] = false
				}
			}
		}
	}

	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			s.updatePossibilities(x, y)
		}
	}
}

func (s *Sudoku) updatePossibilities(x, y int) {
	var zx, zy int
	a := s.Grid[x][y]

	//	cell := gridRef(x, y)
	if a != -1 {
		/* delete the possibility from the rest of the row... */
		for zx := 0; zx < rowSize; zx++ {
			if zx == x {
				continue
			}

			if s.possGrid[zx][y].possibilities[a-1] {
				s.possGrid[zx][y].possibilities[a-1] = false
				s.possGrid[zx][y].count--
				//				fmt.Println(gridRef(zx, y), a, "is not valid as it already appears in the row at", cell)
			}
		}

		/* delete the possibility from the rest of the col... */
		for zy = 0; zy < rowSize; zy++ {
			if zy == y {
				continue
			}

			if s.possGrid[x][zy].possibilities[a-1] {
				s.possGrid[x][zy].possibilities[a-1] = false
				s.possGrid[x][zy].count--
				//				fmt.Println(gridRef(x, zy), a, "is not valid as it already appears in the column at", cell)
			}
		}

		/* delete the possibility from the rest of the box... */
		{
			sx := x - x%GridSize
			sy := y - y%GridSize

			for zy = sy; zy < sy+GridSize; zy++ {
				for zx = sx; zx < sx+GridSize; zx++ {
					if zx == x && zy == y {
						continue
					}

					if s.possGrid[zx][zy].possibilities[a-1] {
						s.possGrid[zx][zy].possibilities[a-1] = false
						s.possGrid[zx][zy].count--
						//						fmt.Println(gridRef(zx, zy), a, "is not valid as it already appears in the box at", cell)
					}
				}
			}
		}
	}
}

func (s *Sudoku) printPossGrid() {
	size := 0
	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			if s.possGrid[x][y].count > size {
				size = s.possGrid[x][y].count
			}
		}
	}

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

			printed := printPossible(s.possGrid[x][y])
			for p := printed; p < size; p++ {
				fmt.Print(" ")
			}
			fmt.Print(" ")

			if s.Grid[x][y] == -1 {
				fmt.Print(" ")
			} else {
				fmt.Print(s.Grid[x][y])
			}
		}

		fmt.Println()
	}
}

func possibleIntersects(seta, setb Possible) bool {
	for i := 0; i < rowSize; i++ {
		if seta.possibilities[i] {
			if setb.possibilities[i] {
				return true
			}
		}
	}
	return false
}

func possibleIntersect(modify *Possible, reference Possible) {
	for i := 0; i < rowSize; i++ {
		if !reference.possibilities[i] {
			if (*modify).possibilities[i] {
				(*modify).possibilities[i] = false
				(*modify).count--
			}
		}
	}
}

func gridRef(x, y int) string {
	ret := string(x + 65)
	return ret + string(y+1+48)
}

func Difficulty(score int) string {
	if score <= ScoreSimple {
		return "simple"
	} else if score <= ScoreEasy {
		return "easy"
	} else if score <= ScoreMedium {
		return "moderate"
	} else if score <= ScoreHard {
		return "hard"
	} else if score == ScoreImpossible {
		return "impossible"
	}

	return "diabolical"
}

func Score(difficulty string) int {
	if difficulty == "simple" {
		return ScoreSimple
	} else if difficulty == "easy" {
		return ScoreEasy
	} else if difficulty == "moderate" {
		return ScoreMedium
	} else if difficulty == "hard" {
		return ScoreHard
	} else if difficulty == "diabolical" {
		return ScoreImpossible
	}

	return -1
}

func printPossible(in Possible) int {
	printed := 0

	for i := 0; i < rowSize; i++ {
		if in.possibilities[i] {
			fmt.Print(i + 1)
			printed++
		}
	}

	return printed
}
