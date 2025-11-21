package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2000-2099/2050-2059/2050/2050G.go"

type treeCase struct {
	n     int
	edges [][2]int
}

type testCase struct {
	name  string
	cases []treeCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()

	for idx, tc := range tests {
		input := buildInput(tc)

		expected := evaluate(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutput(refOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		_ = refAns

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		if !equalAnswers(candAns, expected) {
			fmt.Fprintf(os.Stderr, "candidate mismatch on test %d (%s)\ninput:\n%soutput:\n%s", idx+1, tc.name, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2050G-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2050G.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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
	return out.String(), err
}

func parseOutput(output string, expected int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, token := range fields {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		res[i] = val
	}
	return res, nil
}

func equalAnswers(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.cases)))
	for _, c := range tc.cases {
		sb.WriteString(fmt.Sprintf("%d\n", c.n))
		for _, e := range c.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

func evaluate(tc testCase) []int64 {
	res := make([]int64, len(tc.cases))
	for idx, c := range tc.cases {
		res[idx] = solveCase(c)
	}
	return res
}

func solveCase(c treeCase) int64 {
	n := c.n
	adj := make([][]int, n+1)
	deg := make([]int64, n+1)
	for _, e := range c.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		deg[u]++
		deg[v]++
	}
	weight := make([]int64, n+1)
	for v := 1; v <= n; v++ {
		weight[v] = deg[v] - 2
	}

	bestDown := make([]int64, n+1)
	bestPath := int64(-1 << 60)

	type node struct {
		v, p int
		done bool
	}
	stack := []node{{v: 1, p: 0, done: false}}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !cur.done {
			stack = append(stack, node{v: cur.v, p: cur.p, done: true})
			for _, to := range adj[cur.v] {
				if to == cur.p {
					continue
				}
				stack = append(stack, node{v: to, p: cur.v, done: false})
			}
		} else {
			max1, max2 := int64(0), int64(0)
			for _, to := range adj[cur.v] {
				if to == cur.p {
					continue
				}
				val := bestDown[to]
				if val < 0 {
					val = 0
				}
				if val >= max1 {
					max2 = max1
					max1 = val
				} else if val > max2 {
					max2 = val
				}
			}
			bestDown[cur.v] = weight[cur.v] + max1
			if bestDown[cur.v] < weight[cur.v] {
				bestDown[cur.v] = weight[cur.v]
			}
			candidate := weight[cur.v] + max1 + max2
			if candidate > bestPath {
				bestPath = candidate
			}
		}
	}

	ans := bestPath + 2
	if ans < 0 {
		ans = 0
	}
	return ans
}

func buildTests() []testCase {
	sampleCases := []treeCase{
		{
			n:     2,
			edges: [][2]int{{1, 2}},
		},
		{
			n:     5,
			edges: [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}},
		},
	}
	tests := []testCase{{name: "samples", cases: sampleCases}}

	rng := rand.New(rand.NewSource(123456789))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(5) + 1
	cases := make([]treeCase, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(18) + 2
		edges := make([][2]int, 0, n-1)
		for v := 2; v <= n; v++ {
			parent := rng.Intn(v-1) + 1
			edges = append(edges, [2]int{parent, v})
		}
		cases[i] = treeCase{n: n, edges: edges}
	}
	return testCase{name: fmt.Sprintf("random_%d", idx+1), cases: cases}
}
