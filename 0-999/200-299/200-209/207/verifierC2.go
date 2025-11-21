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

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine verifier location")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	target := os.Args[1]

	for i, input := range tests {
		expected, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %q, got %q\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-207C2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "207C2.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
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

func generateTests() []string {
	var tests []string
	rng := rand.New(rand.NewSource(207))

	tests = append(tests, sampleTest())
	tests = append(tests, treeOnlyTest(1, 25, rng))
	tests = append(tests, treeOnlyTest(2, 25, rng))
	tests = append(tests, alternatingTest(60, rng))
	tests = append(tests, stagedTest(40, 40, rng))
	tests = append(tests, mirroredChainsTest(30))
	tests = append(tests, randomTest(rng, 120, 0.5))
	tests = append(tests, randomTest(rng, 600, 0.35))
	tests = append(tests, randomTest(rng, 5000, 0.6))
	tests = append(tests, randomTest(rng, 30000, 0.55))

	return tests
}

func sampleTest() string {
	return "5\n1 1 a\n2 1 a\n1 2 b\n2 1 b\n2 3 a\n"
}

func treeOnlyTest(tree, ops int, rng *rand.Rand) string {
	var b strings.Builder
	b.Grow(ops * 8)
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
	b.Grow(ops * 8)
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
	b.Grow(total * 9)
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
	b.Grow(total * 9)
	fmt.Fprintf(&b, "%d\n", total)

	letters := make([]byte, length)
	for i := 0; i < length; i++ {
		letters[i] = byte('a' + (i % 26))
	}

	last1 := 1
	size1 := 1
	for i := 0; i < length; i++ {
		fmt.Fprintf(&b, "1 %d %c\n", last1, letters[i])
		size1++
		last1 = size1
	}

	last2 := 1
	size2 := 1
	for i := 0; i < length; i++ {
		fmt.Fprintf(&b, "2 %d %c\n", last2, letters[length-1-i])
		size2++
		last2 = size2
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
	b.Grow(ops * 9)
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
