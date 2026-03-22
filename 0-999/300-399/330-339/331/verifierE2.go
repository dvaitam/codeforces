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

const mod = 1000000007

type edge struct {
	u, v int
	seq  []int
}

type testCase struct {
	n     int
	edges []edge
	desc  string
}

// Embedded solver for 331E2 ported from /tmp/cf_r21_331_E2.go
func solve331E2(input string) string {
	fields := strings.Fields(input)
	pos := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(fields[pos])
		pos++
		return v
	}

	n := nextInt()
	m := nextInt()

	type Edge struct {
		to int
		V  string
	}

	type State struct {
		u   int
		typ int
		Q   string
	}

	adj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		k := nextInt()
		visions := make([]byte, k)
		for j := 0; j < k; j++ {
			visions[j] = byte(nextInt())
		}
		adj[u] = append(adj[u], Edge{to: v, V: string(visions)})
	}

	MOD := 1000000007

	curCount := make(map[State]int)
	curPath := make(map[State][]byte)

	for i := 1; i <= n; i++ {
		s := State{u: i, typ: 0, Q: string([]byte{byte(i)})}
		curCount[s] = 1
		curPath[s] = []byte{byte(i)}
	}

	var e1Path []byte
	ansE2 := make([]int, 2*n+1)

	for step := 1; step <= 2*n; step++ {
		nextCount := make(map[State]int)
		nextPath := make(map[State][]byte)

		for state, count := range curCount {
			path := curPath[state]
			if state.typ == 0 {
				for _, e := range adj[state.u] {
					T := state.Q + string([]byte{byte(e.to)})
					Ve := e.V

					var newState State
					valid := false

					if strings.HasPrefix(T, Ve) {
						newState = State{u: e.to, typ: 0, Q: T[len(Ve):]}
						valid = true
					} else if strings.HasPrefix(Ve, T) {
						newState = State{u: e.to, typ: 1, Q: Ve[len(T):]}
						valid = true
					}

					if valid {
						if newState.typ == 1 && len(newState.Q) > 2*n-step {
							continue
						}
						nextCount[newState] = (nextCount[newState] + count) % MOD
						if _, exists := nextPath[newState]; !exists {
							newP := make([]byte, len(path), len(path)+1)
							copy(newP, path)
							newP = append(newP, byte(e.to))
							nextPath[newState] = newP
						}
					}
				}
			} else {
				w := int(state.Q[0])
				for _, e := range adj[state.u] {
					if e.to == w {
						T := state.Q[1:] + e.V
						var newState State
						if len(T) == 0 {
							newState = State{u: w, typ: 0, Q: ""}
						} else {
							newState = State{u: w, typ: 1, Q: T}
						}

						if newState.typ == 1 && len(newState.Q) > 2*n-step {
							continue
						}

						nextCount[newState] = (nextCount[newState] + count) % MOD
						if _, exists := nextPath[newState]; !exists {
							newP := make([]byte, len(path), len(path)+1)
							copy(newP, path)
							newP = append(newP, byte(w))
							nextPath[newState] = newP
						}
						break
					}
				}
			}
		}

		curCount = nextCount
		curPath = nextPath

		ans := 0
		for state, count := range curCount {
			if state.typ == 0 && state.Q == "" {
				ans = (ans + count) % MOD
				if e1Path == nil && len(curPath[state]) <= 2*n {
					e1Path = curPath[state]
				}
			}
		}
		ansE2[step] = ans
	}

	var sb strings.Builder
	if e1Path == nil {
		sb.WriteString("0\n")
	} else {
		sb.WriteString(fmt.Sprintf("%d\n", len(e1Path)))
		for i, v := range e1Path {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(int(v)))
		}
		sb.WriteByte('\n')
	}

	for i := 1; i <= 2*n; i++ {
		sb.WriteString(strconv.Itoa(ansE2[i]))
		sb.WriteByte('\n')
	}

	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	for idx, tc := range tests {
		input := buildInput(tc)
		expStdout := solve331E2(input)

		expVals, err := parseOutput(expStdout, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.desc, err, expStdout)
			os.Exit(1)
		}

		gotStdout, gotStderr, err := runBinary(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d (%s): %v\nstderr:\n%s\n", idx+1, tc.desc, err, gotStderr)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotStdout, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d (%s): %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.desc, err, input, gotStdout)
			os.Exit(1)
		}
		if err := compareOutputs(expVals, gotVals); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, tc.desc, err, input, formatLines(expVals), formatLines(gotVals))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runBinary(path, input string) (string, string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func parseOutput(out string, n int) ([]int, error) {
	expected := 2 * n
	tokens := strings.Fields(out)
	if len(tokens) < expected {
		return nil, fmt.Errorf("expected at least %d integers for E2 counts, got %d total tokens", expected, len(tokens))
	}
	e2Tokens := tokens[len(tokens)-expected:]
	res := make([]int, expected)
	for i, tok := range e2Tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("E2 token %d: not an integer (%v)", i+1, err)
		}
		norm := int((val%int64(mod) + int64(mod)) % int64(mod))
		res[i] = norm
	}
	return res, nil
}

func compareOutputs(expected, actual []int) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("expected %d numbers, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("line %d: expected %d, got %d", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d %d", e.u, e.v, len(e.seq))
		for _, val := range e.seq {
			fmt.Fprintf(&sb, " %d", val)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func formatLines(vals []int) string {
	var sb strings.Builder
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{
		n: 6,
		edges: []edge{
			{1, 2, []int{1, 2}},
			{2, 3, []int{3}},
			{3, 4, []int{4, 5}},
			{4, 5, []int{}},
			{5, 3, []int{3}},
			{6, 1, []int{6}},
		},
		desc: "sample",
	})
	tests = append(tests, testCase{
		n:    1,
		desc: "single-node",
	})
	tests = append(tests, testCase{
		n: 2,
		edges: []edge{
			{1, 2, []int{1}},
		},
		desc: "single-edge",
	})
	tests = append(tests, testCase{
		n: 3,
		edges: []edge{
			{1, 2, []int{1, 2}},
			{2, 3, []int{2, 3}},
		},
		desc: "chain-with-prefix",
	})
	tests = append(tests, testCase{
		n: 4,
		edges: []edge{
			{1, 2, []int{}},
			{2, 3, []int{1, 2, 3}},
			{3, 4, []int{2, 3, 4}},
		},
		desc: "empty-vision-edge",
	})
	tests = append(tests, denseCase(6))
	tests = append(tests, cycleCase())

	rng := rand.New(rand.NewSource(1331331))
	for i := 0; i < 60; i++ {
		n := rng.Intn(8) + 2
		tc := randomCase(rng, n)
		tc.desc = fmt.Sprintf("random-%d", i+1)
		tests = append(tests, tc)
	}
	return tests
}

func denseCase(n int) testCase {
	var edges []edge
	flag := false
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			u, v := i, j
			if flag {
				u, v = v, u
			}
			flag = !flag
			edges = append(edges, edge{u: u, v: v, seq: []int{u, v}})
		}
	}
	return testCase{
		n:     n,
		edges: edges,
		desc:  "dense",
	}
}

func cycleCase() testCase {
	return testCase{
		n: 5,
		edges: []edge{
			{1, 2, []int{1}},
			{2, 3, []int{}},
			{3, 4, []int{3, 4}},
			{4, 5, []int{5}},
			{5, 1, []int{1}},
		},
		desc: "cycle",
	}
}

func randomCase(rng *rand.Rand, n int) testCase {
	type pair struct{ a, b int }
	var pairs []pair
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			pairs = append(pairs, pair{i, j})
		}
	}
	rng.Shuffle(len(pairs), func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})
	m := 0
	if len(pairs) > 0 {
		m = rng.Intn(len(pairs) + 1)
	}
	edges := make([]edge, 0, m)
	for i := 0; i < m; i++ {
		x, y := pairs[i].a, pairs[i].b
		if rng.Intn(2) == 0 {
			x, y = y, x
		}
		k := rng.Intn(n + 1)
		seq := make([]int, k)
		for j := 0; j < k; j++ {
			seq[j] = rng.Intn(n) + 1
		}
		edges = append(edges, edge{u: x, v: y, seq: seq})
	}
	return testCase{n: n, edges: edges}
}
