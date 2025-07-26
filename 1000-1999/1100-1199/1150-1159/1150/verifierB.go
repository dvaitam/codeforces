package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	input  string
	output string
}

func computeAnswer(n int, grid []string) string {
	board := make([][]byte, n)
	for i := 0; i < n; i++ {
		board[i] = []byte(grid[i])
	}
	for i := 1; i < n-1; i++ {
		for j := 1; j < n-1; j++ {
			if board[i][j] == '.' && board[i-1][j] == '.' && board[i+1][j] == '.' && board[i][j-1] == '.' && board[i][j+1] == '.' {
				board[i][j] = '#'
				board[i-1][j] = '#'
				board[i+1][j] = '#'
				board[i][j-1] = '#'
				board[i][j+1] = '#'
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if board[i][j] == '.' {
				return "NO"
			}
		}
	}
	return "YES"
}

func generateTests() []TestCase {
	var tests []TestCase
	for i := 1; i <= 120; i++ {
		n := (i % 4) + 3 // sizes 3..6
		grid := make([]string, n)
		for r := 0; r < n; r++ {
			row := make([]byte, n)
			for c := 0; c < n; c++ {
				if (i+r+c)%3 == 0 {
					row[c] = '#'
				} else {
					row[c] = '.'
				}
			}
			// ensure at least one dot
			if !bytes.Contains(row, []byte{'.'}) {
				row[0] = '.'
			}
			grid[r] = string(row)
		}
		expect := computeAnswer(n, grid)
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(grid, "\n"))
		tests = append(tests, TestCase{input: input, output: expect + "\n"})
	}
	return tests
}

func runTest(binary string, tc TestCase) error {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %v", err)
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(tc.output)
	if got != want {
		return fmt.Errorf("expected %q, got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB <binary>")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := runTest(binary, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
