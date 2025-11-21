package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct {
	a, b int
}

type testCase struct {
	input string
	n, m  int
	edges []edge
}

func buildReference() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "27D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func makeTestCase(n int, edges []edge) testCase {
	normalized := make([]edge, len(edges))
	for i, e := range edges {
		a, b := e.a, e.b
		if a == b {
			panic("self-loop edge not allowed")
		}
		if a > b {
			a, b = b, a
		}
		normalized[i] = edge{a: a, b: b}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(normalized))
	for _, e := range normalized {
		fmt.Fprintf(&sb, "%d %d\n", e.a, e.b)
	}
	return testCase{
		input: sb.String(),
		n:     n,
		m:     len(normalized),
		edges: normalized,
	}
}

func deterministicCases() []testCase {
	cases := []testCase{
		makeTestCase(4, []edge{{1, 3}}),
		makeTestCase(5, []edge{{1, 3}, {2, 5}}),
		makeTestCase(6, []edge{{1, 4}, {2, 5}, {3, 6}}),
		makeTestCase(8, []edge{{1, 4}, {2, 5}, {3, 6}, {4, 7}, {5, 8}}),
	}
	n := 30
	edges := allPossibleEdges(n)
	if len(edges) >= 100 {
		cases = append(cases, makeTestCase(n, append([]edge(nil), edges[:100]...)))
	}
	n = 100
	edges = allPossibleEdges(n)
	if len(edges) >= 100 {
		cases = append(cases, makeTestCase(n, append([]edge(nil), edges[:100]...)))
	}
	return cases
}

func allPossibleEdges(n int) []edge {
	var edges []edge
	for a := 1; a <= n; a++ {
		for b := a + 1; b <= n; b++ {
			if b == a+1 {
				continue
			}
			if a == 1 && b == n {
				continue
			}
			edges = append(edges, edge{a: a, b: b})
		}
	}
	return edges
}

func generateRandomCases(rng *rand.Rand, count int) []testCase {
	cases := make([]testCase, 0, count)
	for len(cases) < count {
		n := rng.Intn(97) + 4 // [4, 100]
		pool := allPossibleEdges(n)
		if len(pool) == 0 {
			continue
		}
		rng.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		limit := len(pool)
		if limit > 100 {
			limit = 100
		}
		m := rng.Intn(limit) + 1
		edges := make([]edge, m)
		copy(edges, pool[:m])
		cases = append(cases, makeTestCase(n, edges))
	}
	return cases
}

func parseOutput(out string, m int) (string, bool, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return "", false, fmt.Errorf("empty output")
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return "", false, fmt.Errorf("empty output")
	}
	if strings.EqualFold(fields[0], "impossible") {
		if len(fields) != 1 {
			return "", false, fmt.Errorf("unexpected data after Impossible")
		}
		return "", true, nil
	}
	joined := strings.Join(fields, "")
	if len(joined) != m {
		return "", false, fmt.Errorf("expected %d characters, got %d", m, len(joined))
	}
	var sb strings.Builder
	for i, ch := range joined {
		switch ch {
		case 'i', 'I':
			sb.WriteByte('i')
		case 'o', 'O':
			sb.WriteByte('o')
		default:
			return "", false, fmt.Errorf("invalid character %q at position %d", string(ch), i+1)
		}
	}
	return sb.String(), false, nil
}

func inRange(x int, e edge) bool {
	return x > e.a && x < e.b
}

func cross(e1, e2 edge) bool {
	return (inRange(e1.a, e2) || inRange(e1.b, e2)) &&
		(inRange(e2.a, e1) || inRange(e2.b, e1))
}

func validateAssignment(edges []edge, assign string) error {
	for i := 0; i < len(edges); i++ {
		for j := i + 1; j < len(edges); j++ {
			if cross(edges[i], edges[j]) && assign[i] == assign[j] {
				return fmt.Errorf("roads %d and %d cross but are both %c", i+1, j+1, assign[i])
			}
		}
	}
	return nil
}

func verifyCase(candidate, ref string, tc testCase) error {
	refOut, err := runProgram(ref, tc.input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	_, refImpossible, err := parseOutput(refOut, tc.m)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return err
	}
	assign, candImpossible, err := parseOutput(candOut, tc.m)
	if err != nil {
		return err
	}

	if refImpossible {
		if candImpossible {
			return nil
		}
		return fmt.Errorf("expected Impossible, got %q", candOut)
	}
	if candImpossible {
		return fmt.Errorf("candidate reported Impossible but a valid solution exists")
	}
	return validateAssignment(tc.edges, assign)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, generateRandomCases(rng, 200)...)

	for idx, tc := range tests {
		if err := verifyCase(candidate, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
