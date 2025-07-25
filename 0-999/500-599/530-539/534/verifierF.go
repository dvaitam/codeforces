package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func countSegments(line string) int {
	cnt := 0
	prev := '.'
	for i := 0; i < len(line); i++ {
		if line[i] == '*' && prev != '*' {
			cnt++
		}
		prev = line[i]
	}
	return cnt
}

func checkBoard(out string, n, m int, rowSeg []int, colSeg []int) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != n {
		return fmt.Errorf("expected %d lines, got %d", n, len(lines))
	}
	board := make([]string, n)
	for i := 0; i < n; i++ {
		if len(lines[i]) != m {
			return fmt.Errorf("line %d length %d expected %d", i+1, len(lines[i]), m)
		}
		for j := 0; j < m; j++ {
			if lines[i][j] != '.' && lines[i][j] != '*' {
				return fmt.Errorf("invalid character")
			}
		}
		board[i] = lines[i]
		if countSegments(board[i]) != rowSeg[i] {
			return fmt.Errorf("row %d segments mismatch", i+1)
		}
	}
	for j := 0; j < m; j++ {
		col := make([]byte, n)
		for i := 0; i < n; i++ {
			col[i] = board[i][j]
		}
		if countSegments(string(col)) != colSeg[j] {
			return fmt.Errorf("column %d segments mismatch", j+1)
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, []int, []int) {
	n := rng.Intn(5) + 1
	m := rng.Intn(20) + 1
	board := make([][]byte, n)
	for i := 0; i < n; i++ {
		board[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				board[i][j] = '.'
			} else {
				board[i][j] = '*'
			}
		}
	}
	row := make([]int, n)
	for i := 0; i < n; i++ {
		row[i] = countSegments(string(board[i]))
	}
	col := make([]int, m)
	for j := 0; j < m; j++ {
		colLine := make([]byte, n)
		for i := 0; i < n; i++ {
			colLine[i] = board[i][j]
		}
		col[j] = countSegments(string(colLine))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(row[i]))
	}
	sb.WriteByte('\n')
	for j := 0; j < m; j++ {
		if j > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(col[j]))
	}
	sb.WriteByte('\n')
	return sb.String(), row, col
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, row, col := generateCase(rng)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		nmc := strings.Fields(input)
		n, _ := strconv.Atoi(nmc[0])
		m, _ := strconv.Atoi(nmc[1])
		if err := checkBoard(out, n, m, row, col); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%soutput:%s\n", t+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
