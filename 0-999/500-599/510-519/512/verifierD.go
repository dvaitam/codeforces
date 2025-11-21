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

const (
	refSource = "0-999/500-599/510-519/512/512D.go"
	mod       = 1000000009
)

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refVals, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if len(candVals) != len(refVals) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d values got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, len(refVals), len(candVals), tc.input, refOut, candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed at position %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-512D-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref512D.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string) ([]int64, error) {
	line := strings.TrimSpace(out)
	if len(line) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	fields := strings.Fields(line)
	values := make([]int64, len(fields))
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		values[i] = modNormalize(val)
	}
	return values, nil
}

func modNormalize(x int64) int64 {
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "n1_m0", input: "1 0\n0\n"},
		{name: "n2_no_edges", input: "2 0\n0\n"},
		{name: "n2_edge", input: "2 1\n1 2\n"},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		n := rng.Intn(8) + 1
		m := rng.Intn(n*(n-1)/2 + 1)
		tests = append(tests, randomCase(rng, i, n, m))
	}

	// include dense case with n <= 20, random edges
	tests = append(tests, randomCase(rand.New(rand.NewSource(12345)), 1000, 20, 190))
	return tests
}

func randomCase(rng *rand.Rand, idx, n, m int) testCase {
	type edge struct {
		u, v int
	}
	allEdges := make([]edge, 0, n*(n-1)/2)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			allEdges = append(allEdges, edge{i, j})
		}
	}
	// shuffle edges and take first m without duplicates
	rng.Shuffle(len(allEdges), func(i, j int) {
		allEdges[i], allEdges[j] = allEdges[j], allEdges[i]
	})
	selected := allEdges
	if m < len(allEdges) {
		selected = allEdges[:m]
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(selected))
	for _, e := range selected {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	return testCase{
		name:  fmt.Sprintf("random_%d_n%d_m%d", idx+1, n, len(selected)),
		input: sb.String(),
	}
}
