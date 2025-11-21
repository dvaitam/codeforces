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

const refSource = "2138C2.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/candidate")
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
		refVals, err := parseOutputs(refOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.outputs; i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) at case %d: expected %q got %q\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
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

	tmpDir, err := os.MkdirTemp("", "oracle-2138C2-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleC2")

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
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return lines, nil
}

func buildTests() []testCase {
	tests := []testCase{
		manualSmall(),
		pathAndStar(),
		mixedTrees(),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func manualSmall() testCase {
	input := `4
2 0
1
3 2
1 1
4 1
1 2 3
5 3
1 1 1 1
`
	return testCase{name: "manual_small", input: input, outputs: 4}
}

func pathAndStar() testCase {
	input := `3
5 2
1 2 3 4
5 3
1 1 1 1
6 4
1 2 3 3 3
`
	return testCase{name: "path_star", input: input, outputs: 3}
}

func mixedTrees() testCase {
	input := `3
7 3
1 1 2 2 4 4
8 4
1 2 3 4 5 6 7
9 5
1 1 2 2 3 3 4 4
`
	return testCase{name: "mixed", input: input, outputs: 3}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(80) + 2
		k := rng.Intn(n + 1)
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for v := 2; v <= n; v++ {
			p := rng.Intn(v-1) + 1
			if v > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(p))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		outputs: t,
	}
}
