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
	"time"
)

const refSource2106G1 = "2000-2999/2100-2199/2100-2109/2106/2106G1.go"

type testCase struct {
	n      int
	root   int
	values []int
	edges  [][2]int
}

type namedCase struct {
	name string
	tc   testCase
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource2106G1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runBinary(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	expect, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "internal error: cannot parse reference output: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runBinary(candBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	got, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}

	if err := compareAnswers(expect, got); err != nil {
		fmt.Fprintf(os.Stderr, "mismatch: %v\ninput:\n%sreference:\n%v\ncandidate:\n%v\n", err, input, expect, got)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier2106G1-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []namedCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return []namedCase{
		{name: "tiny_star", tc: starCase(2)},
		{name: "small_star", tc: starCase(5)},
		{name: "alternating", tc: fixedValuesCase([]int{1, -1, 1, -1, 1, -1})},
		{name: "random_small", tc: randomCase(rng, 10)},
		{name: "random_medium", tc: randomCase(rng, 100)},
		{name: "random_max", tc: randomCase(rng, 1000)},
	}
}

func starCase(n int) testCase {
	values := make([]int, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			values[i] = 1
		} else {
			values[i] = -1
		}
	}
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		edges = append(edges, [2]int{1, v})
	}
	return testCase{n: n, root: 1, values: values, edges: edges}
}

func fixedValuesCase(vals []int) testCase {
	n := len(vals)
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		edges = append(edges, [2]int{1, v})
	}
	return testCase{n: n, root: 1, values: append([]int(nil), vals...), edges: edges}
}

func randomCase(rng *rand.Rand, n int) testCase {
	values := make([]int, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			values[i] = -1
		} else {
			values[i] = 1
		}
	}
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		edges = append(edges, [2]int{1, v})
	}
	root := rng.Intn(n) + 1
	return testCase{n: n, root: root, values: values, edges: edges}
}

func buildInput(tests []namedCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for idx, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.tc.n, tc.tc.root)
		for i, v := range tc.tc.values {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, e := range tc.tc.edges {
			fmt.Fprintf(&sb, "%d %d", e[0], e[1])
			if i+1 != len(tc.tc.edges) {
				sb.WriteByte('\n')
			}
		}
		if idx+1 != len(tests) {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string, tests []namedCase) ([][]int, error) {
	fields := strings.Fields(out)
	want := 0
	for _, tc := range tests {
		want += tc.tc.n
	}
	if len(fields) != want {
		return nil, fmt.Errorf("expected %d numbers, got %d", want, len(fields))
	}
	ans := make([][]int, len(tests))
	pos := 0
	for i, tc := range tests {
		ans[i] = make([]int, tc.tc.n)
		for j := 0; j < tc.tc.n; j++ {
			val, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[pos])
			}
			if val != 1 && val != -1 {
				return nil, fmt.Errorf("value %d is not ±1", val)
			}
			ans[i][j] = val
			pos++
		}
	}
	return ans, nil
}

func compareAnswers(expect, got [][]int) error {
	if len(expect) != len(got) {
		return fmt.Errorf("expected %d testcases, got %d", len(expect), len(got))
	}
	for i := range expect {
		if len(expect[i]) != len(got[i]) {
			return fmt.Errorf("test %d: expected %d values, got %d", i+1, len(expect[i]), len(got[i]))
		}
		for j := range expect[i] {
			if expect[i][j] != got[i][j] {
				return fmt.Errorf("test %d node %d mismatch: expected %d got %d", i+1, j+1, expect[i][j], got[i][j])
			}
		}
	}
	return nil
}
