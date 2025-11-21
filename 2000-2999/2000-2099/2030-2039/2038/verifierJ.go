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

const refSource = "2000-2999/2000-2099/2030-2039/2038/2038J.go"

type testCase struct {
	name  string
	input string
	buses int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/candidate")
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
		exp, err := parseAnswers(refOut, tc.buses)
		if err != nil {
			fail("failed to parse reference output on test %d (%s): %v\nOutput:\n%s", idx+1, tc.name, err, refOut)
		}

		candOut, err := runProgramCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
		}
		got, err := parseAnswers(candOut, tc.buses)
		if err != nil {
			fail("failed to parse candidate output on test %d (%s): %v\nOutput:\n%s", idx+1, tc.name, err, candOut)
		}
		for i := 0; i < tc.buses; i++ {
			if got[i] != exp[i] {
				fail("wrong answer on test %d (%s) bus %d: expected %s got %s\nInput:\n%s",
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
	tmp, err := os.CreateTemp("", "2038J-ref-*")
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
	var tests []testCase
	add := func(name string, input string, buses int) {
		tests = append(tests, testCase{name: name, input: input, buses: buses})
	}

	add("sample", sampleInput(), 6)
	add("single-bus", manualInput([]string{"B 1"}), 1)
	add("only-people-before-bus", manualInput([]string{"P 5", "P 3", "B 10"}), 1)
	add("overcrowded", manualInput([]string{"P 5", "B 3", "B 10"}), 2)
	tests = append(tests, randomTest("random-small", 20, 10))
	tests = append(tests, randomTest("random-large", 1000, 1000000))
	return tests
}

func sampleInput() string {
	lines := []string{
		"10",
		"P 2",
		"P 5",
		"B 8",
		"P 14",
		"B 5",
		"B 9",
		"B 3",
		"P 2",
		"B 1",
		"B 2",
	}
	return strings.Join(lines, "\n") + "\n"
}

func manualInput(events []string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(events))
	for _, ev := range events {
		sb.WriteString(ev)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomTest(name string, n int, maxVal int) testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var events []string
	busCount := 0
	for i := 0; i < n; i++ {
		if rng.Intn(3) == 0 {
			events = append(events, fmt.Sprintf("B %d", rng.Intn(maxVal)+1))
			busCount++
		} else {
			val := rng.Intn(maxVal) + 1
			events = append(events, fmt.Sprintf("P %d", val))
		}
	}
	if busCount == 0 {
		events = append(events, fmt.Sprintf("B %d", rng.Intn(maxVal)+1))
		busCount++
		n++
	}
	input := manualInput(events)
	return testCase{name: name, input: input, buses: busCount}
}
