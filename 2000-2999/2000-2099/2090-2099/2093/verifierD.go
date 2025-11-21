package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const refSource = "2093D.go"

type query struct {
	typ string
	x   int64
	y   int64
	d   int64
}

type testCase struct {
	n       int
	queries []query
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)
	expectedLens := outputLengths(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, expectedLens)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, expectedLens)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "output line count mismatch: expected %d got %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}

	for idx := range refAns {
		exp := refAns[idx]
		got := candAns[idx]
		if len(exp) != len(got) {
			tcIdx, qIdx, q := locateQuery(tests, idx)
			fmt.Fprintf(os.Stderr, "line %d length mismatch for test %d query %d (type %s): expected %d values got %d\n", idx+1, tcIdx+1, qIdx+1, q.typ, len(exp), len(got))
			os.Exit(1)
		}
		for j := range exp {
			if exp[j] != got[j] {
				tcIdx, qIdx, q := locateQuery(tests, idx)
				fmt.Fprintf(os.Stderr, "mismatch at line %d value %d: expected %d got %d (test %d query %d type %s", idx+1, j+1, exp[j], got[j], tcIdx+1, qIdx+1, q.typ)
				if q.typ == "->" {
					fmt.Fprintf(os.Stderr, " x=%d y=%d)\n", q.x, q.y)
				} else {
					fmt.Fprintf(os.Stderr, " d=%d)\n", q.d)
				}
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	refPath, err := referencePath()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "ref_2093D_*.bin")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func referencePath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to locate verifier path")
	}
	dir := filepath.Dir(file)
	return filepath.Join(dir, refSource), nil
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

func buildTests() []testCase {
	var tests []testCase
	add := func(tc testCase) {
		tests = append(tests, tc)
	}

	// Simple deterministic cases
	add(testCase{n: 1, queries: []query{{typ: "->", x: 1, y: 1}, {typ: "->", x: 2, y: 2}, {typ: "<-", d: 4}}})
	add(testCase{n: 2, queries: []query{{typ: "->", x: 3, y: 2}, {typ: "<-", d: 7}, {typ: "<-", d: 10}}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(tests) < 10 {
		n := rng.Intn(10) + 1 // keep sizes reasonable
		size := int64(1) << n
		maxD := int64(1) << (2 * n)
		qCount := rng.Intn(15) + 5
		qs := make([]query, 0, qCount)
		for i := 0; i < qCount; i++ {
			if rng.Intn(2) == 0 {
				x := rng.Int63n(size) + 1
				y := rng.Int63n(size) + 1
				qs = append(qs, query{typ: "->", x: x, y: y})
			} else {
				d := rng.Int63n(maxD) + 1
				qs = append(qs, query{typ: "<-", d: d})
			}
		}
		add(testCase{n: n, queries: qs})
	}

	// One larger n to cover depth
	add(testCase{
		n: 15,
		queries: []query{
			{typ: "->", x: 1, y: 1},
			{typ: "<-", d: 1},
			{typ: "->", x: 1 << 15, y: 1 << 15},
			{typ: "<-", d: (1 << 30)},
		},
	})

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n%d\n", tc.n, len(tc.queries)))
		for _, q := range tc.queries {
			if q.typ == "->" {
				sb.WriteString(fmt.Sprintf("-> %d %d\n", q.x, q.y))
			} else {
				sb.WriteString(fmt.Sprintf("<- %d\n", q.d))
			}
		}
	}
	return sb.String()
}

func outputLengths(tests []testCase) []int {
	var lens []int
	for _, tc := range tests {
		for _, q := range tc.queries {
			if q.typ == "->" {
				lens = append(lens, 1)
			} else {
				lens = append(lens, 2)
			}
		}
	}
	return lens
}

func parseOutput(out string, lens []int) ([][]int64, error) {
	tokens := strings.Fields(out)
	expectedTokens := 0
	for _, l := range lens {
		expectedTokens += l
	}
	if len(tokens) != expectedTokens {
		return nil, fmt.Errorf("expected %d integers, got %d", expectedTokens, len(tokens))
	}
	res := make([][]int64, len(lens))
	pos := 0
	for i, l := range lens {
		line := make([]int64, l)
		for j := 0; j < l; j++ {
			v, err := strconv.ParseInt(tokens[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", tokens[pos])
			}
			line[j] = v
			pos++
		}
		res[i] = line
	}
	return res, nil
}

func locateQuery(tests []testCase, globalIdx int) (int, int, query) {
	idx := globalIdx
	for ti, tc := range tests {
		if idx < len(tc.queries) {
			return ti, idx, tc.queries[idx]
		}
		idx -= len(tc.queries)
	}
	return -1, -1, query{}
}
