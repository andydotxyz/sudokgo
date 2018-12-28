// Package main loads the sudokgo solver / generator
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/andydotxyz/sudokgo"
)

func printUsage() {
	fmt.Println("Usage: sudokgo {generate|solve} [parameters]")
	fmt.Println("  generate takes no params and generates a sudoku puzzle")
	fmt.Println("  solve takes 1 parameter, the puzzle to solve, using numbers or dash 12-4-...")
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	game := sudokgo.NewSudoku()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
	} else if args[0] == "generate" {
		score := sudokgo.ScoreMedium

		if len(args) >= 2 {
			newScore := sudokgo.Score(args[1])

			if newScore != -1 {
				score = newScore
			} else {
				fmt.Println("invalid score, please try simple, easy, moderate, hard or diabolical")
				fmt.Println("using moderate")
				score = sudokgo.ScoreMedium
			}
		} else {
			fmt.Println("Generating a moderate puzzle")
		}

		fmt.Println("Score,", score)
		game.Generate(score)
		game.Print()
	} else if args[0] == "solve" {
		if len(args) < 2 {
			fmt.Println("solve command needs an extra option, the grid to solve (e.g. 12-3-4 etc)")
			os.Exit(1)
		}

		if game.Load(args[1]) != nil {
			os.Exit(1)
		}

		game.Print()

		start := time.Now()
		score, err := game.Solve()
		micro := time.Since(start).Round(time.Microsecond)

		fmt.Println("Finished in ", micro)
		if err == nil {
			fmt.Printf("Difficulty score was %d (%s).\n", score, sudokgo.Difficulty(score))
			game.Print()
		}
	}
}
