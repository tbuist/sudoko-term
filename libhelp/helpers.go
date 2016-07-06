package libhelp

import(
	"fmt"
	"os/exec"
	//"math/rand"
	//"time"
	"strconv"
	"bytes"
	"strings"
)

type Command struct {
	cmd int // 0 = mark, 1 = erase, 2 = check, 3 = quit, 4 = invalid
	val int
	row int
	col int
}

func ResizeTerm(height string, width string) {
	cmd := exec.Command("resize", "-s", height, width)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Resizing terminal failed. Please resize it yourself.")
	}
}

func PrintBoard(board *[9][9]int) {
	var buffer bytes.Buffer
	tmp := "\033[4m                         \033[0m"
	fmt.Println(tmp)


	// for each roach
	for i := range *board {
		// first border
		buffer.WriteString("|")
		for j, v := range (*board)[i] {
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
func FillBoard_junk(board *[9][9]int) {
	//rand.Seed(time.Now().UnixNano())
	for i := range *board {
		for j := range (*board)[i] {
			//tmp := rand.Intn(9) + 1
			(*board)[i][j] = -1
		}
	}
}

// returns (true, _) if valid, returns (true, 9) if complete
func CheckRow(board *[9][9]int, row int) (bool, int) {
	if 0 > row || row > 8 {
		return false, -1
	}

	count := 0
	m := make(map[int]int)
	for i, v := range (*board)[row] {
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

func CheckCol(board *[9][9]int, col int) (bool, int) {
	if 0 > col || col > 8 {
		return false, -1
	}

	count := 0
	m := make(map[int]int)
	for i := 0; i < 9; i++ {
		_, exists := m[(*board)[i][col]]
		if (*board)[i][col] != -1 && !exists {
			m[i] = (*board)[i][col]
			count++
		}
		if (*board)[i][col] != -1 && exists {
			return false, count
		}
	}
	return true, count

}

func CheckBoard_valid(board *[9][9]int) bool {
	for i := 0; i < 9; i++ {
		row_val, _ := CheckRow(board, i)
		col_val, _ := CheckCol(board, i)
		
		if !row_val || !col_val {
			return false
		}
	}
	return true
}

func CheckBoard_complete(board *[9][9]int) bool {
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

	fmt.Print("Command: ")
	var input string
	fmt.Scanln(&input)
	args := strings.Split(input, " ")

	// mark command
	switch len(args) {
		case 4:
			val, err1 := strconv.Atoi(args[1])
			row, err2 := strconv.Atoi(args[2])
			col, err3 := strconv.Atoi(args[3])
			if args[0] == "m" && inRange(val, 1, 9) && inRange(row, 1, 9) && inRange(col, 1, 9) && err1 == nil && err2 == nil && err3 == nil {
				cmd.cmd = 0
				cmd.val = val
				cmd.row = row
				cmd.col = col
				return &cmd
			}
			fallthrough
		case 3:
			row, err1 := strconv.Atoi(args[1])
			col, err2 := strconv.Atoi(args[2])
			if args[0] == "e" && inRange(row, 1, 9) && inRange(col, 1, 9) && err1 == nil && err2 == nil {
				cmd.cmd = 1
				cmd.val = -1
				cmd.row = row
				cmd.col = col
				return &cmd
			}
			fallthrough
		case 1:
			switch args[0] {
				case "c":
					cmd.cmd = 2
					cmd.val = -1
					cmd.row = -1
					cmd.col = -1
					return &cmd
				case "h":
					cmd.cmd = 3
					cmd.val = -1
					cmd.row = -1
					cmd.col = -1
					return &cmd
				case "q":
					cmd.cmd = 4
					cmd.val = -1
					cmd.row = -1
					cmd.col = -1
					return &cmd
			}
			fallthrough
		default:
			fmt.Println("Invalid command\n")
			return ReadCommand()
	}
}

func PrintInstructions() {
	fmt.Println("Instructions:")
	fmt.Println("	1. Comands")
	fmt.Println("		a. \"m 3 4 5\" marks [4][5] as a 3")
	fmt.Println("		b. \"e 4 5\" erases [4][5]")
	fmt.Println("		c. \"c\" checks board for validity/complete")
	fmt.Println("		d. \"h\" prints this again")
	fmt.Println("		e. \"q\" quits")
}

// inclusive
func inRange(num int, left int, right int) bool {
	return num >= left && num <= right
}