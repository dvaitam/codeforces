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

const refSource = "0-999/300-399/330-339/331/331B1.go"

type testCase struct {
	input string
}

type scenario struct {
	n       int
	p       []int
	q       int
	queries []query
}

type query struct {
	t int
	x int
	y int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}

	tests := generateTests()

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}

		if !equalTokens(refOut, gotOut) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc.input, refOut, gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "331B1-ref-*")
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

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
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

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(33103310))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest(buildScenario(2, []int{1, 2}, []query{
		{1, 1, 2},
		{2, 1, 2},
		{1, 1, 2},
	})))

	for i := 0; i < 35; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(6)+1, rng.Intn(10)+1))
	}

	tests = append(tests, limitCase())

	return tests
}

func sampleTest() testCase {
	sc := scenario{
		n: 5,
		p: []int{1, 3, 4, 2, 5},
		q: 6,
		queries: []query{
			{1, 1, 5},
			{1, 3, 4},
			{2, 2, 3},
			{1, 1, 5},
			{2, 1, 5},
			{1, 1, 5},
		},
	}
	return makeTest(sc)
}

func randomCase(rng *rand.Rand, n int, q int) testCase {
	if n < 2 {
		n = 2
	}
	if n > 100 {
		n = 100
	}
	if q < 1 {
		q = 1
	}
	sc := scenario{
		n:       n,
		p:       randPerm(rng, n),
		q:       q,
		queries: make([]query, q),
	}
	for i := 0; i < q; i++ {
		t := rng.Intn(2) + 1
		x := rng.Intn(n-1) + 1
		y := rng.Intn(n-x) + x + 1
		sc.queries[i] = query{t, x, y}
	}
	return makeTest(sc)
}

func limitCase() testCase {
	n := 100
	q := 200
	sc := scenario{
		n:       n,
		p:       randPerm(rand.New(rand.NewSource(123456789)), n),
		q:       q,
		queries: make([]query, q),
	}
	for i := 0; i < q; i++ {
		t := 1
		if i%3 == 0 {
			t = 2
		}
		x := (i % (n - 1)) + 1
		y := n
		sc.queries[i] = query{t, x, y}
	}
	return makeTest(sc)
}

func randPerm(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func buildScenario(n int, p []int, queries []query) scenario {
	return scenario{n: n, p: p, q: len(queries), queries: queries}
}

func makeTest(sc scenario) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, sc.n)
	for i, v := range sc.p {
		if i > 0 {
			fmt.Fprint(&b, " ")
		}
		fmt.Fprint(&b, v)
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, sc.q)
	for _, q := range sc.queries {
		fmt.Fprintf(&b, "%d %d %d\n", q.t, q.x, q.y)
	}
	return testCase{input: b.String()}
}
