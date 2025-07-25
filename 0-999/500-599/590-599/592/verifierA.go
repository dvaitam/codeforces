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

func expected(board []string) string {
	minA := 8
	minB := 8
	for c := 0; c < 8; c++ {
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
		return "A"
	}
	return "B"
}

func generateCase(rng *rand.Rand) (string, string) {
	board := make([]string, 8)
	countW, countB := 0, 0
	for r := 0; r < 8; r++ {
		row := make([]byte, 8)
		for c := 0; c < 8; c++ {
			if r == 0 {
				row[c] = '.'
			} else if r == 7 {
				row[c] = '.'
			} else {
				v := rng.Intn(3)
				if v == 0 {
					row[c] = '.'
				} else if v == 1 {
					row[c] = 'W'
					countW++
				} else {
					row[c] = 'B'
					countB++
				}
			}
		}
		board[r] = string(row)
	}
	if countW == 0 {
		board[1] = "W" + board[1][1:]
	}
	if countB == 0 {
		board[6] = "B" + board[6][1:]
	}
	var sb strings.Builder
	for _, row := range board {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	exp := expected(board)
	return sb.String(), exp
}

func runCase(exe, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
