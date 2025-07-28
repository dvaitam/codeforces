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
	grid [][]byte
}

func simulate(grid [][]byte) []string {
	n := len(grid)
	if n == 0 {
		return []string{}
	}
	m := len(grid[0])
	result := make([][]byte, n)
	for i := range result {
		result[i] = make([]byte, m)
		for j := range result[i] {
			result[i][j] = '.'
		}
	}
	for c := 0; c < m; c++ {
		pos := n - 1
		for r := n - 1; r >= 0; r-- {
			switch grid[r][c] {
			case 'o':
				result[r][c] = 'o'
				pos = r - 1
			case '*':
				result[pos][c] = '*'
				pos--
			}
		}
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = string(result[i])
	}
	return out
}

func solveCase(tc testCase) string {
	res := simulate(tc.grid)
	return strings.Join(res, "\n")
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	letters := []byte{'.', '*', 'o'}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			row[j] = letters[rng.Intn(3)]
		}
		grid[i] = row
	}
	tc := testCase{grid: grid}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	input := sb.String()
	output := solveCase(tc) + "\n"
	return input, output
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\n\ngot \n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
