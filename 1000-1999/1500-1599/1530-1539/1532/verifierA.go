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

const refSource = "./1532A.go"

type testCase struct {
	name  string
	input string
	t     int
}

type pair struct {
	a int
	b int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
		}
		exp, err := parseAnswers(refOut, tc.t)
		if err != nil {
			fail("failed to parse reference output on test %d (%s): %v\nOutput:\n%s", idx+1, tc.name, err, refOut)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
		}
		got, err := parseAnswers(candOut, tc.t)
		if err != nil {
			fail("failed to parse candidate output on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
		}
		for i := 0; i < tc.t; i++ {
			if got[i] != exp[i] {
				fail("wrong answer on test %d (%s) case %d: expected %d got %d\nInput:\n%s", idx+1, tc.name, i+1, exp[i], got[i], tc.input)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1532A-ref-*")
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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
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

func parseAnswers(out string, expected int) ([]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(tokens))
	}
	res := make([]int, expected)
	for i, tok := range tokens {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	var tests []testCase
	tests = append(tests, makeCase("basic-positives", []pair{{1, 5}, {3, 14}, {7, 8}, {10, 0}}))
	tests = append(tests, makeCase("including-negatives", []pair{{-5, -7}, {-5, 7}, {12, -12}, {0, 0}, {1000, -1000}}))

	rng := rand.New(rand.NewSource(1532))
	tests = append(tests, randomCase("random-50", rng, 50))
	tests = append(tests, randomCase("random-500", rng, 500))
	tests = append(tests, edgeSweep())
	return tests
}

func makeCase(name string, pairs []pair) testCase {
	return testCase{name: name, input: buildInput(pairs), t: len(pairs)}
}

func randomCase(name string, rng *rand.Rand, cnt int) testCase {
	pairs := make([]pair, cnt)
	for i := 0; i < cnt; i++ {
		pairs[i] = pair{
			a: rng.Intn(2001) - 1000,
			b: rng.Intn(2001) - 1000,
		}
	}
	return makeCase(name, pairs)
}

func edgeSweep() testCase {
	var pairs []pair
	for v := -1000; v <= 1000; v += 200 {
		pairs = append(pairs, pair{a: -1000, b: v})
		pairs = append(pairs, pair{a: 1000, b: v})
		pairs = append(pairs, pair{a: v, b: -1000})
		pairs = append(pairs, pair{a: v, b: 1000})
	}
	return makeCase("edge-sweep", pairs)
}

func buildInput(pairs []pair) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(pairs))
	for _, p := range pairs {
		fmt.Fprintf(&sb, "%d %d\n", p.a, p.b)
	}
	return sb.String()
}
