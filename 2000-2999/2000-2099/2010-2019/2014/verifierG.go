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

type entry struct {
	day int
	amt int
}

type milkCase struct {
	m       int
	k       int
	entries []entry
}

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot locate verifier path")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "cf-2014G-ref-*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2014G.go")
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
	rng := rand.New(rand.NewSource(2014007))
	var tests []testCase

	tests = append(tests, simpleTest())
	tests = append(tests, expiryEdgeTest())
	tests = append(tests, randomPack("random-small", rng, 8, 40, 80, 50, 1000, 8))
	tests = append(tests, randomPack("random-medium", rng, 10, 400, 200, 20, 5000, 20))
	tests = append(tests, randomPack("wide-gaps", rng, 5, 500, 300, 250, 5000, 15))
	tests = append(tests, randomPack("heavy", rng, 3, 60000, 100000, 5, 1000000, 100000))

	return tests
}

func simpleTest() testCase {
	cases := []milkCase{
		{
			m: 1, k: 3,
			entries: []entry{{day: 1, amt: 5}},
		},
		{
			m: 3, k: 2,
			entries: []entry{
				{day: 1, amt: 4},
				{day: 2, amt: 2},
				{day: 4, amt: 5},
			},
		},
	}
	return packCases("simple", cases)
}

func expiryEdgeTest() testCase {
	cases := []milkCase{
		{
			m: 5, k: 1,
			entries: []entry{
				{day: 1, amt: 4},
				{day: 2, amt: 7},
				{day: 3, amt: 9},
				{day: 5, amt: 3},
				{day: 8, amt: 5},
			},
		},
		{
			m: 7, k: 3,
			entries: []entry{
				{day: 2, amt: 10},
				{day: 4, amt: 12},
				{day: 6, amt: 11},
				{day: 7, amt: 14},
				{day: 10, amt: 6},
			},
		},
	}
	return packCases("expiry-edge", cases)
}

func randomPack(name string, rng *rand.Rand, t, maxN, mMax, maxGap, amtMax, kMax int) testCase {
	cases := make([]milkCase, 0, t)
	for i := 0; i < t; i++ {
		n := 1 + rng.Intn(maxN)
		cases = append(cases, randomCase(rng, n, mMax, kMax, maxGap, amtMax))
	}
	return packCases(name, cases)
}

func randomCase(rng *rand.Rand, n, mMax, kMax, maxGap, amtMax int) milkCase {
	m := 1 + rng.Intn(mMax)
	k := 1 + rng.Intn(kMax)
	entries := make([]entry, 0, n)
	day := 1 + rng.Intn(3)
	for i := 0; i < n; i++ {
		day += rng.Intn(maxGap) + 1
		if day > 1_000_000 {
			day = 1_000_000 - rng.Intn(maxGap) - 1
		}
		entries = append(entries, entry{
			day: day,
			amt: 1 + rng.Intn(amtMax),
		})
	}
	return milkCase{m: m, k: k, entries: entries}
}

func packCases(name string, cases []milkCase) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d %d %d\n", len(cs.entries), cs.m, cs.k)
		for _, e := range cs.entries {
			fmt.Fprintf(&b, "%d %d\n", e.day, e.amt)
		}
	}
	return testCase{name: name, input: b.String()}
}
