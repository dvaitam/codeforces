package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

func solveRef(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return "", err
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}
	const INF int = int(1e9)
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			bottom := INF
			right := INF
			if i+1 < n {
				bottom = a[i+1][j]
			}
			if j+1 < m {
				right = a[i][j+1]
			}
			if a[i][j] == 0 {
				val := bottom - 1
				if right-1 < val {
					val = right - 1
				}
				if val <= 0 {
					return "-1", nil
				}
				a[i][j] = val
			} else {
				if a[i][j] >= bottom || a[i][j] >= right {
					return "-1", nil
				}
			}
		}
	}
	var sum int64
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j+1 < m && a[i][j] >= a[i][j+1] {
				return "-1", nil
			}
			if i+1 < n && a[i][j] >= a[i+1][j] {
				return "-1", nil
			}
			sum += int64(a[i][j])
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func makeCase(name string, matrix [][]int) testCase {
	n := len(matrix)
	m := len(matrix[0])
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", matrix[i][j])
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("example1", [][]int{
			{1, 2, 3},
			{2, 0, 5},
			{3, 4, 6},
		}),
		makeCase("example2", [][]int{
			{1, 3, 5},
			{2, 0, 6},
			{4, 7, 8},
		}),
		makeCase("impossible", [][]int{
			{1, 5, 9},
			{3, 0, 10},
			{6, 8, 11},
		}),
		makeCase("allZerosInterior", [][]int{
			{1, 5, 9, 12},
			{2, 0, 0, 13},
			{3, 0, 0, 14},
			{4, 6, 7, 15},
		}),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for t := 0; t < 80; t++ {
		n := rng.Intn(5) + 3
		m := rng.Intn(5) + 3
		matrix := make([][]int, n)
		for i := 0; i < n; i++ {
			matrix[i] = make([]int, m)
			for j := 0; j < m; j++ {
				if i == 0 || j == 0 || i == n-1 || j == m-1 {
					matrix[i][j] = rng.Intn(8000) + 1
				} else {
					if rng.Intn(3) == 0 {
						matrix[i][j] = 0
					} else {
						matrix[i][j] = rng.Intn(8000) + 1
					}
				}
			}
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", t+1), matrix))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to solve reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("test %d (%s) mismatch\ninput:\n%s\nexpect:%s\nactual:%s\n", idx+1, tc.name, tc.input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
