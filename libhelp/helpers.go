package libhelp

import(
	"fmt"
	"os/exec"
	//"math/rand"
	//"time"
	_"strconv"
	"bytes"
	"strings"
)

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

func CheckRow_valid(board *[9][9]int, row int) (bool, int) {
	count := 0
	m := make(map[int]int)
	for _, v := range (*board)[row] {
		tmp, exists := m[v]
		if v != -1 && exists {
			return false
		}
	}
	return true
}

func CheckRow_complete(board *[9][9]int, row int) bool {
	
}

func CheckCol_valid(board *[9][9]int, col int) bool {
	return true
}

func CheckCol_complete(board *[9][9]int, col int) bool {
	return true
}

func CheckBoard_valid(board *[9][9]int) bool {
	for i := 0; i < 9; i++ {
		if !CheckRow_valid(board, i) || !CheckCol_valid(board, i) {
			return false
		}
	}
	return true
}

func CheckBoard_complete(board *[9][9]int) bool {
	for i := 0; i < 9; i++ {
		if !CheckRow_complete(board, i) || !CheckCol_complete(board, i) {
			return false
		}
	}
	return true
}