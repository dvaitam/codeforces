package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./2129E.go"

type testCase struct {
	n, m    int
	edges   [][2]int
	queries [][3]int
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[len(args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	exp, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range exp {
		if len(exp[i]) != len(got[i]) {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "case %d (n=%d m=%d q=%d): expected %d answers, got %d\n", i+1, tc.n, tc.m, len(tc.queries), len(exp[i]), len(got[i]))
			os.Exit(1)
		}
		for j := range exp[i] {
			if exp[i][j] != got[i][j] {
				tc := tests[i]
				fmt.Fprintf(os.Stderr, "wrong answer on case %d query %d (n=%d m=%d): expected %d got %d\n", i+1, j+1, tc.n, tc.m, exp[i][j], got[i][j])
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2129E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func buildTests() []testCase {
	var tcs []testCase
	// Sample-based small graph.
	tcs = append(tcs, sampleCase())

	// Random and crafted cases.
	rng := rand.New(rand.NewSource(2129))

	tcs = append(tcs, randomCase(rng, 6, 8, 6))
	tcs = append(tcs, randomCase(rng, 20, 40, 30))
	tcs = append(tcs, randomCase(rng, 1000, 2000, 1500))
	tcs = append(tcs, randomCase(rng, 8000, 12000, 9000))

	// Large stress while keeping totals within constraints.
	tcs = append(tcs, randomCase(rng, 50000, 70000, 30000))

	return tcs
}

func sampleCase() testCase {
	n, m := 4, 5
	edges := [][2]int{{1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}}
	queries := [][3]int{
		{1, 2, 2},
		{1, 3, 1},
		{2, 4, 3},
	}
	return testCase{n: n, m: m, edges: edges, queries: queries}
}

func randomCase(rng *rand.Rand, n, m, q int) testCase {
	if m > n*(n-1)/2 {
		m = n * (n - 1) / 2
	}
	type pair struct{ u, v int }
	seen := make(map[pair]struct{})
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		p := pair{u, v}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		edges = append(edges, [2]int{u, v})
	}

	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		k := rng.Intn(r-l+1) + 1
		queries[i] = [3]int{l, r, k}
	}

	return testCase{n: n, m: m, edges: edges, queries: queries}
}

func buildInput(tcs []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tcs)))
	sb.WriteByte('\n')
	for _, tc := range tcs {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		sb.WriteString(strconv.Itoa(len(tc.queries)))
		sb.WriteByte('\n')
		for _, q := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", q[0], q[1], q[2]))
		}
	}
	return sb.String()
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutputs(out string, tests []testCase) ([][]int, error) {
	tokens := strings.Fields(out)
	pos := 0
	res := make([][]int, len(tests))
	for i, tc := range tests {
		q := len(tc.queries)
		if pos+q > len(tokens) {
			return nil, fmt.Errorf("case %d: expected %d answers, got %d", i+1, q, len(tokens)-pos)
		}
		ans := make([]int, q)
		for j := 0; j < q; j++ {
			val, err := strconv.Atoi(tokens[pos+j])
			if err != nil {
				return nil, fmt.Errorf("case %d answer %d: invalid integer %q", i+1, j+1, tokens[pos+j])
			}
			ans[j] = val
		}
		res[i] = ans
		pos += q
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("extra tokens detected (%d unused)", len(tokens)-pos)
	}
	return res, nil
}
