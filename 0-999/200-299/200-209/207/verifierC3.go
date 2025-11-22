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

var problemDir string

func init() {
	_, file, ok := runtime.Caller(0)
	if !ok {
		panic("failed to locate verifier path")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC3.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "cf-207C3-ref-*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "207C3.go")
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
	rng := rand.New(rand.NewSource(2070303))
	var tests []testCase

	tests = append(tests, testCase{"sample", sampleTest()})
	tests = append(tests, testCase{"tree1-only-small", treeOnlyTest(1, 30, rng)})
	tests = append(tests, testCase{"tree2-only-small", treeOnlyTest(2, 30, rng)})
	tests = append(tests, testCase{"alternating", alternatingTest(80, rng)})
	tests = append(tests, testCase{"staged", stagedTest(70, 60, rng)})
	tests = append(tests, testCase{"mirrored-chains", mirroredChainsTest(45)})
	tests = append(tests, testCase{"repeating-letters", repeatingLettersTest(120, rng)})
	tests = append(tests, testCase{"random-1k", randomTest(rng, 1000, 0.55)})
	tests = append(tests, testCase{"random-5k", randomTest(rng, 5000, 0.45)})
	tests = append(tests, testCase{"random-25k", randomTest(rng, 25000, 0.5)})

	return tests
}

func sampleTest() string {
	return "5\n1 1 a\n2 1 a\n1 2 b\n2 1 b\n2 3 a\n"
}

func treeOnlyTest(tree, ops int, rng *rand.Rand) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", ops)
	size1, size2 := 1, 1
	for i := 0; i < ops; i++ {
		if tree == 1 {
			v := 1 + rng.Intn(size1)
			fmt.Fprintf(&b, "1 %d %c\n", v, randomChar(rng))
			size1++
		} else {
			v := 1 + rng.Intn(size2)
			fmt.Fprintf(&b, "2 %d %c\n", v, randomChar(rng))
			size2++
		}
	}
	return b.String()
}

func alternatingTest(ops int, rng *rand.Rand) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", ops)
	size1, size2 := 1, 1
	for i := 0; i < ops; i++ {
		if i%2 == 0 {
			v := 1 + rng.Intn(size1)
			fmt.Fprintf(&b, "1 %d %c\n", v, randomChar(rng))
			size1++
		} else {
			v := 1 + rng.Intn(size2)
			fmt.Fprintf(&b, "2 %d %c\n", v, randomChar(rng))
			size2++
		}
	}
	return b.String()
}

func stagedTest(tree1Ops, tree2Ops int, rng *rand.Rand) string {
	total := tree1Ops + tree2Ops
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", total)
	size1, size2 := 1, 1
	for i := 0; i < tree1Ops; i++ {
		v := 1 + rng.Intn(size1)
		fmt.Fprintf(&b, "1 %d %c\n", v, randomChar(rng))
		size1++
	}
	for i := 0; i < tree2Ops; i++ {
		v := 1 + rng.Intn(size2)
		fmt.Fprintf(&b, "2 %d %c\n", v, randomChar(rng))
		size2++
	}
	return b.String()
}

func mirroredChainsTest(length int) string {
	total := 2 * length
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", total)

	chain := make([]byte, length)
	for i := 0; i < length; i++ {
		chain[i] = byte('a' + (i % 26))
	}

	last1, size1 := 1, 1
	for i := 0; i < length; i++ {
		fmt.Fprintf(&b, "1 %d %c\n", last1, chain[i])
		size1++
		last1 = size1
	}

	last2, size2 := 1, 1
	for i := 0; i < length; i++ {
		fmt.Fprintf(&b, "2 %d %c\n", last2, chain[length-1-i])
		size2++
		last2 = size2
	}

	return b.String()
}

func repeatingLettersTest(ops int, rng *rand.Rand) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", ops)
	size1, size2 := 1, 1
	for i := 0; i < ops; i++ {
		t := 1
		if rng.Intn(3) == 0 {
			t = 2
		}
		letter := byte('a' + rng.Intn(3))
		if t == 1 {
			v := 1 + rng.Intn(size1)
			fmt.Fprintf(&b, "1 %d %c\n", v, letter)
			size1++
		} else {
			v := 1 + rng.Intn(size2)
			fmt.Fprintf(&b, "2 %d %c\n", v, letter)
			size2++
		}
	}
	return b.String()
}

func randomTest(rng *rand.Rand, ops int, tree1Prob float64) string {
	if tree1Prob < 0 {
		tree1Prob = 0
	} else if tree1Prob > 1 {
		tree1Prob = 1
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", ops)
	size1, size2 := 1, 1
	for i := 0; i < ops; i++ {
		t := 1
		if tree1Prob == 0 {
			t = 2
		} else if tree1Prob == 1 {
			t = 1
		} else if rng.Float64() >= tree1Prob {
			t = 2
		}
		if t == 1 {
			v := 1 + rng.Intn(size1)
			fmt.Fprintf(&b, "1 %d %c\n", v, randomChar(rng))
			size1++
		} else {
			v := 1 + rng.Intn(size2)
			fmt.Fprintf(&b, "2 %d %c\n", v, randomChar(rng))
			size2++
		}
	}
	return b.String()
}

func randomChar(rng *rand.Rand) byte {
	return byte('a' + rng.Intn(26))
}
