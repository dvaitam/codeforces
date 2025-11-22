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
		panic("unable to locate verifier path")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "cf-2106C-ref-*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2106C.go")
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

// ----------------- test generation -----------------

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(2106001))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, manualCases())
	tests = append(tests, randomPack("random-small", rng, 8, 20, 20))
	tests = append(tests, randomPack("random-mid", rng, 6, 2000, 1_000_000_000))
	tests = append(tests, randomPack("random-large", rng, 3, 50000, 1_000_000_000))

	return tests
}

func sampleTest() testCase {
	input := "7\n3 10\n1 3 2\n-1 -1 1\n5 10\n1 0 0 1 0\n-1 0 1 0 -1\n5 10\n1 0 0 1 0\n-1 1 -1 1 -1\n5 10\n1 3 2 5 4\n-1 -1 -1 -1 -1\n5 4\n1 3 2 1 3\n1 -1 -1 1 -1\n5 4\n1 3 2 1 3\n2 -1 -1 2 0\n5 5\n5 0 5 4 3\n5 -1 -1 -1 -1\n"
	return testCase{name: "sample", input: input}
}

func manualCases() testCase {
	var cases []testInstance
	// No missing, complementary
	cases = append(cases, newInstance([]int64{2, 1}, []int64{3, 4}, 5))
	// No missing, not complementary
	cases = append(cases, newInstance([]int64{1, 2}, []int64{0, 5}, 5))
	// All missing, interval empty
	cases = append(cases, newInstance([]int64{0, 5}, []int64{-1, -1}, 4))
	// Mixed missing with fixed x within bounds
	cases = append(cases, testInstance{
		n: 4, k: 10,
		a: []int64{3, 1, 4, 2},
		b: []int64{7, -1, 6, -1},
	})
	// Mixed missing with impossible because fixed mismatch
	cases = append(cases, testInstance{
		n: 3, k: 10,
		a: []int64{2, 2, 2},
		b: []int64{5, 6, -1},
	})
	return packCases("manual", cases)
}

type testInstance struct {
	n int
	k int64
	a []int64
	b []int64
}

func newInstance(a, b []int64, k int64) testInstance {
	return testInstance{n: len(a), k: k, a: a, b: b}
}

func randomPack(name string, rng *rand.Rand, t, maxN int, maxK int64) testCase {
	cases := make([]testInstance, 0, t)
	for i := 0; i < t; i++ {
		n := 1 + rng.Intn(maxN)
		k := rng.Int63n(maxK + 1)
		cases = append(cases, randomCase(rng, n, k))
	}
	return packCases(name, cases)
}

func randomCase(rng *rand.Rand, n int, k int64) testInstance {
	a := make([]int64, n)
	b := make([]int64, n)

	// Decide whether to make it valid or invalid
	valid := rng.Intn(4) != 0 // 75% valid

	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(k + 1)
	}

	if valid {
		L := a[0]
		R := a[0] + k
		for i := 1; i < n; i++ {
			if a[i] > L {
				L = a[i]
			}
			if a[i]+k < R {
				R = a[i] + k
			}
		}
		if L > R {
			valid = false
		} else {
			x := L + rng.Int63n(R-L+1)
			for i := 0; i < n; i++ {
				val := x - a[i]
				if rng.Intn(4) == 0 {
					b[i] = -1
				} else {
					b[i] = val
				}
			}
		}
	}

	if !valid {
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				b[i] = rng.Int63n(k + 1)
			} else {
				b[i] = -1
			}
		}
		// Force inconsistency by setting two fixed different sums if possible.
		if n >= 2 {
			b[0] = 0
			b[1] = k
		}
	}

	return testInstance{n: n, k: k, a: a, b: b}
}

func packCases(name string, cases []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.k)
		for i, v := range cs.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		for i, v := range cs.b {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return testCase{name: name, input: b.String()}
}
