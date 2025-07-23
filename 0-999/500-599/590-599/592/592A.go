package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	board := make([]string, 8)
	for i := 0; i < 8; i++ {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		board[i] = line
	}

	minA := 8
	minB := 8

	for c := 0; c < 8; c++ {
		// check column for white pawn reaching top
		for r := 0; r < 8; r++ {
			ch := board[r][c]
			if ch == 'W' {
				if r < minA {
					minA = r
				}
				break
			} else if ch == 'B' {
				break
			}
		}
		// check column for black pawn reaching bottom
		for r := 7; r >= 0; r-- {
			ch := board[r][c]
			if ch == 'B' {
				steps := 7 - r
				if steps < minB {
					minB = steps
				}
				break
			} else if ch == 'W' {
				break
			}
		}
	}

	if minA <= minB {
		fmt.Println("A")
	} else {
		fmt.Println("B")
	}
}
