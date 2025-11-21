package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "1000-1999/1600-1699/1660-1669/1666/1666D.go"

type testCase struct {
	name  string
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
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
			fail("reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
		}
		exp, err := parseAnswers(refOut, tc.t)
		if err != nil {
			fail("failed to parse reference output on test %d (%s): %v\nOutput:\n%s",
				idx+1, tc.name, err, refOut)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
		}
		got, err := parseAnswers(candOut, tc.t)
		if err != nil {
			fail("failed to parse candidate output on test %d (%s): %v\nInput:\n%sOutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
		}

		for i := 0; i < tc.t; i++ {
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
	tmp, err := os.CreateTemp("", "1666D-ref-*")
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

func parseAnswers(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d tokens, got %d", expected, len(tokens))
	}
	res := make([]string, expected)
	for i, tok := range tokens {
		res[i] = strings.ToUpper(tok)
		if res[i] != "YES" && res[i] != "NO" {
			return nil, fmt.Errorf("invalid answer %q", tok)
		}
	}
	return res, nil
}

func buildTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{
		name: "sample",
		input: "6\n" +
			"DETERMINED TRME\n" +
			"DETERMINED TERM\n" +
			"PSEUDOPSEUDOHYPOPARATHYROIDISM PEPA\n" +
			"DEINSTITUTIONALIZATION DONATION\n" +
			"CONTEST CODE\n" +
			"SOLUTION SOLUTION\n",
		t: 6,
	})
	tests = append(tests, testCase{
		name: "identity",
		input: "3\n" +
			"A A\n" +
			"ABC ABC\n" +
			"XYZ XYZ\n",
		t: 3,
	})
	tests = append(tests, testCase{
		name: "impossible-short",
		input: "2\n" +
			"AB BA\n" +
			"A B\n",
		t: 2,
	})

	rng := rand.New(rand.NewSource(1666))
	tests = append(tests, randomCase("random-20", rng, 20))
	tests = append(tests, randomCase("random-200", rng, 200))
	return tests
}

func randomCase(name string, rng *rand.Rand, cases int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", cases)
	for i := 0; i < cases; i++ {
		s := randomWord(rng, rng.Intn(30)+1)
		t := mutateWord(rng, s)
		fmt.Fprintf(&sb, "%s %s\n", s, t)
	}
	return testCase{name: name, input: sb.String(), t: cases}
}

func randomWord(rng *rand.Rand, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('A' + rng.Intn(26))
	}
	return string(b)
}

func mutateWord(rng *rand.Rand, s string) string {
	switch rng.Intn(3) {
	case 0:
		return s
	case 1:
		// delete random subset preserving order
		var sb strings.Builder
		for i := 0; i < len(s); i++ {
			if rng.Intn(2) == 0 {
				continue
			}
			sb.WriteByte(s[i])
		}
		res := sb.String()
		if res == "" {
			return string(s[len(s)-1])
		}
		return res
	default:
		// random word unrelated
		return randomWord(rng, rng.Intn(len(s))+1)
	}
}
