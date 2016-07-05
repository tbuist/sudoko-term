package libhelp

import(
	"fmt"
	"os/exec"
	"math/rand"
	"time"
	_"strconv"
	"bytes"
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
		fmt.Println(tmp)
		buffer.Reset()
	}
}

// not a valid board
func FillBoard_junk(board *[9][9]int) {
	rand.Seed(time.Now().UnixNano())
	for _, v := range *board {
		for j := range v {
			tmp := rand.Intn(9) + 1
			v[j] = tmp
		}
	}
}