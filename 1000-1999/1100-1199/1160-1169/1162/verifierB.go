package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func isIncreasing(mat [][]int) bool {
	n := len(mat)
	if n == 0 {
		return true
	}
	m := len(mat[0])
	for i := 0; i < n; i++ {
		for j := 1; j < m; j++ {
			if mat[i][j] <= mat[i][j-1] {
				return false
			}
		}
	}
	for j := 0; j < m; j++ {
		for i := 1; i < n; i++ {
			if mat[i][j] <= mat[i-1][j] {
				return false
			}
		}
	}
	return true
}

func solveCase(input string) string {
	in := bufio.NewReader(strings.NewReader(strings.TrimSpace(input)))
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return ""
	}
	a := make([][]int, n)
	b := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}
	for i := 0; i < n; i++ {
		b[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &b[i][j])
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if a[i][j] < b[i][j] {
				a[i][j], b[i][j] = b[i][j], a[i][j]
			}
		}
	}
	if isIncreasing(a) && isIncreasing(b) {
		return "Possible"
	}
	return "Impossible"
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	var tests []test
	// simple deterministic cases
	tests = append(tests, test{input: "1 1\n1\n1\n", expected: solveCase("1 1\n1\n1\n")})
	tests = append(tests, test{input: "2 2\n1 4\n5 6\n2 3\n4 5\n", expected: solveCase("2 2\n1 4\n5 6\n2 3\n4 5\n")})
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				sb.WriteString(fmt.Sprintf("%d ", rng.Intn(1000)+1))
			}
			sb.WriteByte('\n')
		}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				sb.WriteString(fmt.Sprintf("%d ", rng.Intn(1000)+1))
			}
			sb.WriteByte('\n')
		}
		inp := sb.String()
		tests = append(tests, test{inp, solveCase(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
