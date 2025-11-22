package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	limitN      = 500
	limitM      = 400000
	limitQ      = 1000
	maxValue    = 1_000_000_000
	maxEdgeCost = 100000
	refBinName  = "./2112F_ref.bin"
)

type edge struct {
	x int
	y int
	z int
}

type query struct {
	k int64
	a []int64
}

type testInput struct {
	name    string
	n       int
	edges   []edge
	queries []query
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinName, "2112F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refBinName, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
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

func buildInput(t testInput) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t.n))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(len(t.edges)))
	sb.WriteByte('\n')
	for _, e := range t.edges {
		sb.WriteString(strconv.Itoa(e.x + 1))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e.y + 1))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e.z))
		sb.WriteByte('\n')
	}
	sb.WriteString(strconv.Itoa(len(t.queries)))
	sb.WriteByte('\n')
	for _, q := range t.queries {
		sb.WriteString(strconv.FormatInt(q.k, 10))
		sb.WriteByte('\n')
		for i, v := range q.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutputs(output string, t testInput) ([]string, error) {
	tokens := strings.Fields(output)
	idx := 0
	lines := make([]string, 0, len(t.queries))

	for qi := 0; qi < len(t.queries); qi++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("missing output for query %d", qi+1)
		}
		var sb strings.Builder
		for sb.Len() < t.n && idx < len(tokens) {
			sb.WriteString(tokens[idx])
			idx++
		}
		if sb.Len() != t.n {
			return nil, fmt.Errorf("query %d: expected %d characters, got %d", qi+1, t.n, sb.Len())
		}
		line := sb.String()
		for j, ch := range line {
			if ch != '0' && ch != '1' {
				return nil, fmt.Errorf("query %d: invalid character %q at position %d", qi+1, ch, j+1)
			}
		}
		lines = append(lines, line)
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra output detected (%d tokens)", len(tokens)-idx)
	}
	return lines, nil
}

func compareOutputs(exp, got []string) error {
	if len(exp) != len(got) {
		return fmt.Errorf("expected %d output lines, got %d", len(exp), len(got))
	}
	for i := range exp {
		if exp[i] != got[i] {
			return fmt.Errorf("query %d mismatch: expected %s, got %s", i+1, exp[i], got[i])
		}
	}
	return nil
}

func sampleInput() testInput {
	return testInput{
		name: "sample",
		n:    4,
		edges: []edge{
			{1 - 1, 0, 10},
			{2 - 1, 1, 5},
			{0, 3, 8},
			{0, 1, 6},
			{2, 0, 17},
		},
		queries: []query{
			{k: 30, a: []int64{20, 0, 15, 5}},
			{k: 10, a: []int64{20, 0, 15, 5}},
			{k: 30, a: []int64{20, 0, 15, 5}},
		},
	}
}

func craftedPathInput() testInput {
	n := 4
	edges := []edge{
		// Direct 3->1 costly, but path 3->2->1 cheaper.
		{0, 2, 20},
		{1, 2, 2},
		{0, 1, 2},
		// Another alternative path vs direct.
		{3, 1, 50},
		{2, 3, 1},
	}
	queries := []query{
		{k: 0, a: []int64{5, 5, 5, 5}},
		{k: 5, a: []int64{7, 1, 10, 0}},
	}
	return testInput{name: "better-paths", n: n, edges: edges, queries: queries}
}

func randomEdges(rng *rand.Rand, n, m int) []edge {
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		x := rng.Intn(n)
		y := rng.Intn(n - 1)
		if y >= x {
			y++
		}
		z := rng.Intn(maxEdgeCost + 1)
		edges[i] = edge{x: x, y: y, z: z}
	}
	return edges
}

func randomQueries(rng *rand.Rand, n, q int) []query {
	queries := make([]query, q)
	for i := 0; i < q; i++ {
		k := rng.Int63n(1_000_000_001)
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Int63n(maxValue + 1)
		}
		queries[i] = query{k: k, a: a}
	}
	return queries
}

func randomInput(rng *rand.Rand, name string, nMin, nMax, qMax int, dense bool) testInput {
	n := rng.Intn(nMax-nMin+1) + nMin
	maxEdges := n * (n - 1)
	if maxEdges > limitM {
		maxEdges = limitM
	}
	m := rng.Intn(maxEdges) + 1
	if dense {
		m = maxEdges
	}
	q := rng.Intn(qMax) + 1
	if q > limitQ {
		q = limitQ
	}
	return testInput{
		name:    name,
		n:       n,
		edges:   randomEdges(rng, n, m),
		queries: randomQueries(rng, n, q),
	}
}

func largeStress(rng *rand.Rand) testInput {
	n := 150
	m := 20000
	q := 120
	return testInput{
		name:    "stress",
		n:       n,
		edges:   randomEdges(rng, n, m),
		queries: randomQueries(rng, n, q),
	}
}

func buildTests() []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []testInput{
		sampleInput(),
		craftedPathInput(),
		randomInput(rng, "small-random", 2, 8, 5, false),
		randomInput(rng, "medium-random", 20, 60, 20, false),
		randomInput(rng, "dense-random", 25, 40, 10, true),
		randomInput(rng, "large-random", 60, 120, 80, false),
		largeStress(rng),
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for idx, t := range tests {
		input := buildInput(t)

		expRaw, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\n", idx+1, t.name, err)
			os.Exit(1)
		}
		exp, err := parseOutputs(expRaw, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, t.name, err, expRaw)
			os.Exit(1)
		}

		gotRaw, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): runtime error: %v\ninput:\n%s", idx+1, t.name, err, input)
			os.Exit(1)
		}
		got, err := parseOutputs(gotRaw, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): invalid output: %v\ninput:\n%soutput:\n%s\n", idx+1, t.name, err, input, gotRaw)
			os.Exit(1)
		}

		if err := compareOutputs(exp, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, t.name, err, input, expRaw, gotRaw)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
