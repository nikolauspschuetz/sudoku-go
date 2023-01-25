package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
nice post about logging
https://www.honeybadger.io/blog/golang-logging/
can use os.Stout/os.Stderr as well as a file
*/

var puzzleFileDefault = 
	// "puzzles/ppg-20221205.txt"
	// "puzzles/rps-20221103.txt"
	// "puzzles/sudoku-com-evil-20221206T173500.txt"
	"puzzles/sudoku-com-evil-20221216T173500.txt"
	// "puzzles/sudoku-com-hard-20221206T153200.txt"
	// "puzzles/sudoku-com-hard-20221206T173200.txt"


var fileFlag string

func init() {
	flag.StringVar(&fileFlag, "file", puzzleFileDefault, "Puzzle file")
}


func Board(filename string) [9][9]int {
	fmt.Println("parsing board from file ", fileFlag)
	f, err := os.Open(fileFlag)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	reader := bufio.NewReader(f)

	var board [9][9]int

	i := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if (len(line) > 0) && (line[0] != '#') {
			a := strings.Split(string(line), ",")
			if len(a) != 9 {
				msg := fmt.Sprintf("Incorrect row %d length %d != 9", i, len(a))
				log.Fatal(msg)
			}
			for j, c := range a {
				v, err := strconv.Atoi(c)
				if err != nil {
					log.Fatal(err)
				}
				board[i][j] = v
			}
			i++
		}
	}
	if i != 9 {
		msg := fmt.Sprintf("Incorrect number of rows %d, expected 9", i)
		log.Fatal(msg)
	}
	return board
}

func notContains(array [9]int, n int) bool {
	var result bool = true
	// iterate over the array using for loop and break if values match
	for i := 0; i < len(array); i++ {
		// checking if the array contains the given value
		if array[i] == n {
			// changing the boolean variable
			result = false
			break
		}
	}
	return result
}

func Solutions(board *[9][9]int, r int, c int) []int {
	// fmt.Printf("Getting solutions for %d, %d\n", r, c)
	var row [9]int
	var col [9]int
	var box [9]int
	for i := 0; i < 9; i++ {
		row[i] = board[r][i]
		col[i] = board[i][c]
	}
	k := 0
	r3 := 3 * int(r/3)
	c3 := 3 * int(c/3)
	for i := r3; i < r3+3; i++ {
		for j := c3; j < c3+3; j++ {
			// fmt.Printf("%d, %d box[%d], board[%d][%d] = %d\n", r3, c3, k, i, j, board[i][j])
			box[k] = board[i][j]
			k++
		}
	}
	var sol []int
	for i := 1; i <= 9; i++ {
		if notContains(row, i) {
			if notContains(col, i) {
				if notContains(box, i) {
					sol = append(sol, i)
				}
			}
		}
	}
	// fmt.Printf("Got solutions for %d, %d: [", r, c)
	// for i := 0; i < len(sol); i++ {
	//     fmt.Printf("%d", sol[i])
	// }
	// fmt.Printf("]\n")
	return sol
}

func Solution(board *[9][9]int, r int, c int) (int, error) {
	sol := Solutions(board, r, c)
	n := len(sol)
	if n == 0 {
		// PrintBoard(*board)
		return 0, errors.New(fmt.Sprintf("No solutions for %d, %d", r, c))
	}
	if n == 1 {
		// fmt.Printf("Got solution for %d, %d: %d\n", r, c, sol[0])
		return sol[0], nil
	} else {
		// fmt.Printf("No solution for %d, %d\n", r, c)
		return 0, nil
	}
}

func PrintBoard(board [9][9]int) {
	for i := 0; i < 9; i++ {
		fmt.Println(board[i])
	}
}

type Cell struct {
	row     int
	col     int
	options []int
}

func (c Cell) String() string {
	solJson, _ := json.Marshal(c.options)
	return fmt.Sprintf("Cell{%d, %d, %s}", c.row, c.col, solJson)
}

func Options(board *[9][9]int, sorted bool) []Cell {
	var cells []Cell
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] < 1 {
				sol := Solutions(board, i, j)
				cells = append(cells, Cell{row: i, col: j, options: sol})
			}
		}
	}

	if sorted {
		sort.Slice(cells, func(i, j int) bool {
			return len(cells[i].options) < len(cells[j].options)
		})
	}
	return cells
}

func Backtrack(board *[9][9]int) (bool, error) {
	fmt.Printf("Backtrack with Board:\n")
	PrintBoard(*board)

	options := Options(board, true)

	// fmt.Printf("Options:\n")
	// for i, c := range options {
	//     fmt.Printf("%d:\t%s\n", i, c)
	// }

	o := options[0]
	for n, s := range o.options {
		fmt.Printf("Backtracking on %d, %d options[%d]: %d\n", o.row, o.col, n, s)
		PrintBoard(*board)
		var backup [9][9]int
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				backup[i][j] = int(board[i][j])
			}
		}
		board[o.row][o.col] = s
		solved := false
		var err error
		for !solved {
			solved, err = Solve(board)
			if err != nil {
				fmt.Printf("Reverting board due to %s\n", err)
				PrintBoard(*board)
				for i := 0; i < 9; i++ {
					for j := 0; j < 9; j++ {
						if board[i][j] != backup[i][j] {
							board[i][j] = backup[i][j]
						}
					}
				}
				fmt.Printf("Reverted board\n")
				PrintBoard(*board)
				break
			}
		}
		if solved {
			return solved, nil
		}
	}
	return false, errors.New("Not solved with backtracking")
}

func Solve(board *[9][9]int) (bool, error) {
	solved := true
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] < 1 {
				// fmt.Printf("Solving %d, %d (%d)\n", i, j, board[i][j])
				solved = false
				answer, err := Solution(board, i, j)
				if err != nil {
					return solved, err
				}
				if answer > 0 {
					fmt.Printf("Solved %d, %d = %d\n", i, j, answer)
					board[i][j] = answer
					return solved, nil
				}
			}
		}
	}
	if !solved {
		var err error
		solved, err = Backtrack(board)
		return solved, err
	} else {
		return solved, nil
	}
}

func main() {
	flag.Parse()
	board := Board(fileFlag)
	fmt.Println("Starting Board:")
	PrintBoard(board)
	fmt.Println("Solving...")
	solved, err := Solve(&board)
	if err != nil {
		log.Fatal(err)
	}
	for !solved {
		solved, err = Solve(&board)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Solved Board:")
	PrintBoard(board)
}
