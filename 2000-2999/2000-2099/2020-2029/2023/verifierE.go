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

const referenceSource = "2000-2999/2000-2099/2020-2029/2023/2023E.go"

type testCase struct {
	name  string
	input string
	cases int
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
	refBin, cleanup, err := buildReferenceBinary(refPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}

		for i := 0; i < tc.cases; i++ {
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\n", idx+1, tc.name, i+1, refVals[i], candVals[i])
				fmt.Fprintln(os.Stderr, previewInput(tc.input))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func candidatePathFromArgs() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("usage: go run verifierE.go /path/to/binary-or-source")
	}
	return os.Args[1], nil
}

func referencePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Join(filepath.Dir(file), "2023E.go")
	}
	return referenceSource
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2023E")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2023e")
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
	tmpDir, err := os.MkdirTemp("", "candidate2023E")
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

func parseOutputs(out string, cases int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != cases {
		return nil, fmt.Errorf("expected %d integers, got %d", cases, len(tokens))
	}
	res := make([]int64, cases)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	var tests []testCase

	sample := strings.Builder{}
	sample.WriteString("5\n")
	sample.WriteString("4\n1 2\n2 3\n3 4\n")
	sample.WriteString("2\n1 2\n")
	sample.WriteString("4\n1 2\n1 3\n1 4\n")
	sample.WriteString("8\n3 7\n2 4\n1 2\n5 3\n6 1\n3 8\n2 3\n")
	sample.WriteString("6\n2 3\n1 2\n3 6\n1 5\n1 4\n")
	tests = append(tests, testCase{name: "sample-mix", input: sample.String(), cases: 5})

	tests = append(tests, testCase{
		name:  "line-and-star",
		cases: 3,
		input: "3\n" +
			"3\n1 2\n2 3\n" +
			"5\n1 2\n2 3\n3 4\n4 5\n" +
			"7\n1 2\n1 3\n1 4\n1 5\n1 6\n1 7\n",
	})

	tests = append(tests, testCase{
		name:  "small-repeated-center",
		cases: 2,
		input: "2\n" +
			"6\n1 2\n2 3\n2 4\n4 5\n4 6\n" +
			"6\n3 1\n2 3\n3 4\n3 5\n3 6\n",
	})

	tests = append(tests, testCase{
		name:  "single-case-minimal",
		cases: 1,
		input: "1\n2\n1 2\n",
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 120; i++ {
		tests = append(tests, randomTest(rng, i, 20, 6))
	}
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTest(rng, i, 500, 4))
	}
	for i := 0; i < 10; i++ {
		tests = append(tests, randomTest(rng, i, 3000, 3))
	}

	return tests
}

func randomTest(rng *rand.Rand, idx, maxN, maxCases int) testCase {
	t := rng.Intn(maxCases) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN-1) + 2
		edges := randomTree(rng, n)
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return testCase{
		name:  fmt.Sprintf("random-%d-%d", maxN, idx+1),
		input: sb.String(),
		cases: t,
	}
}

func randomTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		if rng.Intn(2) == 0 {
			edges = append(edges, [2]int{u, v})
		} else {
			edges = append(edges, [2]int{v, u})
		}
	}
	return edges
}

func previewInput(in string) string {
	const limit = 500
	if len(in) <= limit {
		return in
	}
	return in[:limit] + "..."
}
