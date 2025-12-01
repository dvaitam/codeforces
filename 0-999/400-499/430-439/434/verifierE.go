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

const refSource = "./434E.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		want, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		if normalize(got) != normalize(want) {
			fail("wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s", i+1, tc.input, want, got)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "434E-ref-*")
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
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", filepath.Clean(bin))
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20240602))
	primes := []int64{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29,
		97, 101, 251, 509, 911, 1009, 5003, 10007,
		65537, 99991, 199999, 499979, 999983, 1999969,
		9999917, 19999999, 39999997, 79999993, 15485863,
		32452843, 49999991, 67867967, 86028121, 99999989, 999999937,
	}
	var tests []testCase

	tests = append(tests,
		makeTestCase(1, 2, 1, 0, []int64{0}, nil),
		makeTestCase(2, 3, 1, 1, []int64{0, 1}, [][2]int{{1, 2}}),
		makeTestCase(3, 5, 2, 3, []int64{1, 2, 3}, [][2]int{{1, 2}, {2, 3}}),
		makeStarTest(6, 7, 3, 2),
		makePathTest(8, 11, 4, 1),
	)

	// random small
	for i := 0; i < 40; i++ {
		n := rng.Intn(8) + 1
		tests = append(tests, randomTest(rng, n, primes))
	}

	// random medium
	for i := 0; i < 25; i++ {
		n := rng.Intn(300) + 50
		tests = append(tests, randomTest(rng, n, primes))
	}

	// random large
	tests = append(tests, randomTestWithSize(rng, 5000, primes))
	tests = append(tests, randomTestWithSize(rng, 50000, primes))
	tests = append(tests, makePathTest(100000, 99999989, 12345678%99999989+1, 98765432%99999989))
	tests = append(tests, makeStarTest(100000, 999999937, 34567890%999999937+1, 123456789%999999937))
	tests = append(tests, randomTestWithSize(rng, 100000, primes))

	return tests
}

func randomTest(rng *rand.Rand, n int, primes []int64) testCase {
	return randomTestWithParams(rng, n, primes[rng.Intn(len(primes))])
}

func randomTestWithSize(rng *rand.Rand, n int, primes []int64) testCase {
	p := primes[rng.Intn(len(primes))]
	return randomTestWithParams(rng, n, p)
}

func randomTestWithParams(rng *rand.Rand, n int, mod int64) testCase {
	if mod <= 1 {
		mod = 2
	}
	k := int64(rng.Int63n(mod-1) + 1)
	x := int64(rng.Int63n(mod))
	values := make([]int64, n)
	for i := range values {
		values[i] = int64(rng.Int63n(mod))
	}
	edges := randomTreeEdges(rng, n)
	return makeTestCase(n, mod, k, x, values, edges)
}

func randomTreeEdges(rng *rand.Rand, n int) [][2]int {
	if n == 1 {
		return nil
	}
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	return edges
}

func makePathTest(n int, mod int64, k int64, x int64) testCase {
	values := make([]int64, n)
	for i := range values {
		values[i] = int64((i * 1234567) % int(mod))
	}
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	k = normalizeK(k, mod)
	x %= mod
	return makeTestCase(n, mod, k, x, values, edges)
}

func makeStarTest(n int, mod int64, k int64, x int64) testCase {
	values := make([]int64, n)
	for i := range values {
		values[i] = int64((i * 7654321) % int(mod))
	}
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	k = normalizeK(k, mod)
	x %= mod
	return makeTestCase(n, mod, k, x, values, edges)
}

func normalizeK(k, mod int64) int64 {
	k %= mod
	if k == 0 {
		k = 1
	}
	if k == mod {
		k = mod - 1
	}
	if k <= 0 {
		k += mod - 1
	}
	if k >= mod {
		k = mod - 1
	}
	return k
}

func makeTestCase(n int, mod int64, k int64, x int64, values []int64, edges [][2]int) testCase {
	if len(values) != n {
		panic("values length mismatch")
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, mod, k, x%mod)
	for i, v := range values {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v%mod)
	}
	sb.WriteByte('\n')
	if len(edges) == 0 && n > 1 {
		panic("missing edges for non-trivial tree")
	}
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return testCase{input: sb.String()}
}
