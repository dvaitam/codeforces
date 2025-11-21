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

const (
	refSource2164E = "2164E.go"
	refBinary2164E = "ref2164E.bin"
	maxTests       = 120
	maxTotalN      = 300000
	maxTotalM      = 300000
)

type edge struct {
	u int
	v int
	w int64
}

type testCase struct {
	n     int
	m     int
	edges []edge
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on test %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2164E, refSource2164E)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2164E), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d not an integer: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2164))
	var tests []testCase
	totalN := 0
	totalM := 0

	addCase := func(tc testCase) {
		tests = append(tests, tc)
		totalN += tc.n
		totalM += tc.m
	}

	addCase(testCase{n: 1, m: 0})
	addCase(testCase{
		n: 2, m: 1,
		edges: []edge{{u: 1, v: 2, w: 5}},
	})
	addCase(testCase{
		n: 3, m: 3,
		edges: []edge{
			{1, 2, 3},
			{2, 3, 4},
			{3, 1, 5},
		},
	})

	for len(tests) < maxTests && totalN < maxTotalN && totalM < maxTotalM {
		maxN := maxTotalN - totalN
		if maxN <= 0 {
			break
		}
		limitN := 500
		if maxN < limitN {
			limitN = maxN
		}
		n := rnd.Intn(limitN) + 1

		minEdges := 0
		if n > 1 {
			minEdges = n - 1
		}
		remainingEdges := maxTotalM - totalM
		if remainingEdges < minEdges {
			break
		}
		maxExtra := remainingEdges - minEdges
		extraLimit := n
		if maxExtra > extraLimit {
			maxExtra = extraLimit
		}
		extra := 0
		if maxExtra > 0 {
			extra = rnd.Intn(maxExtra + 1)
		}
		m := minEdges + extra

		tc := testCase{n: n, m: m, edges: make([]edge, 0, m)}
		if n > 1 {
			for v := 2; v <= n; v++ {
				u := rnd.Intn(v-1) + 1
				tc.edges = append(tc.edges, edge{
					u: u,
					v: v,
					w: rndWeight(rnd),
				})
			}
		}
		for len(tc.edges) < m {
			u := rnd.Intn(n) + 1
			v := rnd.Intn(n) + 1
			tc.edges = append(tc.edges, edge{
				u: u,
				v: v,
				w: rndWeight(rnd),
			})
		}
		addCase(tc)
	}
	return tests
}

func rndWeight(rnd *rand.Rand) int64 {
	base := rnd.Int63n(1_000_000_000) + 1
	return base
}
