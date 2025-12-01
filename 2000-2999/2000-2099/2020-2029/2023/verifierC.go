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

const (
	// refSource points to the local reference solution to avoid GOPATH resolution.
	refSource     = "2023C.go"
	totalNLimit   = 200000
	totalEdgesLim = 500000
)

type edge struct {
	u int
	v int
}

type testCase struct {
	n      int
	k      int
	a      []int
	b      []int
	graph1 []edge
	graph2 []edge
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	candidate := os.Args[1]
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
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected := tokenize(refOut)

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	got := tokenize(candOut)

	if len(expected) != len(got) {
		fmt.Fprintf(os.Stderr, "wrong number of tokens: expected %d got %d\n", len(expected), len(got))
		os.Exit(1)
	}
	for i := range expected {
		if !strings.EqualFold(expected[i], got[i]) {
			fmt.Fprintf(os.Stderr, "mismatch at token %d: expected %q got %q\n", i+1, expected[i], got[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2023C-ref-*")
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
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func tokenize(s string) []string {
	return strings.Fields(strings.TrimSpace(s))
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20230203))
	var tests []testCase
	totalN := 0
	totalM1 := 0
	totalM2 := 0

	addTest := func(n, k int, mode string) bool {
		if n < 2 || k < 2 || k > n {
			return false
		}
		edgeCount := computeEdgeCount(n, k)
		if totalN+n > totalNLimit || totalM1+edgeCount > totalEdgesLim || totalM2+edgeCount > totalEdgesLim {
			return false
		}
		g1 := buildGraph(n, k, rng)
		g2 := buildGraph(n, k, rng)
		a, b := buildTypes(n, rng, mode)
		tests = append(tests, testCase{
			n:      n,
			k:      k,
			a:      a,
			b:      b,
			graph1: g1,
			graph2: g2,
		})
		totalN += n
		totalM1 += len(g1)
		totalM2 += len(g2)
		return true
	}

	// Deterministic small tests
	addTest(4, 2, "out_in")
	addTest(5, 3, "in_out")
	addTest(6, 6, "")

	for attempts := 0; attempts < 20000 && totalN < totalNLimit; attempts++ {
		remainingN := totalNLimit - totalN
		if remainingN < 2 {
			break
		}
		maxN := remainingN
		if maxN > 400 {
			maxN = 400
		}
		if maxN < 2 {
			break
		}
		n := rng.Intn(maxN-1) + 2

		var k int
		if n <= 20 {
			k = rng.Intn(n-1) + 2
		} else {
			if rng.Intn(5) == 0 {
				upper := n - 1
				if upper > 10 {
					upper = 10
				}
				if upper < 2 {
					upper = 2
				}
				k = rng.Intn(upper-1) + 2
			} else {
				span := n / 3
				if span < 1 {
					span = 1
				}
				minK := n - span
				if minK < 2 {
					minK = 2
				}
				if minK > n {
					minK = n
				}
				k = minK + rng.Intn(n-minK+1)
			}
		}

		mode := ""
		switch rng.Intn(8) {
		case 0:
			mode = "out_in"
		case 1:
			mode = "in_out"
		default:
			mode = ""
		}

		if addTest(n, k, mode) && totalN >= totalNLimit {
			break
		}
	}

	if len(tests) == 0 {
		// Fallback to minimal test to avoid empty input.
		addTest(2, 2, "out_in")
	}
	return tests
}

func computeEdgeCount(n, k int) int {
	base := n / k
	extra := n % k
	sizes := make([]int, k)
	for i := 0; i < k; i++ {
		sizes[i] = base
		if i < extra {
			sizes[i]++
		}
	}
	count := 0
	for i := 0; i < k; i++ {
		next := (i + 1) % k
		count += sizes[i] * sizes[next]
	}
	return count
}

func buildGraph(n, k int, rng *rand.Rand) []edge {
	groups := make([][]int, k)
	order := rng.Perm(n)
	for idx, v := range order {
		label := idx % k
		groups[label] = append(groups[label], v+1)
	}
	var edges []edge
	for i := 0; i < k; i++ {
		next := (i + 1) % k
		for _, from := range groups[i] {
			for _, to := range groups[next] {
				edges = append(edges, edge{u: from, v: to})
			}
		}
	}
	return edges
}

func buildTypes(n int, rng *rand.Rand, mode string) ([]int, []int) {
	a := make([]int, n)
	b := make([]int, n)
	switch mode {
	case "out_in":
		for i := 0; i < n; i++ {
			a[i] = 1
			b[i] = 0
		}
	case "in_out":
		for i := 0; i < n; i++ {
			a[i] = 0
			b[i] = 1
		}
	default:
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(2)
			b[i] = rng.Intn(2)
		}
	}
	return a, b
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.k)
		for i, val := range tc.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", val)
		}
		b.WriteByte('\n')
		fmt.Fprintf(&b, "%d\n", len(tc.graph1))
		for _, e := range tc.graph1 {
			fmt.Fprintf(&b, "%d %d\n", e.u, e.v)
		}
		for i, val := range tc.b {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", val)
		}
		b.WriteByte('\n')
		fmt.Fprintf(&b, "%d\n", len(tc.graph2))
		for _, e := range tc.graph2 {
			fmt.Fprintf(&b, "%d %d\n", e.u, e.v)
		}
	}
	return b.String()
}
