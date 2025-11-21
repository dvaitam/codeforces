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
	refSource2018C = "2018C.go"
	refBinary2018C = "ref2018C.bin"
	maxTests       = 100
	maxTotalN      = 500000
)

type testCase struct {
	n     int
	edges [][2]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
			fmt.Printf("Mismatch on case %d: expected %d got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2018C, refSource2018C)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2018C), nil
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
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2018))
	var tests []testCase
	totalN := 0

	addCase := func(tc testCase) {
		tests = append(tests, tc)
		totalN += tc.n
	}

	// hand-crafted small cases
	addCase(testCase{
		n: 3,
		edges: [][2]int{
			{1, 2},
			{1, 3},
		},
	})
	addCase(testCase{
		n: 4,
		edges: [][2]int{
			{1, 2},
			{2, 3},
			{3, 4},
		},
	})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		maxN := 10000
		if remain < maxN {
			maxN = remain
		}
		n := rnd.Intn(maxN-2) + 3
		edges := randomTree(n, rnd)
		addCase(testCase{n: n, edges: edges})
	}
	return tests
}

func randomTree(n int, rnd *rand.Rand) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rnd.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	// random shuffle for variety
	rnd.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})
	return edges
}
