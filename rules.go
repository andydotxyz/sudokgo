package sudokgo

import (
	"fmt"
)

var rules []func(*Sudoku) (bool, int)

func init() {
	rules = []func(*Sudoku) (bool, int){
		ruleLastPossibility,
		ruleOnlyPossiblePlaceInRow,
		ruleOnlyPossiblePlaceInCol,
		ruleOnlyPossiblePlaceInBox,
		ruleMustBeInCertainBox,
		ruleEliminateSubsetExtras,
		ruleXwingRow,
		ruleXwingCol,
	}
}

func ruleLastPossibility(s *Sudoku) (bool, int) {
	for y := 0; y < RowSize; y++ {
		for x := 0; x < RowSize; x++ {
			if s.possGrid[x][y].count == 1 {
				for z := 0; z < RowSize; z++ {
					if s.possGrid[x][y].possibilities[z] {
						s.Grid[x][y] = z + 1
						s.possGrid[x][y].possibilities[z] = false
						s.possGrid[x][y].count = 0

						fmt.Println(gridRef(x, y), z+1, "is the last possibility")
						s.updatePossibilities(x, y)
						return true, RuleSimple
					}
				}
			}
		}
	}

	return false, 0
}

func ruleOnlyPossiblePlaceInRow(s *Sudoku) (bool, int) {
	for z := 0; z < RowSize; z++ {
		for y := 0; y < RowSize; y++ {
			found := false
			possCount := 0
			lastPlace := 0

			for x := 0; x < RowSize; x++ {
				if s.Grid[x][y] == z+1 {
					found = true
					break
				} else if s.possGrid[x][y].possibilities[z] {
					possCount++
					lastPlace = x
				}
			}

			if !found && possCount == 1 {
				s.Grid[lastPlace][y] = z + 1
				for a := 0; a < RowSize; a++ {
					s.possGrid[lastPlace][y].possibilities[a] = false
				}
				s.possGrid[lastPlace][y].count = 0
				fmt.Println(gridRef(lastPlace, y), "must be", z+1, "it is the only possible place in the row")

				s.updatePossibilities(lastPlace, y)
				return true, RuleEasy
			}
		}
	}

	return false, 0
}

func ruleOnlyPossiblePlaceInCol(s *Sudoku) (bool, int) {
	for z := 0; z < RowSize; z++ {
		for x := 0; x < RowSize; x++ {
			found := false
			possCount := 0
			lastPlace := 0

			for y := 0; y < RowSize; y++ {
				if s.Grid[x][y] == z+1 {
					found = true
					break
				} else if s.possGrid[x][y].possibilities[z] {
					possCount++
					lastPlace = y
				}
			}

			if !found && possCount == 1 {
				s.Grid[x][lastPlace] = z + 1
				for a := 0; a < RowSize; a++ {
					s.possGrid[x][lastPlace].possibilities[a] = false
				}
				s.possGrid[x][lastPlace].count = 0
				fmt.Println(gridRef(x, lastPlace), "must be", z+1, "it is the only possible place in the column")
				s.updatePossibilities(x, lastPlace)
				return true, RuleEasy
			}
		}
	}

	return false, 0
}

func ruleOnlyPossiblePlaceInBox(s *Sudoku) (bool, int) {
	for z := 0; z < RowSize; z++ {
		for y := 0; y < RowSize; y += GridSize {
			for x := 0; x < RowSize; x += GridSize {
				found := false
				possCount := 0
				lastPlaceX, lastPlaceY := 0, 0

				for xx := x; xx < x+GridSize; xx++ {
					for yy := y; yy < y+GridSize; yy++ {
						if s.Grid[xx][yy] == z+1 {
							found = true
							break
						} else if s.possGrid[xx][yy].possibilities[z] {
							possCount++
							lastPlaceX = xx
							lastPlaceY = yy
						}
					}
				}

				if !found && possCount == 1 {
					s.Grid[lastPlaceX][lastPlaceY] = z + 1
					for a := 0; a < RowSize; a++ {
						s.possGrid[lastPlaceX][lastPlaceY].possibilities[a] = false
					}
					s.possGrid[lastPlaceX][lastPlaceY].count = 0
					fmt.Println(gridRef(lastPlaceX, lastPlaceY), "must be", z+1, "it is the only possible place in the box")
					s.updatePossibilities(lastPlaceX, lastPlaceY)
					return true, RuleEasy
				}
			}
		}
	}

	return false, 0
}

func ruleMustBeInCertainBox(s *Sudoku) (bool, int) {
	for z := 0; z < RowSize; z++ {
		for y := 0; y < RowSize; y += GridSize {
			for x := 0; x < RowSize; x += GridSize {
				row, col := -1, -1
				breakingRow, breakingCol := false, false

				for xx := x; xx < x+GridSize && (!breakingRow || !breakingCol); xx++ {
					for yy := y; yy < y+GridSize && (!breakingRow || !breakingCol); yy++ {
						if s.possGrid[xx][yy].possibilities[z] {
							if row == -1 {
								row = yy
							} else if row != yy {
								breakingRow = true
							}

							if col == -1 {
								col = xx
							} else if col != xx {
								breakingCol = true
							}
						}
					}
				}

				if !breakingRow && row != -1 {
					for xx := 0; xx < RowSize; xx++ {
						if xx < x || xx >= x+GridSize {
							/* if not in that box, remove */
							if s.possGrid[xx][row].possibilities[z] {
								fmt.Println(gridRef(xx, row), "cannot be", z+1, "as it must appear elsewhere in the row")
								s.possGrid[xx][row].possibilities[z] = false
								s.possGrid[xx][row].count--
								return true, RuleMedium
							}
						}
					}
				}
				if !breakingCol && col != -1 {
					for yy := 0; yy < RowSize; yy++ {
						if yy < y || yy >= y+GridSize {
							/* if not in that box, remove */
							if s.possGrid[col][yy].possibilities[z] {
								fmt.Println(gridRef(col, yy), "cannot be", z+1, "as it must appear elsewhere in the column")
								s.possGrid[col][yy].possibilities[z] = false
								s.possGrid[col][yy].count--
								return true, RuleMedium
							}
						}
					}
				}
			}
		}
	}

	return false, 0
}

func ruleEliminateSubsetExtrasSlave(s *Sudoku, set Possible, x, y int) (bool, int) {
	ret := false
	score := 0

	if set.count <= 2 {
		return false, 0
	}

	for ii := 0; ii < RowSize; ii++ {
		if set.possibilities[ii] {
			subset := set
			subset.possibilities[ii] = false
			subset.count--

			subsetstr := ""
			for xx := 0; xx < RowSize; xx++ {
				if subset.possibilities[xx] {
					subsetstr += string(xx + 1 + 48)
				}
			}

			applied, points := ruleEliminateSubsetExtrasSlave(s, subset, x, y)
			ret = ret || applied
			score += points

			matches := 0
			for xx := 0; xx < RowSize; xx++ {
				if possibleIntersects(subset, s.possGrid[xx][y]) {
					matches++
				}
			}
			if matches == subset.count {
				for xx := 0; xx < RowSize; xx++ {
					if subset.count < s.possGrid[xx][y].count &&
						possibleIntersects(subset, s.possGrid[xx][y]) {
						possibleIntersect(&s.possGrid[xx][y], subset)
						fmt.Println(gridRef(xx, y), "subset", subsetstr, "cycle on row, eliminate extras")
						score += RuleHard
					}
				}
				ret = true
			}

			matches = 0
			for yy := 0; yy < RowSize; yy++ {
				if possibleIntersects(subset, s.possGrid[x][yy]) {
					matches++
				}
			}
			if matches == subset.count {
				for yy := 0; yy < RowSize; yy++ {
					if subset.count < s.possGrid[x][yy].count &&
						possibleIntersects(subset, s.possGrid[x][yy]) {
						possibleIntersect(&s.possGrid[x][yy], subset)
						fmt.Println(gridRef(x, yy), "subset", subsetstr, "cycle on col, eliminate extras")
						score += ScoreHard
					}
				}
				ret = true
			}

			matches = 0
			startx := x - (x % GridSize)
			starty := y - (y % GridSize)
			for xx := startx; xx < startx+GridSize; xx++ {
				for yy := starty; yy < starty+GridSize; yy++ {
					if possibleIntersects(subset, s.possGrid[xx][yy]) {
						matches++
					}
				}
			}

			if matches == subset.count {
				for xx := startx; xx < startx+GridSize; xx++ {
					for yy := starty; yy < starty+GridSize; yy++ {
						if subset.count < s.possGrid[xx][yy].count &&
							possibleIntersects(subset, s.possGrid[xx][yy]) {
							possibleIntersect(&s.possGrid[xx][yy], subset)
							fmt.Println(gridRef(xx, yy), "subset", subsetstr, "cycle on box, eliminate extras")
							score += RuleHard
						}
					}
				}
				ret = true
			}
		}
	}
	return ret, score
}

func ruleEliminateSubsetExtras(s *Sudoku) (bool, int) {
	ret := false
	score := 0

	for y := 0; y < RowSize; y++ {
		for x := 0; x < RowSize; x++ {
			if s.possGrid[x][y].count > 1 {
				/* if it is not the largest possible set in the area */
				applied, points := ruleEliminateSubsetExtrasSlave(s, s.possGrid[x][y], x, y)

				ret = ret || applied
				score += points
			}
		}
	}

	return ret, score
}

func ruleXwingRow(s *Sudoku) (bool, int) {
	for z := 0; z < RowSize; z++ {
		foundX1, foundX2 := 0, 0
		foundY1, foundY2 := 0, 0
		yCount := 0

		for y := 0; y < RowSize; y++ {
			xCount := 0
			x1, x2 := 0, 0

			for x := 0; x < RowSize; x++ {
				if s.possGrid[x][y].possibilities[z] {
					xCount++

					if yCount != 0 {
						if x != foundX1 && x != foundX2 {
							xCount = 0
							break
						}
					} else {
						if x1 != 0 {
							x2 = x
						} else {
							x1 = x
						}
					}
				}
			}

			if xCount == 2 {
				if yCount == 0 {
					foundX1 = x1
					foundX2 = x2
				}
				yCount++

				if yCount > 2 {
					yCount = 0
					break
				}

				if foundY1 != 0 {
					foundY2 = y
				} else {
					foundY1 = y
				}
			}
		}

		if yCount == 2 {
			found := false
			for yy := 0; yy < RowSize; yy++ {
				if yy == foundY1 || yy == foundY2 {
					continue
				}

				if s.possGrid[foundX1][yy].possibilities[z] {
					s.possGrid[foundX1][yy].possibilities[z] = false
					s.possGrid[foundX1][yy].count--
					found = true
				}

				if s.possGrid[foundX2][yy].possibilities[z] {
					s.possGrid[foundX2][yy].possibilities[z] = false
					s.possGrid[foundX2][yy].count--
					found = true
				}
			}

			if found {
				fmt.Println("XWING rows found with diagonals", gridRef(foundX1, foundY1), "and",
					gridRef(foundX2, foundY2), "- eliminating", z+1, "on cols")

				return true, RuleDiabolical
			}
		}
	}

	return false, 0
}

func ruleXwingCol(s *Sudoku) (bool, int) {
	for z := 0; z < RowSize; z++ {
		foundX1, foundX2 := 0, 0
		foundY1, foundY2 := 0, 0
		xCount := 0

		for x := 0; x < RowSize; x++ {
			yCount := 0
			y1, y2 := 0, 0

			for y := 0; y < RowSize; y++ {
				if s.possGrid[x][y].possibilities[z] {
					yCount++

					if xCount != 0 {
						if y != foundY1 && y != foundY2 {
							yCount = 0
							break
						}
					} else {
						if y1 != 0 {
							y2 = y
						} else {
							y1 = y
						}
					}
				}
			}

			if yCount == 2 {
				if xCount != 0 {
					foundY1 = y1
					foundY2 = y2
				}
				xCount++

				if xCount > 2 {
					xCount = 0
					break
				}

				if foundX1 != 0 {
					foundX2 = x
				} else {
					foundX1 = x
				}
			}
		}

		if xCount == 2 {
			found := false
			for xx := 0; xx < RowSize; xx++ {
				if xx == foundX1 || xx == foundX2 {
					continue
				}

				if s.possGrid[xx][foundY1].possibilities[z] {
					s.possGrid[xx][foundY1].possibilities[z] = false
					s.possGrid[xx][foundY1].count--
					found = true
				}

				if s.possGrid[xx][foundY2].possibilities[z] {
					s.possGrid[xx][foundY2].possibilities[z] = false
					s.possGrid[xx][foundY2].count--
					found = true
				}
			}

			if found {
				fmt.Println("XWING cols found with diagonals", gridRef(foundX1, foundY1), "and",
					gridRef(foundX2, foundY2), "- eliminating", z+1, "on rows")
				return true, RuleDiabolical
			}
		}
	}

	return false, 0
}
