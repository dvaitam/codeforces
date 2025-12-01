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

const refSource = "2000-2999/2100-2199/2130-2139/2131/2131H.go"

type testInput struct {
	raw   string
	cases []caseSpec
}

type caseSpec struct {
	n int
	m int
	a []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}

	// Build reference just to ensure it compiles, but validation is agnostic to its output.
	refBin, err := buildReference()
	if err == nil {
		defer os.Remove(refBin)
	}

	candidate := os.Args[1]
	tests := generateTests()

	for i, ti := range tests {
		expect := analyze(ti)

		got, err := runCandidate(candidate, ti.raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, ti.raw, got)
			os.Exit(1)
		}

		ok, reason := validateOutput(ti, expect.exists, got)
		if !ok {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %s\ninput:\n%s\n", i+1, reason, ti.raw)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2131H-ref-*")
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

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type analysis struct {
	exists []bool
}

func analyze(ti testInput) analysis {
	res := analysis{exists: make([]bool, len(ti.cases))}
	for idx, cs := range ti.cases {
		exists, p1, p2 := findTwoPairs(cs.a)
		if exists {
			_ = p1
			_ = p2
			res.exists[idx] = true
		}
	}
	return res
}

func findTwoPairs(arr []int) (bool, [2]int, [2]int) {
	n := len(arr)
	if n < 4 {
		return false, [2]int{}, [2]int{}
	}

	adj := make([]bool, n*n)
	deg := make([]int, n)
	edges := make([][2]int, 0)
	total := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if gcd(arr[i], arr[j]) == 1 {
				idx := i*n + j
				adj[idx] = true
				deg[i]++
				deg[j]++
				edges = append(edges, [2]int{i, j})
				total++
			}
		}
	}

	if total < 2 {
		return false, [2]int{}, [2]int{}
	}

	for _, e := range edges {
		u, v := e[0], e[1]
		rem := total - deg[u] - deg[v]
		if adj[u*n+v] {
			rem++
		}
		if rem <= 0 {
			continue
		}
		for _, e2 := range edges {
			x, y := e2[0], e2[1]
			if x == u || x == v || y == u || y == v {
				continue
			}
			return true, [2]int{u, v}, [2]int{x, y}
		}
	}
	return false, [2]int{}, [2]int{}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func validateOutput(ti testInput, expectExists []bool, output string) (bool, string) {
	tok := strings.Fields(output)
	pos := 0
	for caseIdx, cs := range ti.cases {
		if pos >= len(tok) {
			return false, "insufficient tokens in output"
		}
		if tok[pos] == "0" {
			if expectExists[caseIdx] {
				return false, fmt.Sprintf("test case %d: output 0 despite existence of a valid quadruple", caseIdx+1)
			}
			pos++
		} else {
			if pos+3 >= len(tok) {
				return false, "not enough indices for quadruple"
			}
			vals := [4]int{}
			for i := 0; i < 4; i++ {
				v, err := atoi(tok[pos+i])
				if err != nil {
					return false, "non-integer output"
				}
				if v < 1 || v > cs.n {
					return false, "index out of bounds"
				}
				vals[i] = v - 1
			}
			pos += 4

			if vals[0] == vals[1] || vals[0] == vals[2] || vals[0] == vals[3] ||
				vals[1] == vals[2] || vals[1] == vals[3] || vals[2] == vals[3] {
				return false, "indices are not distinct"
			}
			if gcd(cs.a[vals[0]], cs.a[vals[1]]) != 1 || gcd(cs.a[vals[2]], cs.a[vals[3]]) != 1 {
				return false, "gcd constraint violated"
			}
			if !expectExists[caseIdx] {
				return false, fmt.Sprintf("test case %d: quadruple provided but none exists", caseIdx+1)
			}
		}
	}
	if pos != len(tok) {
		return false, "extra tokens in output"
	}
	return true, ""
}

func atoi(s string) (int, error) {
	var x int
	_, err := fmt.Sscan(s, &x)
	return x, err
}

func generateTests() []testInput {
	var tests []testInput
	rng := rand.New(rand.NewSource(21312131))

	tests = append(tests, buildInput([]caseSpec{
		{n: 4, m: 15, a: []int{4, 7, 9, 15}},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 4, m: 10, a: []int{6, 10, 14, 22}},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 6, m: 15, a: []int{6, 10, 11, 12, 15, 14}},
	}))

	for i := 0; i < 10; i++ {
		tests = append(tests, randomBatch(rng, 5, 80))
	}
	tests = append(tests, randomBatch(rng, 8, 200))

	return tests
}

func buildInput(cases []caseSpec) testInput {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.m)
		for i, v := range cs.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return testInput{raw: b.String(), cases: cases}
}

func randomBatch(rng *rand.Rand, maxCases int, maxN int) testInput {
	t := rng.Intn(maxCases) + 1
	remaining := 800 // keep n small for exhaustive check
	var cases []caseSpec
	for i := 0; i < t; i++ {
		minRemaining := t - i - 1
		maxAllowed := remaining - minRemaining
		if maxAllowed < 4 {
			break
		}
		if maxAllowed > maxN {
			maxAllowed = maxN
		}
		n := rng.Intn(maxAllowed-3) + 4
		remaining -= n
		m := rng.Intn(500) + n
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(m) + 1
		}
		cases = append(cases, caseSpec{n: n, m: m, a: a})
	}
	if len(cases) == 0 {
		cases = append(cases, caseSpec{n: 4, m: 10, a: []int{2, 3, 4, 5}})
	}
	return buildInput(cases)
}
