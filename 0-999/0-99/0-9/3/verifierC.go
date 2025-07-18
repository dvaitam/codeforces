package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func win(board [3]string, ch byte) bool {
	for i := 0; i < 3; i++ {
		if board[i][0] == ch && board[i][1] == ch && board[i][2] == ch {
			return true
		}
		if board[0][i] == ch && board[1][i] == ch && board[2][i] == ch {
			return true
		}
	}
	if board[0][0] == ch && board[1][1] == ch && board[2][2] == ch {
		return true
	}
	if board[0][2] == ch && board[1][1] == ch && board[2][0] == ch {
		return true
	}
	return false
}

func verdict(board [3]string) string {
	xCnt, oCnt := 0, 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch board[i][j] {
			case 'X':
				xCnt++
			case '0':
				oCnt++
			}
		}
	}
	if oCnt > xCnt || xCnt-oCnt > 1 {
		return "illegal"
	}
	xWin := win(board, 'X')
	oWin := win(board, '0')
	if xWin && oWin {
		return "illegal"
	}
	if xWin && xCnt != oCnt+1 {
		return "illegal"
	}
	if oWin && xCnt != oCnt {
		return "illegal"
	}
	if xWin {
		return "the first player won"
	}
	if oWin {
		return "the second player won"
	}
	if xCnt+oCnt == 9 {
		return "draw"
	}
	if xCnt == oCnt {
		return "first"
	}
	return "second"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Printf("invalid testcase format on line %d\n", idx)
			os.Exit(1)
		}
		board := [3]string{parts[0], parts[1], parts[2]}
		exp := verdict(board)
		input := fmt.Sprintf("%s\n%s\n%s\n", parts[0], parts[1], parts[2])
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != exp {
			fmt.Printf("Test %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
