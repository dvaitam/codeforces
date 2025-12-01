package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "763E.go"

type testCase struct {
	name    string
	input   string
	answers int
}

type pair struct {
	x int
	y int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		expOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expVals, err := parseAnswers(expOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVals, err := parseAnswers(gotOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, gotOut)
			os.Exit(1)
		}

		if len(gotVals) != len(expVals) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d numbers, got %d\ninput:\n%s\n", idx+1, tc.name, len(expVals), len(gotVals), tc.input)
			os.Exit(1)
		}
		for i := range expVals {
			if expVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) at position %d: expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, i+1, expVals[i], gotVals[i], tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "763E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseAnswers(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int, expected)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	tests := []testCase{
		makeCase("single-tree", 1, 1, nil, []pair{{1, 1}, {1, 1}, {1, 1}}),
		chainCase("chain-8", 8),
		noFriendCase("isolated-5", 5, 3),
		denseWindowCase("dense-20", 20, 4),
	}

	tests = append(tests, largeRandomCase())

	rng := rand.New(rand.NewSource(76307630))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("random-%d", i+1)))
	}

	return tests
}

func makeCase(name string, n, k int, edges []pair, queries []pair) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	fmt.Fprintf(&sb, "%d\n", len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.x, e.y)
	}
	fmt.Fprintf(&sb, "%d\n", len(queries))
	for _, q := range queries {
		fmt.Fprintf(&sb, "%d %d\n", q.x, q.y)
	}
	return testCase{name: name, input: sb.String(), answers: len(queries)}
}

func chainCase(name string, n int) testCase {
	edges := make([]pair, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, pair{i, i + 1})
	}
	queries := []pair{
		{1, n},
		{2, n - 1},
		{1, 1},
		{n, n},
		{1, n / 2},
	}
	return makeCase(name, n, 1, edges, queries)
}

func noFriendCase(name string, n, k int) testCase {
	queries := []pair{
		{1, n},
		{2, n},
		{1, 1},
	}
	if n > 2 {
		queries = append(queries, pair{2, 2})
	}
	return makeCase(name, n, k, nil, queries)
}

func denseWindowCase(name string, n, k int) testCase {
	var edges []pair
	for u := 1; u <= n; u++ {
		for d := 1; d <= k && u+d <= n; d++ {
			edges = append(edges, pair{u, u + d})
		}
	}
	queries := []pair{
		{1, n},
		{1, n / 2},
		{n/2 + 1, n},
		{2, n - 1},
	}
	return makeCase(name, n, k, edges, queries)
}

func largeRandomCase() testCase {
	n := 300
	k := 5
	var edges []pair
	for u := 1; u <= n; u++ {
		for d := 1; d <= k && u+d <= n; d++ {
			if (u+d)%3 == 0 || (u*d)%5 == 0 {
				edges = append(edges, pair{u, u + d})
			}
		}
	}
	var queries []pair
	for l := 1; l <= n; l += 30 {
		r := l + 40
		if r > n {
			r = n
		}
		queries = append(queries, pair{l, r})
	}
	queries = append(queries, pair{1, n})
	return makeCase("large-random", n, k, edges, queries)
}

func randomCase(rng *rand.Rand, name string) testCase {
	n := rng.Intn(70) + 1
	k := rng.Intn(5) + 1
	if k > n {
		k = n
	}
	edgeMap := make(map[[2]int]struct{})
	for u := 1; u <= n; u++ {
		maxD := k
		if u+maxD > n {
			maxD = n - u
		}
		for d := 1; d <= maxD; d++ {
			v := u + d
			if rng.Intn(100) < 35 {
				edgeMap[[2]int{u, v}] = struct{}{}
			}
		}
	}
	edges := make([]pair, 0, len(edgeMap))
	for key := range edgeMap {
		edges = append(edges, pair{key[0], key[1]})
	}
	q := rng.Intn(20) + 1
	queries := make([]pair, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = pair{l, r}
	}
	return makeCase(name, n, k, edges, queries)
}
