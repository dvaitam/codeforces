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

type matrixCase struct {
	n, m int
	a    [][]int
	b    [][]int
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
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
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
	tmp, err := os.CreateTemp("", "cf-2059E1-ref-*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2059E1.go")
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

// -------------------- test generation --------------------

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(2059001))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, simpleAlignments())
	tests = append(tests, shiftedRowsTest())
	tests = append(tests, randomPack("random-small", rng, 8, 4, 6))
	tests = append(tests, randomPack("random-mid", rng, 6, 40, 50))
	tests = append(tests, randomPack("random-large", rng, 3, 150, 400))
	tests = append(tests, heavyCase(rng))

	return tests
}

func sampleTest() testCase {
	var b strings.Builder
	fmt.Fprintln(&b, 3)
	// Case 1
	fmt.Fprintln(&b, "2 2")
	fmt.Fprintln(&b, "2 6")
	fmt.Fprintln(&b, "3 4")
	fmt.Fprintln(&b, "1 2")
	fmt.Fprintln(&b, "7 8")
	// Case 2
	fmt.Fprintln(&b, "2 3")
	fmt.Fprintln(&b, "1 4 7")
	fmt.Fprintln(&b, "8 5 9")
	fmt.Fprintln(&b, "2 3 6")
	fmt.Fprintln(&b, "10 11 12")
	// Case 3
	fmt.Fprintln(&b, "3 3")
	fmt.Fprintln(&b, "1 2 3")
	fmt.Fprintln(&b, "4 5 6")
	fmt.Fprintln(&b, "7 8 9")
	fmt.Fprintln(&b, "10 11 12")
	fmt.Fprintln(&b, "13 14 15")
	fmt.Fprintln(&b, "16 17 18")

	return testCase{name: "sample", input: b.String()}
}

func simpleAlignments() testCase {
	cases := []matrixCase{
		buildSequentialCase(1, 5),
		buildSequentialCase(3, 2),
		buildSequentialCase(2, 4),
	}
	return packCases("simple-alignments", cases)
}

func shiftedRowsTest() testCase {
	// Construct matrices where each b row is a rotation of the corresponding a row.
	n, m := 3, 5
	a := make([][]int, n)
	b := make([][]int, n)
	val := 1
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		b[i] = make([]int, m)
		for j := 0; j < m; j++ {
			a[i][j] = val
			val++
		}
		shift := (i + 2) % m
		for j := 0; j < m; j++ {
			b[i][(j+shift)%m] = a[i][j] + n*m // keep b values distinct from a
		}
		val += m
	}
	return packCases("shifted-rows", []matrixCase{{n: n, m: m, a: a, b: b}})
}

func randomPack(name string, rng *rand.Rand, t, maxN, maxM int) testCase {
	cases := make([]matrixCase, 0, t)
	for i := 0; i < t; i++ {
		n := 1 + rng.Intn(maxN)
		m := 1 + rng.Intn(maxM)
		cases = append(cases, randomCase(n, m, rng))
	}
	return packCases(name, cases)
}

func heavyCase(rng *rand.Rand) testCase {
	// Near the total limit n*m <= 3e5
	n := 300
	m := 1000
	cases := []matrixCase{randomCase(n, m, rng)}
	return packCases("heavy", cases)
}

func buildSequentialCase(n, m int) matrixCase {
	total := n * m
	aVals := make([]int, total)
	for i := 0; i < total; i++ {
		aVals[i] = i + 1
	}
	bVals := make([]int, total)
	for i := 0; i < total; i++ {
		bVals[i] = total + i + 1
	}
	return matrixCase{
		n: n, m: m,
		a: fillMatrix(n, m, aVals),
		b: fillMatrix(n, m, bVals),
	}
}

func randomCase(n, m int, rng *rand.Rand) matrixCase {
	total := n * m
	values := make([]int, 2*total)
	for i := 0; i < 2*total; i++ {
		values[i] = i + 1
	}
	rng.Shuffle(len(values), func(i, j int) { values[i], values[j] = values[j], values[i] })
	aVals := make([]int, total)
	bVals := make([]int, total)
	copy(aVals, values[:total])
	copy(bVals, values[total:])
	return matrixCase{
		n: n, m: m,
		a: fillMatrix(n, m, aVals),
		b: fillMatrix(n, m, bVals),
	}
}

func fillMatrix(n, m int, vals []int) [][]int {
	mat := make([][]int, n)
	idx := 0
	for i := 0; i < n; i++ {
		mat[i] = make([]int, m)
		for j := 0; j < m; j++ {
			mat[i][j] = vals[idx]
			idx++
		}
	}
	return mat
}

func packCases(name string, cases []matrixCase) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.m)
		for i := 0; i < cs.n; i++ {
			for j := 0; j < cs.m; j++ {
				if j > 0 {
					b.WriteByte(' ')
				}
				fmt.Fprintf(&b, "%d", cs.a[i][j])
			}
			b.WriteByte('\n')
		}
		for i := 0; i < cs.n; i++ {
			for j := 0; j < cs.m; j++ {
				if j > 0 {
					b.WriteByte(' ')
				}
				fmt.Fprintf(&b, "%d", cs.b[i][j])
			}
			b.WriteByte('\n')
		}
	}
	return testCase{name: name, input: b.String()}
}
