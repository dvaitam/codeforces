package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Query struct{ l, r, x int }

func solveC(n int, m int, a []int, qs []Query) []int {
	res := make([]int, m)
	for i, q := range qs {
		pos := -1
		for j := q.l - 1; j <= q.r-1; j++ {
			if a[j] != q.x {
				pos = j + 1
				break
			}
		}
		res[i] = pos
	}
	return res
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	type Test struct {
		n, m int
		a    []int
		qs   []Query
	}
	var tests []Test
	// small deterministic tests
	tests = append(tests, Test{n: 1, m: 1, a: []int{1}, qs: []Query{{1, 1, 1}}})
	tests = append(tests, Test{n: 2, m: 2, a: []int{1, 2}, qs: []Query{{1, 2, 1}, {1, 2, 2}}})
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		m := rand.Intn(20) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rand.Intn(100) + 1
		}
		qs := make([]Query, m)
		for i := range qs {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			x := rand.Intn(100) + 1
			qs[i] = Query{l, r, x}
		}
		tests = append(tests, Test{n: n, m: m, a: a, qs: qs})
	}

	for idx, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.m))
		for i, val := range t.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
		for _, q := range t.qs {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", q.l, q.r, q.x))
		}
		input := sb.String()
		expRes := solveC(t.n, t.m, t.a, t.qs)
		expected := strings.TrimSpace(strings.Join(func() []string {
			out := make([]string, len(expRes))
			for i, v := range expRes {
				out[i] = fmt.Sprintf("%d", v)
			}
			return out
		}(), "\n"))
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
