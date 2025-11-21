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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %s, got %s\ninput:\n%s", idx+1, tc.name, expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1221G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1221G.go")
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

	tests = append(tests, makeTestCase("no_edges", 4, [][2]int{}))
	tests = append(tests, makeTestCase("single_edge", 2, [][2]int{{1, 2}}))
	tests = append(tests, makeTestCase("triangle", 3, [][2]int{{1, 2}, {2, 3}, {1, 3}}))
	tests = append(tests, makeTestCase("line4", 4, [][2]int{{1, 2}, {2, 3}, {3, 4}}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomCase(rng, "random_sparse", 12, 15))
	tests = append(tests, randomCase(rng, "random_medium", 25, 80))
	tests = append(tests, randomCase(rng, "random_dense", 35, 400))
	tests = append(tests, randomCase(rng, "random_max", 40, 780))

	// additional stress randoms
	for i := 0; i < 5; i++ {
		n := rng.Intn(40) + 1
		maxEdges := n * (n - 1) / 2
		m := rng.Intn(maxEdges + 1)
		tests = append(tests, randomCase(rng, fmt.Sprintf("random_%d", i+1), n, m))
	}

	return tests
}

func makeTestCase(name string, n int, edges [][2]int) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
	}
	return testCase{name: name, input: b.String()}
}

func randomCase(rng *rand.Rand, name string, n, m int) testCase {
	maxEdges := n * (n - 1) / 2
	if m > maxEdges {
		m = maxEdges
	}
	type pair struct{ a, b int }
	used := make(map[pair]struct{})
	edges := make([][2]int, 0, m)
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
		p := pair{a, b}
		if _, ok := used[p]; ok {
			continue
		}
		used[p] = struct{}{}
		edges = append(edges, [2]int{u, v})
	}
	return makeTestCase(name, n, edges)
}
