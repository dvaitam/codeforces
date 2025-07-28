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

func run(bin, input string) (string, error) {
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

func parseMatrix(out string, n, m int) ([][]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != n {
		return nil, fmt.Errorf("expected %d lines, got %d", n, len(lines))
	}
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		fields := strings.Fields(lines[i])
		if len(fields) != m {
			return nil, fmt.Errorf("line %d: expected %d values, got %d", i+1, m, len(fields))
		}
		row := make([]int, m)
		for j, f := range fields {
			if f != "0" && f != "1" {
				return nil, fmt.Errorf("line %d: invalid value %q", i+1, f)
			}
			if f == "1" {
				row[j] = 1
			}
		}
		mat[i] = row
	}
	return mat, nil
}

func checkMatrix(mat [][]int) error {
	n := len(mat)
	m := len(mat[0])
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			diff := 0
			v := mat[i][j]
			if i > 0 && mat[i-1][j] != v {
				diff++
			}
			if i+1 < n && mat[i+1][j] != v {
				diff++
			}
			if j > 0 && mat[i][j-1] != v {
				diff++
			}
			if j+1 < m && mat[i][j+1] != v {
				diff++
			}
			if diff != 2 {
				return fmt.Errorf("cell %d,%d diff %d", i+1, j+1, diff)
			}
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int, int) {
	n := rng.Intn(10) + 2
	if n%2 == 1 {
		n++
	}
	m := rng.Intn(10) + 2
	if m%2 == 1 {
		m++
	}
	input := fmt.Sprintf("1\n%d %d\n", n, m)
	return input, n, m
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n, m := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		mat, err := parseMatrix(out, n, m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
		if err := checkMatrix(mat); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
