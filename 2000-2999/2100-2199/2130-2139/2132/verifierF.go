package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const refSource = "2132F.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
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
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) at line %d: expected %q got %q\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
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

	tmpDir, err := os.MkdirTemp("", "oracle-2132F-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleF")

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
		sampleLike(),
		tinyGraphs(),
		pathAndStar(),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 60; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleLike() testCase {
	input := `2
4 4
1 2
2 3
3 4
1 4
3
1
2
3
5 4
1 2
2 3
3 4
4 5
2
1
5
`
	return testCase{name: "sample_like", input: input, outputs: 5}
}

func tinyGraphs() testCase {
	input := `3
1 0
1
2 1
1 2
2
1
2
3 3
1 2
2 3
1 3
3
1
2
3
`
	return testCase{name: "tiny", input: input, outputs: 6}
}

func pathAndStar() testCase {
	input := `2
6 5
1 2
2 3
3 4
4 5
5 6
3
1
3
6
6 5
1 2
1 3
1 4
1 5
1 6
3
1
4
6
`
	return testCase{name: "path_star", input: input, outputs: 6}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(15) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	outCount := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(30) + 1
		m := rng.Intn(n*(n-1)/2 + 1)
		edges := make(map[[2]int]struct{})
		// ensure connectivity by building a tree first
		for v := 2; v <= n; v++ {
			p := rng.Intn(v-1) + 1
			edges[[2]int{p, v}] = struct{}{}
		}
		// add extra edges
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			a, b := u, v
			if a > b {
				a, b = b, a
			}
			edges[[2]int{a, b}] = struct{}{}
		}
		fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
		for e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		q := rng.Intn(10) + 1
		fmt.Fprintf(&sb, "%d\n", q)
		for j := 0; j < q; j++ {
			c := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "%d\n", c)
		}
		outCount += q
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		outputs: outCount,
	}
}
