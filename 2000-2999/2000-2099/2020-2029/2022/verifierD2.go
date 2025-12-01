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

const refSource = "./2022D2.go"

type testCase struct {
	n     int
	roles []int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	expect, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		if expect[i] != got[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d (impostor at %d)\n", i+1, expect[i], got[i], impostorIndex(tc.roles))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2022D2-ref-*")
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
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

func parseOutputs(output string, t int) ([]int, error) {
	tokens := strings.Fields(output)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(tokens))
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[t])
	}
	ans := make([]int, t)
	for i := 0; i < t; i++ {
		val, err := strconv.Atoi(tokens[i])
		if err != nil {
			return nil, fmt.Errorf("token %q is not integer", tokens[i])
		}
		ans[i] = val
	}
	return ans, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2022))
	var tests []testCase
	total := 0
	add := func(tc testCase) {
		if total+tc.n > 100000 {
			return
		}
		tests = append(tests, tc)
		total += tc.n
	}

	add(makeCase([]int{0, 1, 0, -1, 0, 1, 0}))
	add(makeCase([]int{0, 1, -1, 0}))
	add(makeCase([]int{1, 1, 1, -1, 0, 0}))
	add(makeCase([]int{-1, 1, 1}))
	add(makeCase([]int{0, 0, 0, -1, 0}))

	for total < 100000 {
		n := rng.Intn(200) + 3
		roles := make([]int, n)
		pos := rng.Intn(n)
		for i := range roles {
			if i == pos {
				roles[i] = -1
			} else {
				if rng.Intn(2) == 0 {
					roles[i] = 0
				} else {
					roles[i] = 1
				}
			}
		}
		add(makeCase(roles))
		if len(tests) > 500 {
			break
		}
	}

	if total < 100000 {
		roles := make([]int, 1000)
		for i := 0; i < 1000; i++ {
			roles[i] = 1
		}
		roles[0] = -1
		add(makeCase(roles))
	}

	return tests
}

func makeCase(roles []int) testCase {
	cp := make([]int, len(roles))
	copy(cp, roles)
	return testCase{n: len(cp), roles: cp}
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d manual\n", tc.n)
		for i, v := range tc.roles {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func impostorIndex(arr []int) int {
	for i, v := range arr {
		if v == -1 {
			return i + 1
		}
	}
	return -1
}
