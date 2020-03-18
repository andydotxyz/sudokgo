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
	fmt.Println("Usage: sudokgo [-h] {generate|solve} [parameters]")
	fmt.Println("  generate takes an optional difficulty parameter (default is moderate)")
	fmt.Println("  solve takes 1 parameter, the puzzle to solve, using numbers or dash 12-4-...")
	fmt.Println("  -h prints human readable grids instead of a number list on 1 line")
}

func main() {
	flag.Usage = printUsage
	pretty := flag.Bool("h", false, "Print human readable grids instead of a number list on 1 line")
	verbose := flag.Bool("v", false, "Print verbose information and hints")
	flag.Parse()

	game := sudokgo.NewSudoku()
	game.Verbose = *verbose
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

		game.Generate(score)
		printStringOrGrid(game, *pretty)
	} else if args[0] == "solve" {
		if len(args) < 2 {
			fmt.Println("solve command needs an extra option, the grid to solve (e.g. 12-3-4 etc)")
			os.Exit(1)
		}

		if err := game.Load(args[1]); err != nil {
			fmt.Println("Error during load:", err)
			os.Exit(1)
		}

		if *verbose {
			printStringOrGrid(game, *pretty)
		}

		start := time.Now()
		score, err := game.Solve()
		micro := time.Since(start).Round(time.Microsecond)

		if err != nil {
			fmt.Println("Error during solve:", err)
		} else {
			if *verbose {
				fmt.Println("Finished in ", micro)
				fmt.Printf("Difficulty score was %d (%s).\n", score, sudokgo.Difficulty(score))
			}
			printStringOrGrid(game, *pretty)
		}
	}
}

func printStringOrGrid(s *sudokgo.Sudoku, pretty bool) {
	if pretty {
		s.Print()
	} else {
		fmt.Println(s.String())
	}
}
