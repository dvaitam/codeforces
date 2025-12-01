package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "2114E.go"

type testCase struct {
	name  string
	cases []oneCase
	input string
}

type oneCase struct {
	n     int
	a     []int64
	edges [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		out, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if err := validate(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2114E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func validate(tc testCase, out string) error {
	reader := bufio.NewReader(strings.NewReader(out))
	for idx, c := range tc.cases {
		exp := computeThreats(c)
		for i := 0; i < c.n; i++ {
			var v int64
			if _, err := fmt.Fscan(reader, &v); err != nil {
				return fmt.Errorf("case %d: failed to read value %d: %v", idx+1, i+1, err)
			}
			if v != exp[i] {
				return fmt.Errorf("case %d: wrong value at vertex %d: expected %d, got %d", idx+1, i+1, exp[i], v)
			}
		}
	}

	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("unexpected extra output: %s", extra)
	}
	return nil
}

func computeThreats(c oneCase) []int64 {
	n := c.n
	adj := make([][]int, n+1)
	for _, e := range c.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	type frame struct {
		v, parent int
		parity    int
		diff      int64
		minDiff   int64
		maxDiff   int64
		idx       int
	}
	ans := make([]int64, n+1)
	stack := []frame{{v: 1, parent: 0, parity: 0, diff: c.a[0], minDiff: 0, maxDiff: 0, idx: 0}}

	// shift a to 1-based for convenience
	a := make([]int64, n+1)
	copy(a[1:], c.a)
	stack[0].diff = a[1]

	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		v := top.v

		if top.idx == 0 {
			if top.parity == 0 {
				ans[v] = top.diff - top.minDiff
			} else {
				ans[v] = top.maxDiff - top.diff
			}
		}

		if top.idx < len(adj[v]) {
			to := adj[v][top.idx]
			top.idx++
			if to == top.parent {
				continue
			}
			childParity := 1 - top.parity
			var childDiff int64
			if childParity == 0 {
				childDiff = top.diff + a[to]
			} else {
				childDiff = top.diff - a[to]
			}
			childMin := top.minDiff
			if top.diff < childMin {
				childMin = top.diff
			}
			childMax := top.maxDiff
			if top.diff > childMax {
				childMax = top.diff
			}
			stack = append(stack, frame{v: to, parent: v, parity: childParity, diff: childDiff, minDiff: childMin, maxDiff: childMax, idx: 0})
		} else {
			stack = stack[:len(stack)-1]
		}
	}

	return ans[1:]
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase(
			"chain-small",
			[][]int64{{4, 5}, {2, 6, 7}},
			[][][]int{{{1, 2}}, {{1, 2}, {1, 3}}},
		),
	}

	rng := rand.New(rand.NewSource(21140307))
	for i := 0; i < 60; i++ {
		cases := rng.Intn(4) + 1
		arr := make([][]int64, 0, cases)
		edges := make([][][]int, 0, cases)
		for j := 0; j < cases; j++ {
			n := rng.Intn(15) + 2
			a := make([]int64, n)
			for k := range a {
				a[k] = int64(rng.Intn(1_000_000_000) + 1)
			}
			edgeList := make([][]int, 0, n-1)
			for v := 2; v <= n; v++ {
				p := rng.Intn(v-1) + 1
				edgeList = append(edgeList, []int{p, v})
			}
			arr = append(arr, a)
			edges = append(edges, edgeList)
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), arr, edges))
	}
	return tests
}

func buildCase(name string, arrays [][]int64, edgeLists [][][]int) testCase {
	var sb strings.Builder
	t := len(arrays)
	fmt.Fprintf(&sb, "%d\n", t)
	cases := make([]oneCase, t)
	for i := 0; i < t; i++ {
		n := len(arrays[i])
		fmt.Fprintf(&sb, "%d\n", n)
		for j, v := range arrays[i] {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for _, e := range edgeLists[i] {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		cases[i] = oneCase{n: n, a: arrays[i], edges: toPairs(edgeLists[i])}
	}
	return testCase{name: name, cases: cases, input: sb.String()}
}

func toPairs(edges [][]int) [][2]int {
	res := make([][2]int, len(edges))
	for i, e := range edges {
		res[i] = [2]int{e[0], e[1]}
	}
	return res
}
