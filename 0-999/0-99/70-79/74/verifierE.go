package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateBoard() []string {
	perm := rand.Perm(len(chars))
	board := make([]string, 6)
	for i := 0; i < 6; i++ {
		row := make([]byte, 6)
		for j := 0; j < 6; j++ {
			row[j] = chars[perm[i*6+j]]
		}
		board[i] = string(row)
	}
	return board
}

func boardToString(b []string) string {
	return strings.Join(b, "\n") + "\n"
}

func applyMove(board [][]byte, op string) bool {
	if len(op) != 2 {
		return false
	}
	dir := op[0]
	idx := int(op[1] - '1')
	if idx < 0 || idx >= 6 {
		return false
	}
	switch dir {
	case 'L':
		row := board[idx]
		first := row[0]
		copy(row, row[1:])
		row[5] = first
	case 'R':
		row := board[idx]
		last := row[5]
		copy(row[1:], row[:5])
		row[0] = last
	case 'U':
		first := board[0][idx]
		for i := 0; i < 5; i++ {
			board[i][idx] = board[i+1][idx]
		}
		board[5][idx] = first
	case 'D':
		last := board[5][idx]
		for i := 5; i > 0; i-- {
			board[i][idx] = board[i-1][idx]
		}
		board[0][idx] = last
	default:
		return false
	}
	return true
}

func checkSolution(initial []string, output string) bool {
	rdr := strings.NewReader(output)
	scanner := bufio.NewScanner(rdr)
	if !scanner.Scan() {
		return false
	}
	var steps int
	if _, err := fmt.Sscan(scanner.Text(), &steps); err != nil {
		return false
	}
	if steps > 10000 || steps < 0 {
		return false
	}
	board := make([][]byte, 6)
	for i := 0; i < 6; i++ {
		board[i] = []byte(initial[i])
	}
	for i := 0; i < steps; i++ {
		if !scanner.Scan() {
			return false
		}
		op := strings.TrimSpace(scanner.Text())
		if !applyMove(board, op) {
			return false
		}
	}
	target := make([][]byte, 6)
	for i := 0; i < 6; i++ {
		target[i] = []byte(chars[i*6 : i*6+6])
	}
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			if board[i][j] != target[i][j] {
				return false
			}
		}
	}
	return true
}

func generateTest() string {
	board := generateBoard()
	return boardToString(board)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	const tests = 100
	for t := 1; t <= tests; t++ {
		inp := generateTest()
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewBufferString(inp)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if !checkSolution(strings.Split(strings.TrimSpace(inp), "\n"), string(out)) {
			fmt.Printf("Test %d failed. Invalid solution.\nInput:\n%s\nOutput:\n%s\n", t, inp, string(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
