package main

import (
	"bufio"
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

const referenceSource = "2000-2999/2000-2099/2060-2069/2068/2068A.go"

type testCase struct {
	name  string
	input string
}

type instance struct {
	n, m int
	edge [][2]int
}

func main() {
	candidatePath, err := candidatePathFromArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	candidate, candCleanup, err := prepareCandidateBinary(candidatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer candCleanup()

	refPath := referencePath()
	refBin, refCleanup, err := buildReferenceBinary(refPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer refCleanup()

	tests := buildTests()
	for idx, tc := range tests {
		inst, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error parsing generated input %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		if err := validateOutput(inst, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\nraw output:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		if err := validateOutput(inst, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\nraw output:\n%s\n", idx+1, tc.name, err, candOut)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func candidatePathFromArgs() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("usage: go run verifierA.go /path/to/binary-or-source")
	}
	return os.Args[1], nil
}

func referencePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Join(filepath.Dir(file), "2068A.go")
	}
	return referenceSource
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2068A-ref")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2068a")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func prepareCandidateBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmpDir, err := os.MkdirTemp("", "verifier2068A-cand")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "candidate")
	cmd := exec.Command("go", "build", "-o", bin, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build candidate: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseInput(in string) (instance, error) {
	var inst instance
	r := bufio.NewReader(strings.NewReader(in))
	if _, err := fmt.Fscan(r, &inst.n, &inst.m); err != nil {
		return inst, fmt.Errorf("failed to read n,m: %v", err)
	}
	if inst.n < 2 || inst.n > 50 || inst.m < 0 || inst.m > inst.n*(inst.n-1)/2 {
		return inst, fmt.Errorf("invalid bounds for n or m")
	}
	inst.edge = make([][2]int, inst.m)
	seen := make(map[[2]int]struct{}, inst.m)
	for i := 0; i < inst.m; i++ {
		var a, b int
		if _, err := fmt.Fscan(r, &a, &b); err != nil {
			return inst, fmt.Errorf("failed to read edge %d: %v", i+1, err)
		}
		if a < 1 || a > inst.n || b < 1 || b > inst.n || a == b {
			return inst, fmt.Errorf("edge %d invalid values", i+1)
		}
		key := [2]int{a, b}
		if _, ok := seen[[2]int{b, a}]; ok {
			return inst, fmt.Errorf("duplicate unordered pair on edge %d", i+1)
		}
		seen[key] = struct{}{}
		inst.edge[i] = key
	}
	return inst, nil
}

func validateOutput(inst instance, out string) error {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return fmt.Errorf("empty output")
	}
	ptr := 0
	first := strings.ToLower(tokens[ptr])
	ptr++
	if first != "yes" {
		return fmt.Errorf("first token must be YES/Yes/yes, got %q", tokens[0])
	}
	if ptr >= len(tokens) {
		return fmt.Errorf("missing k")
	}
	k, err := strconv.Atoi(tokens[ptr])
	ptr++
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}
	if k < 1 || k > 50000 {
		return fmt.Errorf("k out of range: %d", k)
	}
	expectedPermTokens := k * inst.n
	if len(tokens)-ptr != expectedPermTokens {
		return fmt.Errorf("expected %d permutation integers, got %d", expectedPermTokens, len(tokens)-ptr)
	}

	pos := make([]int, inst.n+1)
	countAhead := make([]int, len(inst.edge))

	for vote := 0; vote < k; vote++ {
		seen := make([]bool, inst.n+1)
		for i := 0; i < inst.n; i++ {
			val, err := strconv.Atoi(tokens[ptr])
			if err != nil {
				return fmt.Errorf("invalid permutation value %q", tokens[ptr])
			}
			ptr++
			if val < 1 || val > inst.n {
				return fmt.Errorf("vote %d has value out of range: %d", vote+1, val)
			}
			if seen[val] {
				return fmt.Errorf("vote %d repeats candidate %d", vote+1, val)
			}
			seen[val] = true
			pos[val] = i
		}
		for idx, e := range inst.edge {
			if pos[e[0]] < pos[e[1]] {
				countAhead[idx]++
			}
		}
	}

	half := k / 2
	for idx, e := range inst.edge {
		if countAhead[idx] <= half {
			return fmt.Errorf("pair %d: candidate %d before %d only %d times out of %d", idx+1, e[0], e[1], countAhead[idx], k)
		}
	}
	return nil
}

func buildTests() []testCase {
	var tests []testCase

	tests = append(tests, testCase{
		name:  "sample-1",
		input: "2 1\n1 2\n",
	})
	tests = append(tests, testCase{
		name:  "sample-2",
		input: "3 3\n1 2\n2 3\n3 1\n",
	})
	tests = append(tests, testCase{
		name:  "no-edges",
		input: "4 0\n",
	})
	tests = append(tests, testCase{
		name:  "single-node-pair",
		input: "50 1\n50 1\n",
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 120; i++ {
		tests = append(tests, randomTest(rng, i, 8))
	}
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTest(rng, i, 20))
	}
	for i := 0; i < 10; i++ {
		tests = append(tests, randomTest(rng, i, 50))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx, maxN int) testCase {
	n := rng.Intn(maxN-1) + 2
	maxPairs := n * (n - 1) / 2
	m := rng.Intn(maxPairs + 1)
	edges := make(map[[2]int]struct{})
	var edgeList [][2]int
	for len(edgeList) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		key := [2]int{a, b}
		if _, ok := edges[key]; ok {
			continue
		}
		if _, ok := edges[[2]int{b, a}]; ok {
			continue
		}
		edges[key] = struct{}{}
		edgeList = append(edgeList, key)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edgeList)))
	for _, e := range edgeList {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}

	return testCase{
		name:  fmt.Sprintf("random-n%d-%d", maxN, idx+1),
		input: sb.String(),
	}
}

func previewInput(in string) string {
	const limit = 500
	if len(in) <= limit {
		return in
	}
	return in[:limit] + "..."
}
