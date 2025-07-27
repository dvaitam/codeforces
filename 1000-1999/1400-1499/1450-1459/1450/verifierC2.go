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
	k := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] != '.' {
				k++
			}
		}
	}
	limit := k / 3
	choices := [][3]byte{{'*', 'O', 'X'}, {'*', 'X', 'O'}, {'O', '*', 'X'}, {'X', '*', 'O'}, {'O', 'X', '*'}, {'X', 'O', '*'}}
	var use [3]byte
	for _, ch := range choices {
		cnt := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				c := grid[i][j]
				if c == '.' {
					continue
				}
				d := (i + j) % 3
				if ch[d] != '*' && c != ch[d] {
					cnt++
				}
			}
		}
		if cnt <= limit {
			use = ch
			break
		}
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		row := []byte(grid[i])
		for j := 0; j < n; j++ {
			c := row[j]
			if c == '.' {
				continue
			}
			d := (i + j) % 3
			if use[d] != '*' {
				row[j] = use[d]
			}
		}
		res[i] = string(row)
	}
	return res
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
	expected := strings.Join(solve(grid), "\n")
	return testCase{input: sb.String(), expected: expected}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			r := rng.Intn(3)
			switch r {
			case 0:
				b[j] = '.'
			case 1:
				b[j] = 'X'
			default:
				b[j] = 'O'
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/binary")
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
