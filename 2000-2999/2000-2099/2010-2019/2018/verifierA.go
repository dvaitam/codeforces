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

type caseData struct {
	n int
	k int64
	a []int64
}

type testInput struct {
	text     string
	ansCount int
}

func buildReference() (string, error) {
	refDir := filepath.Join("2000-2999", "2000-2099", "2010-2019", "2018")
	tmp, err := os.CreateTemp("", "ref2018A")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2018A.go")
	cmd.Dir = refDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return tmpPath, nil
}

func commandForPath(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	ans := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		ans[i] = val
	}
	return ans, nil
}

func buildInput(cases []caseData) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.k))
		for i, v := range cs.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return testInput{text: sb.String(), ansCount: len(cases)}
}

func fixedTests() []testInput {
	// Sample input from statement
	sample := "9\n3 1\n3 2 2\n5 4\n2 6 1 2 4\n2 100\n1410065408 1000000000\n10 87\n4 6 6 9 3 10 2 8 7 2\n12 22\n2 22 70 0 11 0 13 0 2 1 2 3\n10 3\n3 3 0 0 0 0 0 0 0 0\n3 2\n3 3 3\n"
	return []testInput{
		buildInput([]caseData{
			{n: 1, k: 0, a: []int64{5}},
			{n: 2, k: 1, a: []int64{1, 0}},
		}),
		{text: sample, ansCount: 9},
	}
}

func randomCase(rng *rand.Rand, maxN int) caseData {
	n := rng.Intn(maxN) + 1
	k := rng.Int63n(1_000_000_000_000_0000) // up to 1e16
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(1_000_000_0000) // up to 1e10
	}
	// ensure sum >= 1
	hasNonZero := false
	for _, v := range a {
		if v > 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		a[rng.Intn(n)] = 1
	}
	return caseData{n: n, k: k, a: a}
}

func largeDeterministicCase(n int) caseData {
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64((i % 5) + 1)
	}
	return caseData{n: n, k: 1_000_000_000_000_0000, a: a}
}

func generateTests() []testInput {
	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Add a near-limit single test covering n ~ 2e5
	tests = append(tests, buildInput([]caseData{largeDeterministicCase(200000)}))

	for len(tests) < 60 {
		numCases := rng.Intn(3) + 1
		cases := make([]caseData, numCases)
		totalN := 0
		for i := 0; i < numCases; i++ {
			maxN := 2000
			if rng.Intn(4) == 0 {
				maxN = 50000
			}
			cs := randomCase(rng, maxN)
			totalN += cs.n
			if totalN > 200000 {
				cs.n = 1
				cs.a = []int64{1}
				cs.k = 0
			}
			cases[i] = cs
		}
		tests = append(tests, buildInput(cases))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, input := range tests {
		expectOut, err := runBinary(ref, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		expect, err := parseOutput(expectOut, input.ansCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, expectOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		got, err := parseOutput(gotOut, input.ansCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for i := range expect {
			if expect[i] != got[i] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d case %d: expected %d got %d\ninput:\n%s\n", idx+1, i+1, expect[i], got[i], preview(input.text))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func preview(s string) string {
	if len(s) <= 400 {
		return s
	}
	return s[:400] + "...\n"
}
