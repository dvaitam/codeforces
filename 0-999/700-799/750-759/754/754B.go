package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var board [4]string
	for i := 0; i < 4; i++ {
		if _, err := fmt.Fscan(in, &board[i]); err != nil {
			return
		}
	}

	check := func(a, b, c byte) bool {
		x := 0
		dot := 0
		if a == 'x' {
			x++
		} else if a == '.' {
			dot++
		}
		if b == 'x' {
			x++
		} else if b == '.' {
			dot++
		}
		if c == 'x' {
			x++
		} else if c == '.' {
			dot++
		}
		return x == 2 && dot == 1
	}

	yes := false
	// Horizontal
	for i := 0; i < 4 && !yes; i++ {
		for j := 0; j <= 1; j++ {
			if check(board[i][j], board[i][j+1], board[i][j+2]) {
				yes = true
				break
			}
		}
	}
	// Vertical
	for j := 0; j < 4 && !yes; j++ {
		for i := 0; i <= 1; i++ {
			if check(board[i][j], board[i+1][j], board[i+2][j]) {
				yes = true
				break
			}
		}
	}
	// Diagonal TL-BR
	for i := 0; i <= 1 && !yes; i++ {
		for j := 0; j <= 1; j++ {
			if check(board[i][j], board[i+1][j+1], board[i+2][j+2]) {
				yes = true
				break
			}
		}
	}
	// Diagonal TR-BL
	for i := 0; i <= 1 && !yes; i++ {
		for j := 2; j < 4; j++ {
			if check(board[i][j], board[i+1][j-1], board[i+2][j-2]) {
				yes = true
				break
			}
		}
	}

	if yes {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
}
