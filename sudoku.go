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

	//libhelp.ResizeTerm("18", "40")

	// Initialize main game board and board var to ptr to 2d array
	_board := [9][9]int{}
	board := &_board

	libhelp.FillBoard_junk(board)
	libhelp.PrintBoard(board)
	
	if libhelp.CheckBoard_valid(board) {
		fmt.Println("valid board")
	}

	if libhelp.CheckBoard_complete(board) {
		fmt.Println("complete board")
	} else {
		fmt.Println("incomplete board")
	}

	//var input string
	//fmt.Print("Press enter to continue...")
	//fmt.Scanln(&input)
	
	//libhelp.ResizeTerm("24", "80")

}