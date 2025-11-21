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

const mod = 998244353

type edge struct {
	to int
	w  int
}

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	exp, _ := strconv.ParseInt(expect, 10, 64)
	val, err := strconv.ParseInt(actual, 10, 64)
	if err != nil {
		return fmt.Errorf("output is not integer: %v", err)
	}
	exp = ((exp % mod) + mod) % mod
	val = ((val % mod) + mod) % mod
	if val != exp {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase(3, [][]edge{
			{{to: 2, w: 0}, {to: 3, w: 1}},
			{},
			{},
		}),
		makeCase(4, [][]edge{
			{{to: 2, w: 0}},
			{{to: 3, w: 1}},
			{{to: 4, w: 0}},
			{},
		}),
	}
	for i := 0; i < 50; i++ {
		n := rand.Intn(5) + 2
		tests = append(tests, randomCase(n))
	}
	return tests
}

func randomCase(n int) testCase {
	graph := make([][]edge, n)
	for i := 0; i < n; i++ {
		cnt := rand.Intn(2)
		for j := 0; j < cnt; j++ {
			to := rand.Intn(n-i-1) + i + 1
			graph[i] = append(graph[i], edge{to: to + 1, w: rand.Intn(2)})
		}
	}
	return makeCase(n, graph)
}

func makeCase(n int, graph [][]edge) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d\n", len(graph[i])))
		for _, e := range graph[i] {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.to, e.w))
		}
	}
	return testCase{
		input:  sb.String(),
		expect: fmt.Sprintf("%d", solveReference(n, graph)),
	}
}

func solveReference(n int, graph [][]edge) int64 {
	cnt0 := make([]int64, n+1)
	cnt1 := make([]int64, n+1)
	inv := make([]int64, n+1)
	for u := n; u >= 1; u-- {
		var p0, p1, invU int64
		for _, e := range graph[u-1] {
			if e.w == 1 {
				p1 = (p1 + 1) % mod
				invU = (invU + p1*cnt0[e.to] + inv[e.to]) % mod
				p1 = (p1 + cnt1[e.to]) % mod
				p0 = (p0 + cnt0[e.to]) % mod
			} else {
				invU = (invU + p1*cnt0[e.to] + inv[e.to]) % mod
				p1 = (p1 + cnt1[e.to]) % mod
				p0 = (p0 + cnt0[e.to]) % mod
				invU = (invU + p1) % mod
				p0 = (p0 + 1) % mod
			}
		}
		cnt0[u] = p0 % mod
		cnt1[u] = p1 % mod
		inv[u] = invU % mod
	}
	return inv[1] % mod
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
