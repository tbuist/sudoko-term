/* Terminal sudoku written in Go
 *
 * Author:
 *		Taylor Buist - tbuist@umich.edu
 * Latest revision:
 *		May 11, 2017
 */

package main

import (
	"fmt"
	"os/exec"
	//"time"
	"bufio"
	"bytes"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Command int

const (
	Mark Command = iota + 1
	Erase
	Check
	Generate
	Help
	Quit
	Invalid
)

// 1-indexed?
type Input struct {
	Cmd Command // 0 = mark, 1 = erase, 2 = check, 3 = quit, 4 = invalid
	Val int
	Row int
	Col int
}

type Board struct {
	Arr  [9][9]int
	Done bool
}

func Board_default() *Board {
	var board Board
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			board.Arr[i][j] = -1
		}
	}
	board.Done = false
	return &board
}

func ResizeTerm(height string, width string) {
	cmd := exec.Command("resize", "-s", height, width)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Resizing terminal failed. Please resize it yourself.")
	}
}

func ClearTerm() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func PrintBoard(board *Board) {
	var buffer bytes.Buffer
	tmp := "\033[4m                         \033[0m"
	fmt.Println(tmp)

	// for each roach
	for i := range (*board).Arr {
		// first border
		buffer.WriteString("|")
		for j, v := range (*board).Arr[i] {
			if (j+1)%3 == 0 {
				buffer.WriteString(fmt.Sprintf(" %v |", v))
			} else {
				buffer.WriteString(fmt.Sprintf(" %v", v))
			}
		}

		tmp = buffer.String()
		// if 3rd, 6th, or 9th row, underline
		if (i+1)%3 == 0 {
			tmp = fmt.Sprintf("\033[4m%s\033[0m", tmp)
		}
		tmp = strings.Replace(tmp, "-1", "-", -1)
		fmt.Println(tmp)
		buffer.Reset()
	}
}

// not a valid board
func FillBoard_junk(board *Board) {
	//rand.Seed(time.Now().UnixNano())
	for i := range (*board).Arr {
		for j := range (*board).Arr[i] {
			//tmp := rand.Intn(9) + 1
			(*board).Arr[i][j] = -1
		}
	}
}

// returns (true, _) if valid, returns (true, 9) if complete
func CheckRow(board *Board, row int) (bool, bool) {
	if 0 > row || row > 8 {
		return false, false
	}

	count := 0
	m := [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	for _, v := range board.Arr[row] {
		switch {
		case inRange(v, 0, 8):
			if m[v] > 0 {
				return false, false
			}
			m[v]++
			count++
		case v == -1:
		default:
			return false, false
		}
	}
	return true, count == 9
}

// returns (true, _) if valid, returns (true, 9) if complete
func CheckCol(board *Board, col int) (bool, bool) {
	if 0 > col || col > 8 {
		return false, false
	}

	count := 0
	m := [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}

	for i := 0; i < 9; i++ {
		tmp := board.Arr[i][col]
		switch {
		case inRange(tmp, 0, 8):
			if m[tmp] > 0 {
				return false, false
			}
			m[tmp]++
			count++
		case tmp == -1:
		default:
			return false, false
		}
	}
	return true, count == 9
}

func CheckSquare(board *Board, row int, col int) (bool, bool) {
	if !inRange(row, 0, 8) || !inRange(col, 0, 8) {
		return false, false
	}

	count := 0
	m := [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}

	startRow := row - (row % 3)
	startCol := col - (col % 3)
	for i := startRow; i < startRow+3; i++ {
		for j := startCol; j < startCol+3; j++ {
			tmp := board.Arr[i][j]
			switch {
			case inRange(tmp, 0, 8):
				if m[tmp] > 0 {
					return false, false
				}
				m[tmp]++
				count++
			case tmp == -1:
			default:
				return false, false
			}
		}
	}
	return true, count == 9
}

func CheckBoard_valid(board *Board) bool {
	for i := 0; i < 9; i++ {
		rowVal, _ := CheckRow(board, i)
		colVal, _ := CheckCol(board, i)
		sqVal, _ := CheckSquare(board, i, i)

		if !rowVal || !colVal || !sqVal {
			return false
		}
	}
	return true
}

func CheckBoard_complete(board *Board) bool {
	for i := 0; i < 9; i++ {
		rowVal, rowComp := CheckRow(board, i)
		colVal, colComp := CheckCol(board, i)
		sqVal, sqComp := CheckSquare(board, i, i)

		if !rowVal || !colVal || !rowComp || !colComp || !sqVal || !sqComp {
			return false
		}
	}
	return true
}

func generateBoard(board *Board, diff int) {
	FillBoard_junk(board)
	numValsFilled := 0
	toFill := 15 + 6*(5-diff)
	for numValsFilled < toFill {
		randRow := rand.Int() % 9
		randCol := rand.Int() % 9
		randVal := rand.Int()%9 + 1
		tmp := board.Arr[randRow][randCol]
		if tmp == -1 {
			board.Arr[randRow][randCol] = randVal
			if CheckBoard_valid(board) {
				numValsFilled++
			} else {
				board.Arr[randRow][randCol] = -1
			}
		}
	}
}

func ReadCommand() *Input {
	var cmd Input

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Command: ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	args := strings.Split(text, " ")

	// mark command
	switch len(args) {
	case 4:
		fmt.Println(args)
		val, err1 := strconv.Atoi(args[1])
		row, err2 := strconv.Atoi(args[2])
		col, err3 := strconv.Atoi(args[3])
		fmt.Printf("Mark command: value %d at position (%d,%d)\n", val, row, col)
		if args[0] == "m" && inRange(val, 1, 9) && inRange(row, 1, 9) && inRange(col, 1, 9) && err1 == nil && err2 == nil && err3 == nil {
			cmd.Cmd = Mark
			cmd.Val = val
			cmd.Row = row
			cmd.Col = col
			return &cmd
		}
		fmt.Println("error")
		os.Exit(0)
		fallthrough
	case 3:
		row, err1 := strconv.Atoi(args[1])
		col, err2 := strconv.Atoi(args[2])
		fmt.Printf("Erase command: position (%d,%d)\n", row, col)
		if args[0] == "e" && inRange(row, 1, 9) && inRange(col, 1, 9) && err1 == nil && err2 == nil {
			cmd.Cmd = Erase
			cmd.Val = -1
			cmd.Row = row
			cmd.Col = col
			return &cmd
		}
		fallthrough
	case 2:
		// generate
		cmd.Cmd = Generate

		val, err := strconv.Atoi(args[1])
		if err == nil && inRange(val, 1, 5) && args[0] == "g" {
			cmd.Val = val
			return &cmd
		}
		fallthrough
	case 1:
		switch args[0] {
		case "c":
			cmd.Cmd = Check
			cmd.Val = -1
			cmd.Row = -1
			cmd.Col = -1
			return &cmd
		case "h":
			cmd.Cmd = Help
			cmd.Val = -1
			cmd.Row = -1
			cmd.Col = -1
			return &cmd
		case "q":
			cmd.Cmd = Quit
			cmd.Val = -1
			cmd.Row = -1
			cmd.Col = -1
			return &cmd
		}
		fallthrough
	default:
		fmt.Print("Invalid command\n")
		return nil
	}
}

func PrintInstructions() {
	fmt.Println("Instructions:")
	fmt.Println("	1. Commands")
	fmt.Println("		a. \"m <val> <row> <col>\" marks [<row>][<col>] as a <val>")
	fmt.Println("		b. \"e <row> <col>\" erases [<row>][<col>]")
	fmt.Println("		c. \"c\" checks board for validity/complete")
	fmt.Println("		d. \"g <dif>\" generates a new board with difficulty <dif> (1-5) easy to hard")
	fmt.Println("		e. \"h\" prints this again")
	fmt.Println("		f. \"q\" quits")
}

// inclusive
func inRange(num int, left int, right int) bool {
	return num >= left && num <= right
}

func main() {

	fmt.Println("Welcome to Terminal Sudoku. Let's begin")
	PrintInstructions()

	//libhelp.ResizeTerm("12", "40")

	// Initialize main game board and board var to ptr to 2d array
	board := Board_default()

	for !board.Done {
		PrintBoard(board)
		cmd := ReadCommand()
		ClearTerm()
		for cmd == nil {
			cmd = ReadCommand()
			ClearTerm()
		}

		switch cmd.Cmd {
		case Mark:
			tmp := board.Arr[cmd.Row-1][cmd.Col-1]
			board.Arr[cmd.Row-1][cmd.Col-1] = cmd.Val
			rowValid, _ := CheckRow(board, cmd.Row-1)
			colValid, _ := CheckCol(board, cmd.Col-1)
			if !rowValid || !colValid {
				fmt.Println("Invalid choice")
				board.Arr[cmd.Row-1][cmd.Col-1] = tmp
			}
		case Erase:
			board.Arr[cmd.Row-1][cmd.Col-1] = -1
		case Check:
			valid := CheckBoard_valid(board)
			complete := CheckBoard_complete(board)

			if valid {
				if complete {
					fmt.Println("Board is valid and complete")
					(*board).Done = true
				}
				fmt.Println("Board is valid but not complete")
			} else {
				fmt.Println("Board is not valid")
			}
		case Generate:
			fmt.Printf("Generating new board with difficulty %d\n", cmd.Val)
			generateBoard(board, cmd.Val)
		case Help:
			PrintInstructions()
		case Quit:
			fmt.Println("Goodbye")
			(*board).Done = true
		default:

		}
	}
}
