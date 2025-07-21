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

func solveK(input string) string {
	cmd := exec.Command("go", "run", "39K.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return out.String()
}

func ok(grid [][]byte, r, c int) bool {
	n := len(grid)
	m := len(grid[0])
	for i := r - 1; i <= r+1; i++ {
		for j := c - 1; j <= c+1; j++ {
			if i >= 0 && i < n && j >= 0 && j < m {
				if grid[i][j] == '*' {
					return false
				}
			}
		}
	}
	return true
}

func generateCaseK(rng *rand.Rand) string {
	n := rng.Intn(3) + 2
	m := rng.Intn(3) + 2
	k := rng.Intn(2) + 1
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	placed := 0
	attempts := 0
	for placed < k && attempts < 100 {
		r := rng.Intn(n)
		c := rng.Intn(m)
		if grid[r][c] == '.' && ok(grid, r, c) {
			grid[r][c] = '*'
			placed++
		}
		attempts++
	}
	k = placed
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseK(rng)
	}
	for i, tc := range cases {
		expect := solveK(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
