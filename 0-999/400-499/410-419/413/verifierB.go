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
	n, m, k int
	matrix  [][]bool
	events  [][2]int
}

func expected(tc testCase) []int {
	totalMsgs := make([]int, tc.m)
	userMsgs := make([][]int, tc.n)
	for i := range userMsgs {
		userMsgs[i] = make([]int, tc.m)
	}
	for _, e := range tc.events {
		u := e[0]
		c := e[1]
		totalMsgs[c]++
		userMsgs[u][c]++
	}
	res := make([]int, tc.n)
	for i := 0; i < tc.n; i++ {
		sum := 0
		for j := 0; j < tc.m; j++ {
			if tc.matrix[i][j] {
				sum += totalMsgs[j] - userMsgs[i][j]
			}
		}
		res[i] = sum
	}
	return res
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			val := 0
			if tc.matrix[i][j] {
				val = 1
			}
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(val))
		}
		sb.WriteByte('\n')
	}
	for _, e := range tc.events {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	gotFields := strings.Fields(strings.TrimSpace(out.String()))
	exp := expected(tc)
	if len(gotFields) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(gotFields))
	}
	for i, f := range gotFields {
		var v int
		fmt.Sscan(f, &v)
		if v != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, gotFields)
		}
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 2
	m := rng.Intn(3) + 1
	k := rng.Intn(30)
	matrix := make([][]bool, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]bool, m)
	}
	// ensure each chat has at least two participants
	for j := 0; j < m; j++ {
		p1 := rng.Intn(n)
		p2 := rng.Intn(n)
		for p2 == p1 {
			p2 = rng.Intn(n)
		}
		matrix[p1][j] = true
		matrix[p2][j] = true
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				matrix[i][j] = true
			}
		}
	}
	events := make([][2]int, k)
	for i := 0; i < k; i++ {
		u := rng.Intn(n)
		c := rng.Intn(m)
		for !matrix[u][c] {
			u = rng.Intn(n)
			c = rng.Intn(m)
		}
		events[i] = [2]int{u, c}
	}
	return testCase{n: n, m: m, k: k, matrix: matrix, events: events}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	// simple deterministic case
	matrix := [][]bool{{true}, {true}}
	tests = append(tests, testCase{n: 2, m: 1, k: 1, matrix: matrix, events: [][2]int{{0, 0}}})
	for len(tests) < 100 {
		tests = append(tests, randomCase(rng))
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
