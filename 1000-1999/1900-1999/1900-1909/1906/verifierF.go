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

type operation struct {
	l, r, x int
}

type query struct {
	k, s, t int
}

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
	expLines := strings.Split(strings.TrimSpace(expect), "\n")
	actLines := strings.Split(strings.TrimSpace(actual), "\n")
	if len(expLines) != len(actLines) {
		return fmt.Errorf("expected %d lines but got %d", len(expLines), len(actLines))
	}
	for i := range expLines {
		expVal, _ := strconv.ParseInt(strings.TrimSpace(expLines[i]), 10, 64)
		actVal, err := strconv.ParseInt(strings.TrimSpace(actLines[i]), 10, 64)
		if err != nil {
			return fmt.Errorf("line %d: output not integer: %v", i+1, err)
		}
		if expVal != actVal {
			return fmt.Errorf("line %d: expected %d but got %d", i+1, expVal, actVal)
		}
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase(3, []operation{{1, 3, 5}, {2, 3, 4}}, []query{{2, 1, 2}}),
		makeCase(4, []operation{{1, 2, 3}, {2, 4, -1}, {1, 4, 2}}, []query{{3, 1, 3}, {1, 2, 3}}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		ops := make([]operation, m)
		for j := 0; j < m; j++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			x := rand.Intn(200001) - 100000
			ops[j] = operation{l: l, r: r, x: x}
		}
		q := rand.Intn(5) + 1
		queries := make([]query, q)
		for j := 0; j < q; j++ {
			k := rand.Intn(n) + 1
			s := rand.Intn(m) + 1
			t := rand.Intn(m-s+1) + s
			queries[j] = query{k: k, s: s, t: t}
		}
		tests = append(tests, makeCase(n, ops, queries))
	}
	return tests
}

func makeCase(n int, ops []operation, queries []query) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(ops))
	for _, op := range ops {
		fmt.Fprintf(&sb, "%d %d %d\n", op.l, op.r, op.x)
	}
	fmt.Fprintf(&sb, "%d\n", len(queries))
	for _, q := range queries {
		fmt.Fprintf(&sb, "%d %d %d\n", q.k, q.s, q.t)
	}
	return testCase{
		input:  sb.String(),
		expect: solveReference(n, ops, queries),
	}
}

func solveReference(n int, ops []operation, queries []query) string {
	var sb strings.Builder
	for _, q := range queries {
		best := int64(^uint64(0) >> 1)
		best = -best - 1 // minimal int64
		if best != 0 {
			best = int64(-1 << 63)
		}
		curr := int64(0)
		for i := q.s - 1; i <= q.t-1; i++ {
			val := int64(0)
			if ops[i].l <= q.k && q.k <= ops[i].r {
				val = int64(ops[i].x)
			}
			if curr > 0 {
				curr += val
			} else {
				curr = val
			}
			if curr > best {
				best = curr
			}
		}
		fmt.Fprintf(&sb, "%d\n", best)
	}
	return strings.TrimSpace(sb.String())
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
