package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "2000-2999/2000-2099/2020-2029/2028/2028A.go"

type testCase struct {
	name  string
	input string
	t     int
}

type singleCase struct {
	n int
	a int
	b int
	s string
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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
		}
		exp, err := parseAnswers(refOut, tc.t)
		if err != nil {
			fail("failed to parse reference output on test %d (%s): %v\nOutput:\n%s",
				idx+1, tc.name, err, refOut)
		}

		candOut, err := runProgramCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
		}
		got, err := parseAnswers(candOut, tc.t)
		if err != nil {
			fail("failed to parse candidate output on test %d (%s): %v\nOutput:\n%s",
				idx+1, tc.name, err, candOut)
		}

		for i := range exp {
			if got[i] != exp[i] {
				fail("wrong answer on test %d (%s) case %d: expected %s got %s\nInput:\n%s",
					idx+1, tc.name, i+1, exp[i], got[i], tc.input)
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
	tmp, err := os.CreateTemp("", "2028A-ref-*")
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

func parseAnswers(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(tokens))
	}
	res := make([]string, expected)
	for i, tok := range tokens {
		ans := strings.ToUpper(tok)
		if ans != "YES" && ans != "NO" {
			return nil, fmt.Errorf("invalid answer %q", tok)
		}
		res[i] = ans
	}
	return res, nil
}

func buildTests() []testCase {
	return []testCase{
		makeCase("samples", []singleCase{
			{n: 2, a: 2, b: 2, s: "NE"},
			{n: 3, a: 2, b: 2, s: "NNE"},
			{n: 6, a: 2, b: 1, s: "NNEESW"},
			{n: 6, a: 10, b: 10, s: "NNEESW"},
			{n: 3, a: 4, b: 2, s: "NEE"},
			{n: 4, a: 5, b: 5, s: "NEWS"},
		}),
		makeCase("origin-hit", []singleCase{
			{n: 1, a: 0, b: 0, s: "N"},
			{n: 1, a: 0, b: 0, s: "S"},
		}),
		makeCase("long-runs", []singleCase{
			{n: 10, a: 5, b: 5, s: "NNNNNEEEEE"},
			{n: 10, a: -3, b: -4, s: "SSSSSWWWWN"},
			{n: 10, a: 3, b: -2, s: "EEEESSWWNN"},
		}),
		makeRandomCase("random-1", 10),
		makeRandomCase("random-2", 20),
	}
}

func makeCase(name string, cases []singleCase) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d %d %d\n%s\n", c.n, c.a, c.b, c.s)
	}
	return testCase{name: name, input: sb.String(), t: len(cases)}
}

func makeRandomCase(name string, t int) testCase {
	var cases []singleCase
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 1
		a := rng.Intn(21) - 10
		b := rng.Intn(21) - 10
		var sb strings.Builder
		dirs := []byte{'N', 'E', 'S', 'W'}
		for j := 0; j < n; j++ {
			sb.WriteByte(dirs[rng.Intn(len(dirs))])
		}
		cases = append(cases, singleCase{n: n, a: a, b: b, s: sb.String()})
	}
	return makeCase(name, cases)
}
