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

type testCase struct {
	name  string
	input string
}

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine verifier location")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expected := strings.TrimSpace(refOut)

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got := strings.TrimSpace(candOut)

		if expected != got {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected:\n%s\nbut got:\n%s\ninput:\n%s", idx+1, tc.name, expected, got, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1184E2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1184E2.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, manualCase("simple_tree", 3, [][]int{
		{1, 2, 1},
		{2, 3, 2},
	}))
	tests = append(tests, manualCase("line_plus_extra", 4, [][]int{
		{1, 2, 1},
		{2, 3, 2},
		{3, 4, 3},
		{1, 3, 4},
		{2, 4, 5},
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomCase(rng, "random_small", 6, 10))
	tests = append(tests, randomCase(rng, "random_medium", 200, 500))
	tests = append(tests, randomCase(rng, "random_dense", 500, 2000))
	tests = append(tests, randomCase(rng, "random_big", 2000, 4000))
	tests = append(tests, randomCase(rng, "random_huge", 10000, 20000))

	return tests
}

func manualCase(name string, n int, edges [][]int) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d %d\n", e[0], e[1], e[2])
	}
	return testCase{name: name, input: b.String()}
}

func randomCase(rng *rand.Rand, name string, n, m int) testCase {
	if m < n-1 {
		m = n - 1
	}
	type pair struct{ u, v int }
	edges := make([][3]int, 0, m)
	used := make(map[pair]struct{})
	// build random tree first
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		w := len(edges) + 1
		edges = append(edges, [3]int{u, v, w})
		used[pair{min(u, v), max(u, v)}] = struct{}{}
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		p := pair{min(u, v), max(u, v)}
		if _, ok := used[p]; ok {
			continue
		}
		used[p] = struct{}{}
		w := len(edges) + 1
		edges = append(edges, [3]int{u, v, w})
	}
	// shuffle edges
	for i := len(edges) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		edges[i], edges[j] = edges[j], edges[i]
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d %d\n", e[0], e[1], e[2])
	}
	return testCase{name: name, input: b.String()}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
