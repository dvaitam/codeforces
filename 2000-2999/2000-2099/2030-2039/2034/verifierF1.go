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
)

type testCase struct {
	name  string
	input string
}

type scroll struct {
	r int
	b int
}

type khayyamCase struct {
	n int
	m int
	k int
	s []scroll
}

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to locate verifier")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "reference build failed:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d (%s): %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d (%s): %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on case %d (%s)\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", i+1, tc.name, exp, got, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, func(), error) {
	tmp, err := os.CreateTemp("", "cf-2034F1-ref-*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2034F1.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.Remove(tmp.Name())
	}
	return tmp.Name(), cleanup, nil
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
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(2034001))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, noScrollsTest())
	tests = append(tests, smallScrollsTest())
	tests = append(tests, randomPack("random-small", rng, 6, 2000, 12, 30))
	tests = append(tests, randomPack("random-mid", rng, 5, 5000, 40, 60))
	tests = append(tests, randomPack("random-big", rng, 3, 90000, 120, 120))
	tests = append(tests, heavySingleCase(rng))

	return tests
}

func sampleTest() testCase {
	// Directly use the sample from the statement.
	input := "5\n3 4 0\n1 1 1\n1 0\n3 3 2\n1 1\n2 2\n3 3 2\n2 1\n1 2\n10 4 5\n1 0\n8 0\n6 4\n0 2\n7 4\n"
	return testCase{name: "sample", input: input}
}

func noScrollsTest() testCase {
	cases := []khayyamCase{
		{n: 5, m: 7, k: 0},
		{n: 1, m: 1, k: 0},
		{n: 20000, m: 15000, k: 0},
	}
	return packCases("no-scrolls", cases)
}

func smallScrollsTest() testCase {
	cases := []khayyamCase{
		{
			n: 4, m: 3, k: 2,
			s: []scroll{{r: 2, b: 1}, {r: 0, b: 2}},
		},
		{
			n: 6, m: 5, k: 3,
			s: []scroll{{r: 3, b: 1}, {r: 1, b: 4}, {r: 5, b: 0}},
		},
	}
	return packCases("small-scrolls", cases)
}

func randomPack(name string, rng *rand.Rand, t int, maxTotal int, maxK int, _ int) testCase {
	cases := make([]khayyamCase, 0, t)
	for len(cases) < t {
		n := 1 + rng.Intn(maxTotal/2)
		m := 1 + rng.Intn(maxTotal/2)
		if n+m > maxTotal {
			continue
		}
		k := rng.Intn(maxK + 1)
		cs := randomCase(rng, n, m, k)
		cases = append(cases, cs)
	}
	return packCases(name, cases)
}

func heavySingleCase(rng *rand.Rand) testCase {
	// Keep sums within limits while stressing DP size.
	n := 120000 + rng.Intn(5000)
	m := 70000 + rng.Intn(5000)
	k := 120
	cs := randomCase(rng, n, m, k)
	return packCases("heavy-single", []khayyamCase{cs})
}

func randomCase(rng *rand.Rand, n, m, k int) khayyamCase {
	scrolls := make([]scroll, 0, k)
	seen := make(map[int]struct{})
	for len(scrolls) < k {
		r := rng.Intn(n + 1)
		b := rng.Intn(m + 1)
		if r+b == 0 || r+b == n+m {
			continue
		}
		key := r*(m+1) + b
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		scrolls = append(scrolls, scroll{r: r, b: b})
	}
	return khayyamCase{n: n, m: m, k: len(scrolls), s: scrolls}
}

func packCases(name string, cases []khayyamCase) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d %d %d\n", cs.n, cs.m, cs.k)
		for _, sc := range cs.s {
			fmt.Fprintf(&b, "%d %d\n", sc.r, sc.b)
		}
	}
	return testCase{name: name, input: b.String()}
}
