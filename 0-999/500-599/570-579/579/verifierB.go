package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Test struct {
	n   int
	mat [][]int
}

func generateTests() []Test {
	var tests []Test
	for t := 0; t < 100; t++ {
		rand.Seed(int64(t + 1))
		n := rand.Intn(4) + 2 // n in [2,5], so 2n in [4,10]
		m := 2 * n
		total := m * (m - 1) / 2
		perm := rand.Perm(total)
		idx := 0
		mat := make([][]int, m+1)
		for i := range mat {
			mat[i] = make([]int, m+1)
		}
		for i := 2; i <= m; i++ {
			for j := 1; j < i; j++ {
				mat[i][j] = perm[idx] + 1
				idx++
			}
		}
		tests = append(tests, Test{n: n, mat: mat})
	}
	return tests
}

func solve(n int, mat [][]int) []int {
	m := 2 * n
	type Pair struct{ i, j, v int }
	var pairs []Pair
	for i := 2; i <= m; i++ {
		for j := 1; j < i; j++ {
			pairs = append(pairs, Pair{i, j, mat[i][j]})
		}
	}
	sort.Slice(pairs, func(a, b int) bool { return pairs[a].v > pairs[b].v })
	res := make([]int, m+1)
	for _, p := range pairs {
		if res[p.i] == 0 && res[p.j] == 0 {
			res[p.i] = p.j
			res[p.j] = p.i
		}
	}
	return res[1:]
}

func (t Test) input() string {
	var b strings.Builder
	m := 2 * t.n
	b.WriteString(fmt.Sprintf("%d\n", t.n))
	for i := 2; i <= m; i++ {
		for j := 1; j < i; j++ {
			b.WriteString(fmt.Sprintf("%d ", t.mat[i][j]))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func runBinary(binary, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func intsToString(arr []int) string {
	var b strings.Builder
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprintf("%d", v))
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	tests := generateTests()
	for idx, t := range tests {
		exp := intsToString(solve(t.n, t.mat))
		output, err := runBinary(binary, t.input())
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if output != exp {
			fmt.Printf("Test %d failed:\nInput:\n%sExpected: %s\nGot: %s\n", idx+1, t.input(), exp, output)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(tests))
}
