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

const refSource = "1000-1999/1700-1799/1770-1779/1774/1774E.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswer(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if refAns != candAns {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, refAns, candAns, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1774E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
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
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseAnswer(out string) (int64, error) {
	reader := strings.NewReader(strings.TrimSpace(out))
	var val int64
	if _, err := fmt.Fscan(reader, &val); err != nil {
		return 0, fmt.Errorf("unable to read integer: %w", err)
	}
	// ensure no extra tokens
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return 0, fmt.Errorf("extra output detected")
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildTest("two-nodes-basic", 2, 2, [][2]int{{1, 2}}, []int{1}, []int{2}),
		buildTest("chain-deep-anchor", 5, 2, [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}, []int{5}, []int{2}),
		buildTest("shared-ancestors", 6, 3, [][2]int{{1, 2}, {2, 3}, {3, 4}, {2, 5}, {5, 6}}, []int{4, 6}, []int{3}),
	}

	rng := rand.New(rand.NewSource(17740115))
	for i := 0; i < 70; i++ {
		tests = append(tests, randomTest(fmt.Sprintf("random-%d", i+1), rng))
	}
	return tests
}

func buildTest(name string, n, d int, edges [][2]int, a, b []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, d)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d", len(a))
	for _, v := range a {
		fmt.Fprintf(&sb, " %d", v)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d", len(b))
	for _, v := range b {
		fmt.Fprintf(&sb, " %d", v)
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func randomTest(name string, rng *rand.Rand) testCase {
	n := rng.Intn(95) + 6 // n in [6,100]
	d := rng.Intn(n-1) + 2
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}

	m1 := rng.Intn(n) + 1
	m2 := rng.Intn(n) + 1
	a := pickDistinct(rng, n, m1)
	b := pickDistinct(rng, n, m2)
	return buildTest(name, n, d, edges, a, b)
}

func pickDistinct(rng *rand.Rand, n, cnt int) []int {
	seen := make([]bool, n+1)
	res := make([]int, 0, cnt)
	for len(res) < cnt {
		v := rng.Intn(n) + 1
		if seen[v] {
			continue
		}
		seen[v] = true
		res = append(res, v)
	}
	return res
}
