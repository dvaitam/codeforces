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

func runSolution(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveA(r, c int, grid []string) int {
	rowsFree := make([]bool, r)
	colsFree := make([]bool, c)
	for i := 0; i < r; i++ {
		rowsFree[i] = true
	}
	for j := 0; j < c; j++ {
		colsFree[j] = true
	}
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if grid[i][j] == 'S' {
				rowsFree[i] = false
				colsFree[j] = false
			}
		}
	}
	freeRows := 0
	for i := 0; i < r; i++ {
		if rowsFree[i] {
			freeRows++
		}
	}
	freeCols := 0
	for j := 0; j < c; j++ {
		if colsFree[j] {
			freeCols++
		}
	}
	return freeRows*c + freeCols*(r-freeRows)
}

func generateCaseA(rng *rand.Rand) (string, string) {
	r := rng.Intn(9) + 2
	c := rng.Intn(9) + 2
	grid := make([]string, r)
	for i := 0; i < r; i++ {
		b := make([]byte, c)
		for j := 0; j < c; j++ {
			if rng.Intn(5) == 0 {
				b[j] = 'S'
			} else {
				b[j] = '.'
			}
		}
		grid[i] = string(b)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", r, c)
	for i := 0; i < r; i++ {
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	expected := solveA(r, c, grid)
	return sb.String(), fmt.Sprintf("%d", expected)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct{ in, exp string }{
		{"2 2\n..\n..\n", "4"},
		{"2 2\nSS\nSS\n", "0"},
		{"3 3\nS..\n.S.\n..S\n", "0"},
	}
	for i := len(cases); i < 100; i++ {
		in, exp := generateCaseA(rng)
		cases = append(cases, struct{ in, exp string }{in, exp})
	}

	for i, tc := range cases {
		out, err := runSolution(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(tc.exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.exp, out, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
