package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "2000-2999/2000-2099/2030-2039/2030/2030F.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}

	refBin, refCleanup, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := buildBinary(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := generateTests()
	for idx, tc := range tests {
		expect, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		got, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, got)
			os.Exit(1)
		}

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d (%s)\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildBinary(path string) (string, func(), error) {
	cleanPath := filepath.Clean(path)
	if strings.HasSuffix(cleanPath, ".go") {
		tmp, err := os.CreateTemp("", "verifier2030F-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), cleanPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return stdout.String(), nil
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
	var tests []testCase
	tests = append(tests, testCase{
		name:  "sample",
		input: sampleInput(),
	})
	tests = append(tests, buildCase("single_elem", []queryCase{
		{
			n: 1, q: 1,
			a: []int{1},
			queries: [][2]int{
				{1, 1},
			},
		},
	}))
	tests = append(tests, buildCase("small_manual", []queryCase{
		{
			n: 4, q: 4,
			a:       []int{1, 2, 2, 1},
			queries: [][2]int{{1, 4}, {1, 3}, {2, 3}, {2, 4}},
		},
		{
			n: 5, q: 3,
			a:       []int{1, 1, 2, 2, 3},
			queries: [][2]int{{1, 5}, {2, 4}, {3, 5}},
		},
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		tests = append(tests, randomCase(fmt.Sprintf("random_small_%d", i+1), rng, 5, 10))
	}
	for i := 0; i < 4; i++ {
		tests = append(tests, randomCase(fmt.Sprintf("random_mid_%d", i+1), rng, 40, 60))
	}
	tests = append(tests, randomCase("random_large", rng, 200000, 200000))

	return tests
}

func sampleInput() string {
	return `3
4 2
1 2 2 1
1 4
1 3
5 3
1 2 1 2 1
1 2
1 4
1 2
1 1
1
1 1
`
}

type queryCase struct {
	n       int
	q       int
	a       []int
	queries [][2]int
}

func buildCase(name string, cases []queryCase) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.q)
		for i, val := range cs.a {
			if i > 0 {
				fmt.Fprint(&b, " ")
			}
			fmt.Fprint(&b, val)
		}
		fmt.Fprintln(&b)
		for _, qu := range cs.queries {
			fmt.Fprintf(&b, "%d %d\n", qu[0], qu[1])
		}
	}
	return testCase{name: name, input: b.String()}
}

func randomCase(name string, rng *rand.Rand, maxN, maxQ int) testCase {
	t := rng.Intn(5) + 1
	cases := make([]queryCase, 0, t)
	totalN := 0
	totalQ := 0
	for len(cases) < t {
		n := rng.Intn(maxN) + 1
		q := rng.Intn(maxQ) + 1
		if totalN+n > 200000 {
			n = 1
		}
		if totalQ+q > 200000 {
			q = 1
		}
		a := make([]int, n)
		for i := range a {
			a[i] = rng.Intn(n) + 1
		}
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			queries[i] = [2]int{l, r}
		}
		cases = append(cases, queryCase{n: n, q: q, a: a, queries: queries})
		totalN += n
		totalQ += q
	}
	return buildCase(name, cases)
}
