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

const refSource = "./2039G.go"

type testCase struct {
	n     int
	m     int
	edges [][2]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()

	for i, tc := range tests {
		input := renderInput(tc)
		expOut, err := runWithInput(exec.Command(refBin), input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\noutput:\n%s\n", i+1, err, expOut)
			os.Exit(1)
		}
		gotOut, err := runWithInput(commandFor(candidate), input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\noutput:\n%s\ninput:\n%s", i+1, err, gotOut, input)
			os.Exit(1)
		}

		expVal, err := parseSingle(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, expOut)
			os.Exit(1)
		}
		gotVal, err := parseSingle(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\ninput:\n%s", i+1, err, gotOut, input)
			os.Exit(1)
		}
		if expVal != gotVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\ninput:\n%s", i+1, expVal, gotVal, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2039G-ref-*")
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

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseSingle(output string) (int64, error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return 0, fmt.Errorf("no output")
	}
	if len(fields) > 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse int: %v", err)
	}
	return val, nil
}

func renderInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func generateTests() []testCase {
	tests := []testCase{
		// Samples
		{n: 6, m: 6, edges: edgesFromPairs(6, [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}, {3, 6}})},
		{n: 2, m: 5, edges: edgesFromPairs(2, [][2]int{{1, 2}})},
		{n: 12, m: 69, edges: edgesFromPairs(12, [][2]int{{3, 5}, {1, 4}, {2, 3}, {4, 5}, {5, 6}, {8, 9}, {7, 3}, {4, 8}, {9, 10}, {1, 11}, {12, 1}})},
	}

	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	const maxTotal = 20000
	rng := rand.New(rand.NewSource(2039))

	for totalN < maxTotal {
		n := rng.Intn(3000) + 2
		if totalN+n > maxTotal {
			n = maxTotal - totalN
		}
		m := rng.Intn(1_000_000_000) + 1
		edges := randomTree(rng, n)
		tests = append(tests, testCase{n: n, m: m, edges: edges})
		totalN += n
		if len(tests) > 200 {
			break
		}
	}

	return tests
}

func edgesFromPairs(n int, pairs [][2]int) [][2]int {
	if len(pairs) != n-1 {
		return pairs
	}
	return pairs
}

func randomTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	return edges
}
