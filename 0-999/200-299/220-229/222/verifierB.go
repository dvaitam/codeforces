package main

import (
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

func solve(input string) string {
	reader := strings.NewReader(input)
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return ""
	}
	table := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &row[j])
		}
		table[i] = row
	}
	rowIdx := make([]int, n)
	for i := 0; i < n; i++ {
		rowIdx[i] = i
	}
	colIdx := make([]int, m)
	for j := 0; j < m; j++ {
		colIdx[j] = j
	}
	var out strings.Builder
	for q := 0; q < k; q++ {
		var op string
		var x, y int
		fmt.Fscan(reader, &op, &x, &y)
		switch op[0] {
		case 'r':
			rowIdx[x-1], rowIdx[y-1] = rowIdx[y-1], rowIdx[x-1]
		case 'c':
			colIdx[x-1], colIdx[y-1] = colIdx[y-1], colIdx[x-1]
		case 'g':
			r := rowIdx[x-1]
			c := colIdx[y-1]
			fmt.Fprintf(&out, "%d\n", table[r][c])
		}
	}
	return strings.TrimRight(out.String(), "\n")
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	var tests []test
	fixed := []string{
		"1 1 1\n5\ng 1 1\n",
		"2 2 3\n1 2\n3 4\nr 1 2\nc 1 2\ng 1 1\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(3) + 1
		m := rng.Intn(3) + 1
		k := rng.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				fmt.Fprintf(&sb, "%d ", rng.Intn(10))
			}
			sb.WriteByte('\n')
		}
		ops := []byte{'r', 'c', 'g'}
		for i := 0; i < k; i++ {
			op := ops[rng.Intn(len(ops))]
			var x, y int
			switch op {
			case 'r':
				x = rng.Intn(n) + 1
				y = rng.Intn(n) + 1
			case 'c':
				x = rng.Intn(m) + 1
				y = rng.Intn(m) + 1
			case 'g':
				x = rng.Intn(n) + 1
				y = rng.Intn(m) + 1
			}
			fmt.Fprintf(&sb, "%c %d %d\n", op, x, y)
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
