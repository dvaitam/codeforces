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

type testCase struct {
	input  string
	expect string
}

type edge struct {
	u, v int
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
	exp, _ := strconv.Atoi(expect)
	val, err := strconv.Atoi(actual)
	if err != nil {
		return fmt.Errorf("output not integer: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase(3, []edge{{1, 2}, {2, 3}}),
		makeCase(4, []edge{{1, 2}, {2, 3}, {3, 4}, {4, 1}}),
		makeCase(5, []edge{{1, 2}, {2, 3}, {2, 4}, {3, 5}}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 1
		m := rand.Intn(n*(n-1)/2 + 1)
		edges := randomEdges(n, m)
		tests = append(tests, makeCase(n, edges))
	}
	return tests
}

func randomEdges(n, m int) []edge {
	set := map[[2]int]struct{}{}
	result := make([]edge, 0, m)
	for len(result) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if _, ok := set[key]; ok {
			continue
		}
		set[key] = struct{}{}
		result = append(result, edge{u, v})
	}
	return result
}

func makeCase(n int, edges []edge) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	return testCase{
		input:  sb.String(),
		expect: fmt.Sprintf("%d", solveReference(n, edges)),
	}
}

func solveReference(n int, edges []edge) int {
	l, r := 1, n
	ans := 1
	for l <= r {
		mid := (l + r) / 2
		if feasible(n, edges, mid) {
			ans = mid
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return ans
}

func feasible(n int, edges []edge, D int) bool {
	cover := make([]bool, n)
	for i := range cover {
		cover[i] = true
	}
	limit := D
	if limit < 1 {
		return false
	}
	dp := make([]int, n)
	var dfs func(int, int) bool
	dfs = func(idx, last int) bool {
		if idx == n {
			for _, e := range edges {
				if !cover[e.u-1] && !cover[e.v-1] {
					return false
				}
			}
			return true
		}
		for _, choice := range []bool{true, false} {
			cover[idx] = choice
			if choice && last != -1 && idx-last < limit {
				continue
			}
			if dfs(idx+1, func() int {
				if choice {
					return idx
				}
				return last
			}()) {
				return true
			}
		}
		return false
	}
	return dfs(0, -limit)
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
