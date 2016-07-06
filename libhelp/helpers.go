package libhelp

import(
	"fmt"
	"os/exec"
	//"math/rand"
	//"time"
	"strconv"
	"bytes"
	"strings"
	"bufio"
	"os"
)


// 1-indexed?
type Command struct {
	Cmd int // 0 = mark, 1 = erase, 2 = check, 3 = quit, 4 = invalid
	Val int
	Row int
	Col int
}

type Board struct {
	Arr [9][9]int
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

func PrintBoard(board *Board) {
	var buffer bytes.Buffer
	tmp := "\033[4m                         \033[0m"
	fmt.Println(tmp)


	// for each roach
	for i := range (*board).Arr {
		// first border
		buffer.WriteString("|")
		for j, v := range (*board).Arr[i] {
			if (j+1) % 3 == 0 {
				buffer.WriteString(fmt.Sprintf(" %v |", v))
			} else {
				buffer.WriteString(fmt.Sprintf(" %v", v))
			}
		}

		tmp = buffer.String()
		// if 3rd, 6th, or 9th row, underline
		if (i+1) % 3 == 0 {
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
func CheckRow(board *Board, row int) (bool, int) {
	if 0 > row || row > 8 {
		return false, -1
	}

	count := 0
	m := make(map[int]int)
	for i, v := range (*board).Arr[row] {
		_, exists := m[v]
		if v != -1 && !exists {
			m[i] = v
			count++
		}
		if v != -1 && exists {
			return false, count
		}
	}
	return true, count
}

// returns (true, _) if valid, returns (true, 9) if complete
func CheckCol(board *Board, col int) (bool, int) {
	if 0 > col || col > 8 {
		return false, -1
	}

	count := 0
	arr := [9]int{-1,-1,-1,-1,-1,-1,-1,-1,-1}

	for i := 0; i < 9; i++ {
		tmp := (*board).Arr[i][col]

		if tmp != -1 {
			if arr[tmp-1] == -1 {
				arr[tmp-1] = tmp
				count++
			}
			if arr[tmp-1] > -1 {
				return false, count
			}	
		}
	}
	return true, count
}

func CheckSquare(board *Board, row int, col int) (bool, int) {
	
	return true, -1
}

func CheckBoard_valid(board *Board) bool {
	for i := 0; i < 9; i++ {
		row_val, _ := CheckRow(board, i)
		col_val, _ := CheckCol(board, i)
		
		if !row_val || !col_val {
			return false
		}
	}
	return true
}

func CheckBoard_complete(board *Board) bool {
	for i := 0; i < 9; i++ {
		row_val, row_count := CheckRow(board, i)
		col_val, col_count := CheckCol(board, i)
		
		if !row_val || !col_val || row_count != 9 || col_count != 9 {
			return false
		}
	}
	return true
}

func ReadCommand() *Command {
	var cmd Command

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Command: ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	args := strings.Split(text, " ")

	// mark command
	switch len(args) {
		case 4:
			val, err1 := strconv.Atoi(args[1])
			row, err2 := strconv.Atoi(args[2])
			col, err3 := strconv.Atoi(args[3])
			if args[0] == "m" && inRange(val, 1, 9) && inRange(row, 1, 9) && inRange(col, 1, 9) && err1 == nil && err2 == nil && err3 == nil {
				cmd.Cmd = 0
				cmd.Val = val
				cmd.Row = row
				cmd.Col = col
				return &cmd
			}
			fallthrough
		case 3:
			row, err1 := strconv.Atoi(args[1])
			col, err2 := strconv.Atoi(args[2])
			if args[0] == "e" && inRange(row, 1, 9) && inRange(col, 1, 9) && err1 == nil && err2 == nil {
				cmd.Cmd = 1
				cmd.Val = -1
				cmd.Row = row
				cmd.Col = col
				return &cmd
			}
			fallthrough
		case 1:
			switch args[0] {
				case "c":
					cmd.Cmd = 2
					cmd.Val = -1
					cmd.Row = -1
					cmd.Col = -1
					return &cmd
				case "h":
					cmd.Cmd = 3
					cmd.Val = -1
					cmd.Row = -1
					cmd.Col = -1
					return &cmd
				case "q":
					cmd.Cmd = 4
					cmd.Val = -1
					cmd.Row = -1
					cmd.Col = -1
					return &cmd
			}
			fallthrough
		default:
			fmt.Print("Invalid command")
			return nil
	}
}

func PrintInstructions() {
	fmt.Println("Instructions:")
	fmt.Println("	1. Commands")
	fmt.Println("		a. \"m <val> <row> <col>\" marks [<row>][<col>] as a <val>")
	fmt.Println("		b. \"e <row> <col>\" erases [<row>][<col>]")
	fmt.Println("		c. \"c\" checks board for validity/complete")
	fmt.Println("		d. \"h\" prints this again")
	fmt.Println("		e. \"q\" quits")
}

// inclusive
func inRange(num int, left int, right int) bool {
	return num >= left && num <= right
}