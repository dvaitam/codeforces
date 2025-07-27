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

type testCase struct {
	input    string
	expected string
}

func solve(grid []string) []string {
	n := len(grid)
	cnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 'X' {
				cnt++
			}
		}
	}
	for k := 0; k < 3; k++ {
		c := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if (i+j)%3 == k && grid[i][j] == 'X' {
					c++
				}
			}
		}
		if c <= cnt/3 {
			res := make([]string, n)
			for i := 0; i < n; i++ {
				row := []byte(grid[i])
				for j := 0; j < n; j++ {
					if (i+j)%3 == k && row[j] == 'X' {
						row[j] = 'O'
					}
				}
				res[i] = string(row)
			}
			return res
		}
	}
	return grid
}

func buildCase(grid []string) testCase {
	n := len(grid)
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	out := solve(grid)
	expected := strings.Join(out, "\n")
	return testCase{input: sb.String(), expected: expected}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				b[j] = 'X'
			} else {
				b[j] = '.'
			}
		}
		grid[i] = string(b)
	}
	return buildCase(grid)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected:\n%s\n----\ngot:\n%s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
