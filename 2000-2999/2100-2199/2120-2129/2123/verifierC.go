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

const refSource = "2123C.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refTokens, err := parseOutputs(refOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candTokens, err := parseOutputs(candOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.outputs; i++ {
			if refTokens[i] != candTokens[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) at case %d: expected %q got %q\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refTokens[i], candTokens[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "oracle-2123C-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleC")

	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
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

func runProgram(bin, input string) (string, error) {
	cmd := commandFor(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, expected int) ([]string, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d lines, got %d", expected, len(lines))
	}
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return lines, nil
}

func buildTests() []testCase {
	tests := []testCase{
		sampleTests(),
		smallPermutations(),
		edgeCases(),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 60; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleTests() testCase {
	input := `3
6
1 3 5 4 7 2
4
13 10 12 20
7
1 2 3 4 5 6 7
`
	return testCase{name: "sample_like", input: input, outputs: 3}
}

func smallPermutations() testCase {
	var sb strings.Builder
	type arrCase struct {
		a []int
	}
	var cases []arrCase
	perms := [][]int{
		{1, 2},
		{2, 1},
		{1, 2, 3},
		{1, 3, 2},
		{2, 1, 3},
		{2, 3, 1},
		{3, 1, 2},
		{3, 2, 1},
	}
	for _, p := range perms {
		cases = append(cases, arrCase{a: p})
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		fmt.Fprintf(&sb, "%d\n", len(cs.a))
		for i, v := range cs.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: "small_perms", input: sb.String(), outputs: len(cases)}
}

func edgeCases() testCase {
	input := `3
2
1 1000000
2
1000000 1
5
100 1 50 2 75
`
	return testCase{name: "edges", input: input, outputs: 3}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(15) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(12) + 2
		arr := randPerm(rng, n)
		fmt.Fprintf(&sb, "%d\n", n)
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		outputs: t,
	}
}

func randPerm(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	// random big shift to mix values
	if rng.Intn(3) == 0 {
		shift := rng.Intn(1_000_000)
		for i := 0; i < n; i++ {
			arr[i] += shift
		}
	}
	return arr
}
