package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(board []string) int {
	rows := 0
	for i := 0; i < 8; i++ {
		if !strings.ContainsRune(board[i], 'W') {
			rows++
		}
	}
	if rows == 8 {
		return 8
	}
	cols := 0
	for c := 0; c < 8; c++ {
		all := true
		for r := 0; r < 8 && all; r++ {
			if board[r][c] != 'B' {
				all = false
			}
		}
		if all {
			cols++
		}
	}
	return rows + cols
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
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
		if len(parts) != 8 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		board := parts
		exp := expected(board)
		var input bytes.Buffer
		for i := 0; i < 8; i++ {
			input.WriteString(board[i])
			input.WriteByte('\n')
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outStr := strings.TrimSpace(string(out))
		var got int
		fmt.Sscan(outStr, &got)
		if got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, outStr)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
