package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expected(board [][]int) int {
	n := len(board)
	rowSum := make([]int, n)
	colSum := make([]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := board[i][j]
			rowSum[i] += v
			colSum[j] += v
		}
	}
	cnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if colSum[j] > rowSum[i] {
				cnt++
			}
		}
	}
	return cnt
}

func runCase(bin string, board [][]int) error {
	var input strings.Builder
	n := len(board)
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", board[i][j]))
		}
		input.WriteByte('\n')
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	var got int
	if _, err := fmt.Fscan(&buf, &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := expected(board)
	if got != exp {
		return fmt.Errorf("expected %d got %d\ninput:\n%s", exp, got, input.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(30) + 1
		board := make([][]int, n)
		for r := 0; r < n; r++ {
			board[r] = make([]int, n)
			for c := 0; c < n; c++ {
				board[r][c] = rng.Intn(100) + 1
			}
		}
		if err := runCase(bin, board); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
