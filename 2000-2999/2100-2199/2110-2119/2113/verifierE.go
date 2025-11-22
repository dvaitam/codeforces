package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	randomTests = 200
	maxT        = 8
)

type testInput struct {
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}

	candidate, candCleanup, err := prepareBinary(os.Args[1], "candidate2113E")
	if err != nil {
		fmt.Println("failed to prepare contestant binary:", err)
		return
	}
	if candCleanup != nil {
		defer candCleanup()
	}

	oracle, oracleCleanup, err := prepareOracle()
	if err != nil {
		fmt.Println("failed to prepare reference solution:", err)
		return
	}
	defer oracleCleanup()

	tests := deterministicTests()
	total := 0
	for idx, test := range tests {
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("deterministic test %d failed: %v\ninput:\n%s", idx+1, err, test.input)
			return
		}
		total++
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTests; i++ {
		test := randomTest(rng)
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("random test %d failed: %v\ninput:\n%s", i+1, err, test.input)
			return
		}
		total++
	}

	for idx, test := range largeTests() {
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("large test %d failed: %v\ninput length: %d bytes\n", idx+1, err, len(test.input))
			return
		}
		total++
	}

	fmt.Printf("All %d tests passed.\n", total)
}

func prepareBinary(path, prefix string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", tmp, abs)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, nil, nil
}

func prepareOracle() (string, func(), error) {
	dir := sourceDir()
	src := filepath.Join(dir, "2113E.go")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracle2113E_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func runTest(test testInput, candidate, oracle string) error {
	expectOut, err := runBinary(oracle, test.input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	gotOut, err := runBinary(candidate, test.input)
	if err != nil {
		return fmt.Errorf("contestant runtime error: %v", err)
	}

	expect, err := parseOutput(expectOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse oracle output: %v", err)
	}
	got, err := parseOutput(gotOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse contestant output: %v", err)
	}

	for i := 0; i < test.t; i++ {
		if expect[i] != got[i] {
			return fmt.Errorf("case %d: expected %d got %d", i+1, expect[i], got[i])
		}
	}
	return nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func parseOutput(out string, t int) ([]int, error) {
	reader := strings.NewReader(out)
	res := make([]int, 0, t)
	for len(res) < t {
		var x int
		if _, err := fmt.Fscan(reader, &x); err != nil {
			return nil, fmt.Errorf("need %d integers, got %d (%v)", t, len(res), err)
		}
		res = append(res, x)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("expected %d integers, output has extra data", t)
	}
	return res, nil
}

func deterministicTests() []testInput {
	// Sample from the statement.
	sample := `5
4 1 1 4
1 2
2 3
3 4
4 1
5 1 1 5
1 2
2 3
3 4
4 5
5 1
9 2 1 9
1 2
2 3
3 4
3 5
5 6
6 7
6 8
8 9
1 7
1 9
2 7 2 7
1 4
2 5
3 6
4 5
5 6
6 9
2 8
3 7
2 3
2 1
3 1
1 3 1 2
1 2
2 3
1 2
`
	// Small crafted tests to hit waiting / impossible branches.
	small := buildInput([]caseSpec{
		{
			n: 4,
			m: 2,
			x: 1,
			y: 4,
			edges: [][2]int{
				{1, 2},
				{2, 3},
				{3, 4},
			},
			ab: [][2]int{{2, 1}, {3, 4}},
		},
		{
			n: 3,
			m: 1,
			x: 1,
			y: 3,
			edges: [][2]int{
				{1, 2},
				{2, 3},
			},
			ab: [][2]int{{2, 3}},
		},
	})
	return []testInput{
		{input: sample, t: 5},
		small,
	}
}

type caseSpec struct {
	n, m  int
	x, y  int
	edges [][2]int
	ab    [][2]int
}

func buildInput(cases []caseSpec) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d %d %d %d\n", c.n, c.m, c.x, c.y)
		for _, e := range c.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		for _, p := range c.ab {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
	}
	return testInput{input: sb.String(), t: len(cases)}
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(maxT) + 1
	cases := make([]caseSpec, t)
	for i := 0; i < t; i++ {
		cases[i] = randomCase(rng, 50)
	}
	return buildInput(cases)
}

func randomCase(rng *rand.Rand, maxN int) caseSpec {
	n := rng.Intn(maxN-2) + 2
	edges := make([][2]int, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges[v-2] = [2]int{v, p}
	}
	x := rng.Intn(n) + 1
	y := rng.Intn(n-1) + 1
	if y >= x {
		y++
	}
	m := rng.Intn(8) + 1
	ab := make([][2]int, m)
	for i := 0; i < m; i++ {
		a := rng.Intn(n-1) + 1
		if a >= x {
			a++
		}
		b := rng.Intn(n-1) + 1
		if b >= a {
			b++
		}
		if b > n {
			b = n
			if b == a {
				b = 1
			}
		}
		if b == a {
			if a == 1 {
				b = 2
			} else {
				b = 1
			}
		}
		ab[i] = [2]int{a, b}
	}
	return caseSpec{n: n, m: m, x: x, y: y, edges: edges, ab: ab}
}

func largeTests() []testInput {
	// Chain with maximum size and many enemies.
	n := 100000
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		edges[i-2] = [2]int{i - 1, i}
	}
	ab := make([][2]int, 200)
	for i := 0; i < 200; i++ {
		a := (i*407+7)%(n-1) + 1
		if a == 1 { // ensure a != x
			a = 2
		}
		b := n - a
		if b == a || b < 1 {
			b = n
		}
		ab[i] = [2]int{a, b}
	}
	large1 := buildInput([]caseSpec{
		{n: n, m: 200, x: 1, y: n, edges: edges, ab: ab},
	})

	// Random big branching tree with fewer enemies.
	rng := rand.New(rand.NewSource(1))
	n2 := 80000
	edges2 := make([][2]int, n2-1)
	for v := 2; v <= n2; v++ {
		p := rng.Intn(v-1) + 1
		edges2[v-2] = [2]int{v, p}
	}
	x := rng.Intn(n2) + 1
	y := rng.Intn(n2-1) + 1
	if y >= x {
		y++
	}
	ab2 := make([][2]int, 200)
	for i := 0; i < 200; i++ {
		a := rng.Intn(n2-1) + 1
		if a >= x {
			a++
		}
		b := rng.Intn(n2) + 1
		for b == a {
			b = rng.Intn(n2) + 1
		}
		ab2[i] = [2]int{a, b}
	}
	large2 := buildInput([]caseSpec{
		{n: n2, m: 200, x: x, y: y, edges: edges2, ab: ab2},
	})
	return []testInput{large1, large2}
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}
