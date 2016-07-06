/* Terminal sudoku written in Golang
 *
 * Author:
 *		Taylor Buist - tbuist@umich.edu
 * Latest revision:
 *		July 05, 2016
 */

package main

import (
 	"fmt"
 	"github.com/tbuist/sudoku/libhelp"
)



func main() {

	fmt.Println("Welcome to Terminal Sudoku. Let's begin")
	libhelp.PrintInstructions()

	libhelp.ResizeTerm("12", "40")

	// Initialize main game board and board var to ptr to 2d array
	board := libhelp.Board_default()

	for ; !(*board).Done; {
		libhelp.PrintBoard(board)
		cmd := libhelp.ReadCommand()
		fmt.Println()

		for ; cmd == nil; {
			libhelp.PrintBoard(board)
			cmd = libhelp.ReadCommand()
		}

		switch cmd.Cmd {
			case 0:
				// mark
				tmp := board.Arr[cmd.Row-1][cmd.Col-1]
				board.Arr[cmd.Row-1][cmd.Col-1] = cmd.Val
				row_good, _ := libhelp.CheckRow(board, cmd.Row-1)
				col_good, _ := libhelp.CheckCol(board, cmd.Col-1)
				if !row_good || !col_good {
					fmt.Println("Invalid choice")
					board.Arr[cmd.Row-1][cmd.Col-1] = tmp
				}
			case 1:
				// erase
				board.Arr[cmd.Row-1][cmd.Col-1] = -1
			case 2:
				// check
				valid := libhelp.CheckBoard_valid(board)
				complete := libhelp.CheckBoard_complete(board)
				
				//var input string
				//fmt.Scanln(&input)

				if valid {
					if complete {
						fmt.Println("Board is valid and complete")
						(*board).Done = true
					}
					fmt.Println("Board is valid but not complete")
				} else {
					fmt.Println("Board is not valid")
				}
			case 3:
				libhelp.PrintInstructions()
			case 4:
				// quit
				fmt.Println("Goodbye")
				(*board).Done = true
			default:

		}
	}

	libhelp.ResizeTerm("24", "80")

}