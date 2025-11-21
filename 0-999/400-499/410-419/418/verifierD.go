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

const refSource = "0-999/400-499/410-419/418/418D.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "418D-ref-*")
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
	var tests []testCase

	// small deterministic cases
	tests = append(tests, makeTest(2, [][2]int{{1, 2}}, [][2]int{{1, 2}}))
	tests = append(tests, makeTest(3, [][2]int{{1, 2}, {2, 3}}, [][2]int{{1, 3}, {2, 3}}))
	tests = append(tests, makeTest(4, [][2]int{{1, 2}, {1, 3}, {3, 4}}, [][2]int{{1, 2}, {1, 3}, {2, 4}}))

	tests = append(tests, pathTest(10, 12))
	tests = append(tests, starTest(15, 20))

	// random small
	for i := 0; i < 40; i++ {
		n := rng.Intn(10) + 2
		m := rng.Intn(15) + 1
		tests = append(tests, randomTest(rng, n, m))
	}

	// random medium
	for i := 0; i < 30; i++ {
		n := rng.Intn(2000) + 50
		m := rng.Intn(2000) + 50
		tests = append(tests, randomTest(rng, n, m))
	}

	// random large near limits
	tests = append(tests, randomTestWithLimits(rng, 100000, 100000))
	tests = append(tests, pathTest(100000, 100000))
	tests = append(tests, starTest(100000, 100000))

	return tests
}

func randomTest(rng *rand.Rand, n, m int) testCase {
	return randomTestWithLimits(rng, n, m)
}

func randomTestWithLimits(rng *rand.Rand, n, m int) testCase {
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	queries := make([][2]int, m)
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		queries[i] = [2]int{u, v}
	}
	return makeTest(n, edges, queries)
}

func pathTest(n, m int) testCase {
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		edges[i-1] = [2]int{i, i + 1}
	}
	rng := rand.New(rand.NewSource(int64(n) * 97))
	queries := make([][2]int, m)
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		queries[i] = [2]int{u, v}
	}
	return makeTest(n, edges, queries)
}

func starTest(n, m int) testCase {
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		edges[i-2] = [2]int{1, i}
	}
	rng := rand.New(rand.NewSource(int64(n) * 131))
	queries := make([][2]int, m)
	for i := 0; i < m; i++ {
		u := rng.Intn(n-1) + 2
		v := rng.Intn(n-1) + 2
		if v == u {
			v = 1
		}
		if v == 1 && u == 1 {
			v = 2
		}
		if v == 1 {
			queries[i] = [2]int{1, u}
		} else {
			queries[i] = [2]int{u, v}
		}
	}
	return makeTest(n, edges, queries)
}

func makeTest(n int, edges [][2]int, queries [][2]int) testCase {
	var sb strings.Builder
	sb.Grow((n + len(queries)) * 16)
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d\n", len(queries))
	for _, q := range queries {
		if q[0] == q[1] {
			if q[1] == n {
				q[1] = 1
			} else {
				q[1]++
			}
		}
		fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
	}
	return testCase{input: sb.String()}
}
