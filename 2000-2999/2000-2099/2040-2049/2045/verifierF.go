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

type cell struct {
	r int64
	c int64
	a int64
}

type gameCase struct {
	N   int64
	M   int
	K   int64
	pts []cell
}

var problemDir string

func init() {
	_, file, ok := runtime.Caller(0)
	if !ok {
		panic("unable to locate verifier path")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
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
	tmp, err := os.CreateTemp("", "cf-2045F-ref-*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2045F.go")
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
	rng := rand.New(rand.NewSource(2045007))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, singlePileTests())
	tests = append(tests, layeredRowsTest())
	tests = append(tests, randomPack("random-small", rng, 10, 10, 10))
	tests = append(tests, randomPack("random-mid", rng, 8, 2000, 8000))
	tests = append(tests, randomPack("random-largeK", rng, 5, 50000, 120000))
	tests = append(tests, skewedHeavyCase(rng))

	return tests
}

func sampleTest() testCase {
	input := "3\n2 2 4\n1 1 3\n2 1 2\n100 2 1\n4 1 10\n4 4 10\n10 5 2\n1 1 4\n3 1 2\n4 2 5\n2 2 1\n5 3 4\n"
	return testCase{name: "sample", input: input}
}

func singlePileTests() testCase {
	cases := []gameCase{
		{N: 1, M: 1, K: 1, pts: []cell{{r: 1, c: 1, a: 1}}},
		{N: 5, M: 1, K: 2, pts: []cell{{r: 3, c: 2, a: 7}}},
		{N: 1_000_000_000, M: 1, K: 200000, pts: []cell{{r: 1, c: 1, a: 1}}},
		{N: 1_000_000_000, M: 1, K: 3, pts: []cell{{r: 1_000_000_000, c: 1, a: 5}}},
	}
	return packCases("single-pile", cases)
}

func layeredRowsTest() testCase {
	cases := []gameCase{
		{
			N: 15, M: 4, K: 3,
			pts: []cell{
				{r: 15, c: 1, a: 3},
				{r: 12, c: 5, a: 4},
				{r: 9, c: 2, a: 6},
				{r: 6, c: 3, a: 8},
			},
		},
		{
			N: 50, M: 6, K: 7,
			pts: []cell{
				{r: 5, c: 1, a: 10},
				{r: 10, c: 5, a: 12},
				{r: 20, c: 7, a: 15},
				{r: 30, c: 10, a: 9},
				{r: 40, c: 15, a: 20},
				{r: 45, c: 25, a: 11},
			},
		},
	}
	return packCases("layered-rows", cases)
}

func randomPack(name string, rng *rand.Rand, t int, maxM int, maxN int64) testCase {
	cases := make([]gameCase, 0, t)
	for i := 0; i < t; i++ {
		cases = append(cases, randomCase(rng, maxM, maxN))
	}
	return packCases(name, cases)
}

func randomCase(rng *rand.Rand, maxM int, maxN int64) gameCase {
	M := 1 + rng.Intn(maxM)
	N := int64(1 + rng.Intn(int(maxN)))
	K := int64(1 + rng.Intn(200000))
	pts := make([]cell, 0, M)
	used := make(map[int64]struct{})
	for len(pts) < M {
		r := int64(1 + rng.Intn(int(N)))
		c := int64(1 + rng.Intn(int(r)))
		key := (r << 32) ^ c
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		a := int64(1 + rng.Intn(1_000_000_000))
		pts = append(pts, cell{r: r, c: c, a: a})
	}
	return gameCase{N: N, M: M, K: K, pts: pts}
}

func skewedHeavyCase(rng *rand.Rand) testCase {
	N := int64(1_000_000_000)
	K := int64(200000)
	M := 2000
	pts := make([]cell, 0, M)
	used := make(map[int64]struct{})
	base := N - 5000
	for len(pts) < M {
		r := base + int64(rng.Intn(5000))
		c := int64(1 + rng.Intn(int(r)))
		key := (r << 32) ^ c
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		a := int64(1 + rng.Intn(1_000_000_000))
		pts = append(pts, cell{r: r, c: c, a: a})
	}
	return packCases("skewed-heavy", []gameCase{{N: N, M: M, K: K, pts: pts}})
}

func packCases(name string, cases []gameCase) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d %d %d\n", cs.N, cs.M, cs.K)
		for _, p := range cs.pts {
			fmt.Fprintf(&b, "%d %d %d\n", p.r, p.c, p.a)
		}
	}
	return testCase{name: name, input: b.String()}
}
