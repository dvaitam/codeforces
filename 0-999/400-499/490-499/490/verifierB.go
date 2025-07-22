package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type caseB struct {
	n     int
	pairs [][2]int
	order []int
}

func generateTests() []caseB {
	r := rand.New(rand.NewSource(43))
	var tests []caseB
	tests = append(tests, genCase(r, 1))
	tests = append(tests, genCase(r, 2))
	for len(tests) < 120 {
		n := r.Intn(20) + 1
		tests = append(tests, genCase(r, n))
	}
	return tests
}

func genCase(r *rand.Rand, n int) caseB {
	ids := r.Perm(n)
	order := make([]int, n)
	for i, v := range ids {
		order[i] = v + 1
	}
	pairs := make([][2]int, n)
	for i := 0; i < n-1; i++ {
		pairs[i] = [2]int{order[i], order[i+1]}
	}
	pairs[n-1] = [2]int{order[n-1], 0}
	r.Shuffle(n, func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
	return caseB{n: n, pairs: pairs, order: order}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(pairs [][2]int) []int {
	next := make(map[int]int)
	prev := make(map[int]int)
	for _, pr := range pairs {
		next[pr[0]] = pr[1]
		prev[pr[1]] = pr[0]
	}
	var start int
	for k := range next {
		if prev[k] == 0 {
			start = k
			break
		}
	}
	res := make([]int, 0)
	cur := start
	for cur != 0 {
		res = append(res, cur)
		cur = next[cur]
	}
	return res
}

func verify(tc caseB, out string) error {
	parts := strings.Fields(out)
	if len(parts) != tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(parts))
	}
	got := make([]int, tc.n)
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("invalid integer")
		}
		got[i] = v
	}
	expected := solveCase(tc.pairs)
	for i := range expected {
		if got[i] != expected[i] {
			return fmt.Errorf("mismatch at position %d", i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, pr := range tc.pairs {
			sb.WriteString(fmt.Sprintf("%d %d\n", pr[0], pr[1]))
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(tc, out); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
