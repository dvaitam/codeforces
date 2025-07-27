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

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, m, x, y int, grid []string) string {
	total := 0
	for i := 0; i < n; i++ {
		row := grid[i]
		j := 0
		for j < m {
			if row[j] == '*' {
				j++
				continue
			}
			if j+1 < m && row[j+1] == '.' && y < 2*x {
				total += y
				j += 2
			} else {
				total += x
				j++
			}
		}
	}
	return fmt.Sprintf("%d", total)
}

type testCase struct {
	n, m, x, y int
	grid       []string
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	x := rng.Intn(1000) + 1
	y := rng.Intn(1000) + 1
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				b[j] = '.'
			} else {
				b[j] = '*'
			}
		}
		grid[i] = string(b)
	}
	return testCase{n, m, x, y, grid}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.m, tc.x, tc.y))
		for _, row := range tc.grid {
			sb.WriteString(row + "\n")
		}
		input := sb.String()
		want := expected(tc.n, tc.m, tc.x, tc.y, tc.grid)
		got, err := runProg(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
