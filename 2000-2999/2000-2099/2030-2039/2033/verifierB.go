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

const refSource = "2000-2999/2000-2099/2030-2039/2033/2033B.go"

type testCase struct {
	name  string
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
		}
		exp, err := parseInts(refOut, tc.t)
		if err != nil {
			fail("failed to parse reference output on test %d (%s): %v\nOutput:\n%s", idx+1, tc.name, err, refOut)
		}

		candOut, err := runProgramCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
		}
		got, err := parseInts(candOut, tc.t)
		if err != nil {
			fail("failed to parse candidate output on test %d (%s): %v\nOutput:\n%s", idx+1, tc.name, err, candOut)
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
	tmp, err := os.CreateTemp("", "2033B-ref-*")
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
		return "", fmt.Errorf("go build failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runProgramCandidate(path, input string) (string, error) {
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

func parseInts(out string, expected int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(tokens))
	}
	res := make([]int64, expected)
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
	add := func(name string, datasets [][][]int64) {
		input, t := formatInput(datasets)
		tests = append(tests, testCase{name: name, input: input, t: t})
	}

	add("samples", [][][]int64{
		{
			{1, 1},
			{2, -1},
		},
		{
			{2, 3, 0},
			{3, 1, 2},
			{3, -2, 1},
		},
		{
			{-1, 0, -1},
			{0, 0, -1},
			{-1, -1, 0},
		},
	})

	add("single-cell", [][][]int64{
		{{-5}},
		{{0}},
		{{7}},
	})

	add("all-negative", [][][]int64{
		makeFilledMatrix(3, -1)[0],
		makeFilledMatrix(4, -10)[0],
	})

	tests = append(tests, randomTests("random-small", 5, 5))
	tests = append(tests, randomTests("random-large", 3, 200))
	return tests
}

func formatInput(mats [][][]int64) (string, int) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(mats))
	for _, mat := range mats {
		n := len(mat)
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.FormatInt(mat[i][j], 10))
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String(), len(mats)
}

func makeFilledMatrix(n int, val int64) [][][]int64 {
	mat := make([][]int64, n)
	for i := range mat {
		mat[i] = make([]int64, n)
		for j := range mat[i] {
			mat[i][j] = val
		}
	}
	return [][][]int64{mat}
}

func randomTests(name string, cases, maxN int) testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	mats := make([][][]int64, cases)
	for c := 0; c < cases; c++ {
		n := rng.Intn(maxN) + 1
		if totalN+n > 1000 {
			n = 1
		}
		totalN += n
		mat := make([][]int64, n)
		for i := 0; i < n; i++ {
			mat[i] = make([]int64, n)
			for j := 0; j < n; j++ {
				mat[i][j] = int64(rng.Intn(200001) - 100000)
			}
		}
		mats[c] = mat
	}
	input, t := formatInput(mats)
	return testCase{name: name, input: input, t: t}
}
